// Package telegram wires incoming Telegram updates to the bot's service
// layer and turns its replies back into outgoing messages. It never talks to
// the backend or Gemini directly — that's all in internal/service.
package telegram

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/afandimsr/cashbook-bot/internal/ratelimit"
	"github.com/afandimsr/cashbook-bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// requestTimeout covers one full message: backend calls plus up to one
// attempt per LLM provider in the fallback chain (LLM_PROVIDER_TIMEOUT each).
const requestTimeout = 30 * time.Second

// A 6-character code from a 32-symbol alphabet is already ~1.07B combinations
// (computationally infeasible to guess), but the bot always holds the valid
// internal API key, so it's the one place a code-guessing attempt could reach
// the backend at all. This throttle removes that path entirely rather than
// relying on the code space alone.
const (
	linkAttemptsMax    = 5
	linkAttemptsWindow = 10 * time.Minute
)

type Bot struct {
	api         *tgbotapi.BotAPI
	svc         *service.Bot
	linkLimiter *ratelimit.Limiter
}

func New(api *tgbotapi.BotAPI, svc *service.Bot) *Bot {
	return &Bot{
		api:         api,
		svc:         svc,
		linkLimiter: ratelimit.New(linkAttemptsMax, linkAttemptsWindow),
	}
}

// Run starts long-polling for updates and blocks until the process is
// terminated. Each update is handled synchronously and independently, so a
// slow or failed LLM/backend call for one user never blocks another.
func (b *Bot) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		go b.handleMessage(update.Message)
	}
}

// SendMessage delivers a message outside the update loop, e.g. from the
// scheduler's monthly report job (see internal/scheduler).
func (b *Bot) SendMessage(chatID int64, text string) error {
	_, err := b.api.Send(tgbotapi.NewMessage(chatID, text))
	return err
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	chatID := msg.Chat.ID
	reply, err := b.dispatch(ctx, chatID, msg)
	if err != nil {
		log.Printf("chat %d: %v", chatID, err)
		reply = "Maaf, terjadi kesalahan di server. Coba lagi sebentar lagi."
	}
	if reply == "" {
		return
	}

	if _, err := b.api.Send(tgbotapi.NewMessage(chatID, reply)); err != nil {
		log.Printf("chat %d: send reply: %v", chatID, err)
	}
}

func (b *Bot) dispatch(ctx context.Context, chatID int64, msg *tgbotapi.Message) (string, error) {
	if msg.IsCommand() {
		switch msg.Command() {
		case "start":
			return b.svc.WhoAmI(ctx, chatID, msg.From.FirstName, msg.From.LastName)
		case "help":
			return service.HelpMessage, nil
		case "link":
			code := strings.TrimSpace(msg.CommandArguments())
			if code == "" {
				return "Kirim dengan format: /link <kode> (lihat kode di aplikasi CashBook).", nil
			}
			if !b.linkLimiter.Allow(chatID) {
				return "Terlalu banyak percobaan. Coba lagi dalam beberapa menit.", nil
			}
			return b.svc.Link(ctx, chatID, msg.From.UserName, code)
		case "summary":
			return b.svc.Summary(ctx, chatID)
		case "budget", "anggaran":
			now := time.Now()
			return b.svc.Budgets(ctx, chatID, int(now.Month()), now.Year())
		case "kategori", "categories":
			return b.svc.Categories(ctx, chatID)
		case "report":
			now := time.Now()
			return b.svc.MonthlyReport(ctx, chatID, int(now.Month()), now.Year())
		default:
			return "Perintah tidak dikenal.\n\n" + service.HelpMessage, nil
		}
	}

	text := strings.TrimSpace(msg.Text)
	if text == "" {
		return "", nil
	}
	return b.svc.HandleMessage(ctx, chatID, text)
}
