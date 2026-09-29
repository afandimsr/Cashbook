package llm

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/afandimsr/cashbook-bot/internal/backend"
)

type parser interface {
	ParseTransaction(ctx context.Context, message string, categories []backend.Category) (*ParsedTransaction, error)
}

// Provider is one named entry in a Chain, e.g. {"gemini", *GeminiClient}.
type Provider struct {
	Name   string
	Client parser
}

// Chain tries providers in order so a rate-limited or failing free tier
// (e.g. Gemini's quota) falls through to the next one instead of failing the
// user's message.
type Chain struct {
	providers []Provider
	timeout   time.Duration
}

// NewChain gives each provider attempt its own perProviderTimeout, so one
// hung provider can't consume the whole request budget.
func NewChain(perProviderTimeout time.Duration, providers ...Provider) *Chain {
	return &Chain{providers: providers, timeout: perProviderTimeout}
}

// ParseTransaction returns the first provider's successful result. It stops
// early on ErrNoFunctionCall: that's the model deciding the message isn't a
// transaction, and another model would only repeat the same verdict. Any
// other error (rate limit, 5xx, timeout, malformed output) moves on to the
// next provider.
func (c *Chain) ParseTransaction(ctx context.Context, message string, categories []backend.Category) (*ParsedTransaction, error) {
	var failures []error
	for _, p := range c.providers {
		if err := ctx.Err(); err != nil {
			failures = append(failures, err)
			break
		}

		parsed, err := c.attempt(ctx, p, message, categories)
		if err == nil {
			log.Printf("llm: served by %s", p.Name)
			return parsed, nil
		}
		if errors.Is(err, ErrNoFunctionCall) {
			return nil, err
		}
		log.Printf("llm: %s failed: %v", p.Name, err)
		failures = append(failures, fmt.Errorf("%s: %w", p.Name, err))
	}
	// Joined with %w so callers can still errors.Is(err, ErrInvalidOutput).
	return nil, fmt.Errorf("all llm providers failed: %w", errors.Join(failures...))
}

func (c *Chain) attempt(ctx context.Context, p Provider, message string, categories []backend.Category) (*ParsedTransaction, error) {
	attemptCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return p.Client.ParseTransaction(attemptCtx, message, categories)
}
