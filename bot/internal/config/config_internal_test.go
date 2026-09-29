package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func validConfig() *Config {
	return &Config{
		TelegramBotToken:   "tg",
		BackendInternalKey: "key",
		LLMProviders:       []string{ProviderGemini},
		LLMProviderTimeout: 10 * time.Second,
		GeminiAPIKey:       "gemini-key",
	}
}

func TestParseList(t *testing.T) {
	assert.Equal(t, []string{"gemini", "groq"}, parseList(" Gemini, Groq ,"))
	assert.Nil(t, parseList(""))
}

func TestValidate_GeminiOnly(t *testing.T) {
	assert.NoError(t, validate(validConfig()))
}

func TestValidate_GroqRequiresKey(t *testing.T) {
	cfg := validConfig()
	cfg.LLMProviders = []string{ProviderGemini, ProviderGroq}
	assert.ErrorContains(t, validate(cfg), "GROQ_API_KEY")

	cfg.GroqAPIKey = "groq-key"
	assert.NoError(t, validate(cfg))
}

func TestValidate_GeminiKeyNotRequiredWhenUnlisted(t *testing.T) {
	cfg := validConfig()
	cfg.GeminiAPIKey = ""
	cfg.LLMProviders = []string{ProviderGroq}
	cfg.GroqAPIKey = "groq-key"
	assert.NoError(t, validate(cfg))
}

// OpenCode Zen's free tier rejects any client other than OpenCode itself
// (403 FreeTierError), so it was removed; a leftover config must fail loudly.
func TestValidate_OpenCodeNoLongerSupported(t *testing.T) {
	cfg := validConfig()
	cfg.LLMProviders = []string{ProviderGemini, "opencode"}
	assert.ErrorContains(t, validate(cfg), `unknown provider "opencode"`)
}

func TestValidate_RejectsUnknownDuplicateAndEmpty(t *testing.T) {
	cfg := validConfig()
	cfg.LLMProviders = []string{"claude"}
	assert.ErrorContains(t, validate(cfg), "unknown provider")

	cfg.LLMProviders = []string{ProviderGemini, ProviderGemini}
	assert.ErrorContains(t, validate(cfg), "more than once")

	cfg.LLMProviders = nil
	assert.ErrorContains(t, validate(cfg), "at least one")
}
