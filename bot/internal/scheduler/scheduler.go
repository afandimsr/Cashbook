// Package scheduler sends every linked user an automatic expense report on
// the 1st of each month, covering the month that just closed.
package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/afandimsr/cashbook-bot/internal/backend"
	"github.com/afandimsr/cashbook-bot/internal/service"
)

// Sender delivers a message to a Telegram chat; satisfied by *telegram.Bot.
type Sender interface {
	SendMessage(chatID int64, text string) error
}

// Scheduler tracks "already sent this month" in memory, per chat. That's a
// deliberate trade-off: it resets on restart, but this is a single
// best-effort monthly notification (not a financial record — the underlying
// data is always available on-demand via /report), so it doesn't warrant a
// backend schema change just for idempotency.
type Scheduler struct {
	backend service.BackendClient
	svc     *service.Bot
	sender  Sender
	now     func() time.Time

	mu       sync.Mutex
	lastSent map[int64]string // chatID -> "YYYY-MM"
}

func New(backendClient service.BackendClient, svc *service.Bot, sender Sender) *Scheduler {
	return &Scheduler{
		backend:  backendClient,
		svc:      svc,
		sender:   sender,
		now:      time.Now,
		lastSent: make(map[int64]string),
	}
}

// Run blocks, checking once per checkInterval whether today is the 1st of
// the month and, if so, sending any chat that hasn't been sent this month
// yet. It also checks once immediately on start.
func (s *Scheduler) Run(ctx context.Context, checkInterval time.Duration) {
	s.tick(ctx)

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context) {
	now := s.now()
	if !shouldRunToday(now) {
		return
	}
	monthKey := now.Format("2006-01")
	month, year := previousMonth(now)

	links, err := s.backend.ListLinks(ctx)
	if err != nil {
		log.Printf("scheduler: list links: %v", err)
		return
	}

	for _, link := range links {
		if s.alreadySent(link.TelegramChatID, monthKey) {
			continue
		}
		s.sendReport(ctx, link, month, year, monthKey)
	}
}

func (s *Scheduler) sendReport(ctx context.Context, link backend.TelegramLink, month, year int, monthKey string) {
	msg, err := s.svc.MonthlyReport(ctx, link.TelegramChatID, month, year)
	if err != nil {
		log.Printf("scheduler: report for chat %d: %v", link.TelegramChatID, err)
		return
	}

	if err := s.sender.SendMessage(link.TelegramChatID, "Laporan bulanan otomatis dari CashBook:\n\n"+msg); err != nil {
		log.Printf("scheduler: send to chat %d: %v", link.TelegramChatID, err)
		return
	}

	s.markSent(link.TelegramChatID, monthKey)
}

func (s *Scheduler) alreadySent(chatID int64, monthKey string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastSent[chatID] == monthKey
}

func (s *Scheduler) markSent(chatID int64, monthKey string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastSent[chatID] = monthKey
}

func shouldRunToday(now time.Time) bool {
	return now.Day() == 1
}

// previousMonth returns the month/year that just closed, e.g. on 2026-09-01
// it returns (8, 2026).
func previousMonth(now time.Time) (month, year int) {
	prev := now.AddDate(0, -1, 0)
	return int(prev.Month()), prev.Year()
}
