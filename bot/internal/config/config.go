package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// LLM provider names accepted in LLM_PROVIDERS.
const (
	ProviderGemini = "gemini"
	ProviderGroq   = "groq"
)

type Config struct {
	AppEnv string

	TelegramBotToken string

	// LLMProviders is the ordered fallback list from LLM_PROVIDERS, e.g.
	// ["gemini", "groq"]: the first is tried first, the rest only when
	// the previous one errors (see llm.Chain).
	LLMProviders       []string
	LLMProviderTimeout time.Duration

	GeminiAPIKey string
	GeminiModel  string

	GroqAPIKey  string
	GroqBaseURL string
	GroqModel   string

	BackendBaseURL     string
	BackendInternalKey string
	BackendTimeout     time.Duration
}

func Load() *Config {
	_ = godotenv.Load()

	timeout, err := time.ParseDuration(getEnv("LLM_PROVIDER_TIMEOUT", "10s"))
	if err != nil {
		log.Fatalf("LLM_PROVIDER_TIMEOUT: %v", err)
	}

	cfg := &Config{
		AppEnv: getEnv("APP_ENV", "development"),

		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),

		LLMProviders:       parseList(getEnv("LLM_PROVIDERS", ProviderGemini)),
		LLMProviderTimeout: timeout,

		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
		GeminiModel:  getEnv("GEMINI_MODEL", "gemini-2.5-flash"),

		GroqAPIKey:  getEnv("GROQ_API_KEY", ""),
		GroqBaseURL: getEnv("GROQ_BASE_URL", "https://api.groq.com/openai/v1"),
		GroqModel:   getEnv("GROQ_MODEL", "openai/gpt-oss-120b"),

		BackendBaseURL:     getEnv("BACKEND_BASE_URL", "http://localhost:8181"),
		BackendInternalKey: getEnv("BACKEND_INTERNAL_API_KEY", ""),
		BackendTimeout:     10 * time.Second,
	}

	if err := validate(cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// parseList splits "gemini, Groq" into ["gemini", "groq"].
func parseList(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.ToLower(strings.TrimSpace(p)); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func validate(cfg *Config) error {
	if cfg.TelegramBotToken == "" {
		return errors.New("TELEGRAM_BOT_TOKEN is required")
	}
	if cfg.BackendInternalKey == "" {
		return errors.New("BACKEND_INTERNAL_API_KEY is required")
	}
	if len(cfg.LLMProviders) == 0 {
		return errors.New("LLM_PROVIDERS must list at least one provider")
	}
	if cfg.LLMProviderTimeout <= 0 {
		return errors.New("LLM_PROVIDER_TIMEOUT must be positive")
	}

	seen := make(map[string]bool)
	for _, p := range cfg.LLMProviders {
		if seen[p] {
			return fmt.Errorf("LLM_PROVIDERS lists %q more than once", p)
		}
		seen[p] = true

		switch p {
		case ProviderGemini:
			if cfg.GeminiAPIKey == "" {
				return errors.New("GEMINI_API_KEY is required when LLM_PROVIDERS includes gemini")
			}
		case ProviderGroq:
			if cfg.GroqAPIKey == "" {
				return errors.New("GROQ_API_KEY is required when LLM_PROVIDERS includes groq")
			}
		default:
			return fmt.Errorf("LLM_PROVIDERS: unknown provider %q (supported: %s, %s)", p, ProviderGemini, ProviderGroq)
		}
	}
	return nil
}
