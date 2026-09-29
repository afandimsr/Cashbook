package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/afandimsr/cashbook-bot/internal/backend"
	"github.com/afandimsr/cashbook-bot/internal/config"
	"github.com/afandimsr/cashbook-bot/internal/llm"
	"github.com/afandimsr/cashbook-bot/internal/scheduler"
	"github.com/afandimsr/cashbook-bot/internal/service"
	"github.com/afandimsr/cashbook-bot/internal/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// schedulerCheckInterval controls how often the scheduler wakes up to check
// whether it's the 1st of the month; it doesn't need to be frequent.
const schedulerCheckInterval = 1 * time.Hour

func main() {
	cfg := config.Load()

	tgAPI, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("telegram: %v", err)
	}
	tgAPI.Debug = cfg.AppEnv != "production"

	llmClient, err := buildLLM(context.Background(), cfg)
	if err != nil {
		log.Fatalf("llm: %v", err)
	}

	backendClient := backend.New(cfg.BackendBaseURL, cfg.BackendInternalKey, cfg.BackendTimeout)
	svc := service.New(backendClient, llmClient)
	bot := telegram.New(tgAPI, svc)

	sched := scheduler.New(backendClient, svc, bot)
	go sched.Run(context.Background(), schedulerCheckInterval)

	log.Printf("CashBook bot running as @%s", tgAPI.Self.UserName)
	bot.Run()
}

// buildLLM creates one client per entry in cfg.LLMProviders, in that order.
// Names are already validated by config.Load.
func buildLLM(ctx context.Context, cfg *config.Config) (*llm.Chain, error) {
	providers := make([]llm.Provider, 0, len(cfg.LLMProviders))
	for _, name := range cfg.LLMProviders {
		switch name {
		case config.ProviderGemini:
			c, err := llm.NewGemini(ctx, cfg.GeminiAPIKey, cfg.GeminiModel)
			if err != nil {
				return nil, fmt.Errorf("gemini: %w", err)
			}
			providers = append(providers, llm.Provider{Name: name, Client: c})
		case config.ProviderGroq:
			c := llm.NewOpenAICompat(cfg.GroqBaseURL, cfg.GroqAPIKey, cfg.GroqModel, cfg.LLMProviderTimeout)
			providers = append(providers, llm.Provider{Name: name, Client: c})
		default:
			return nil, fmt.Errorf("unknown provider %q", name)
		}
	}
	log.Printf("llm providers (in fallback order): %s", strings.Join(cfg.LLMProviders, " -> "))
	return llm.NewChain(cfg.LLMProviderTimeout, providers...), nil
}
