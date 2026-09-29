package llm_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/afandimsr/cashbook-bot/internal/backend"
	"github.com/afandimsr/cashbook-bot/internal/llm"
	"github.com/stretchr/testify/assert"
)

// fakeParser returns a fixed result, or blocks until its context is done
// when hang is set.
type fakeParser struct {
	result *llm.ParsedTransaction
	err    error
	hang   bool
	calls  int
}

func (f *fakeParser) ParseTransaction(ctx context.Context, _ string, _ []backend.Category) (*llm.ParsedTransaction, error) {
	f.calls++
	if f.hang {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	return f.result, f.err
}

var parsed = &llm.ParsedTransaction{Amount: 50000, Type: "expense", CategoryName: "Makanan"}

func TestChain_FirstProviderSucceeds(t *testing.T) {
	first := &fakeParser{result: parsed}
	second := &fakeParser{result: parsed}
	c := llm.NewChain(time.Second, llm.Provider{Name: "a", Client: first}, llm.Provider{Name: "b", Client: second})

	got, err := c.ParseTransaction(context.Background(), "beli 50rb", nil)
	assert.NoError(t, err)
	assert.Equal(t, parsed, got)
	assert.Equal(t, 0, second.calls)
}

func TestChain_FallsBackOnError(t *testing.T) {
	first := &fakeParser{err: errors.New("status 429")}
	second := &fakeParser{result: parsed}
	c := llm.NewChain(time.Second, llm.Provider{Name: "a", Client: first}, llm.Provider{Name: "b", Client: second})

	got, err := c.ParseTransaction(context.Background(), "beli 50rb", nil)
	assert.NoError(t, err)
	assert.Equal(t, parsed, got)
	assert.Equal(t, 1, first.calls)
	assert.Equal(t, 1, second.calls)
}

func TestChain_StopsOnNoFunctionCall(t *testing.T) {
	first := &fakeParser{err: llm.ErrNoFunctionCall}
	second := &fakeParser{result: parsed}
	c := llm.NewChain(time.Second, llm.Provider{Name: "a", Client: first}, llm.Provider{Name: "b", Client: second})

	_, err := c.ParseTransaction(context.Background(), "jam 5", nil)
	assert.ErrorIs(t, err, llm.ErrNoFunctionCall)
	assert.Equal(t, 0, second.calls)
}

func TestChain_AllFail_KeepsUnderlyingErrors(t *testing.T) {
	first := &fakeParser{err: errors.New("status 429")}
	second := &fakeParser{err: llm.ErrInvalidOutput}
	c := llm.NewChain(time.Second, llm.Provider{Name: "a", Client: first}, llm.Provider{Name: "b", Client: second})

	_, err := c.ParseTransaction(context.Background(), "beli 50rb", nil)
	assert.Error(t, err)
	assert.ErrorIs(t, err, llm.ErrInvalidOutput)
	assert.Contains(t, err.Error(), "a: status 429")
}

func TestChain_PerProviderTimeoutFallsBack(t *testing.T) {
	first := &fakeParser{hang: true}
	second := &fakeParser{result: parsed}
	c := llm.NewChain(20*time.Millisecond, llm.Provider{Name: "a", Client: first}, llm.Provider{Name: "b", Client: second})

	got, err := c.ParseTransaction(context.Background(), "beli 50rb", nil)
	assert.NoError(t, err)
	assert.Equal(t, parsed, got)
}

func TestChain_StopsWhenParentContextDone(t *testing.T) {
	second := &fakeParser{result: parsed}
	c := llm.NewChain(time.Second, llm.Provider{Name: "a", Client: &fakeParser{err: errors.New("boom")}}, llm.Provider{Name: "b", Client: second})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := c.ParseTransaction(ctx, "beli 50rb", nil)
	assert.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, 0, second.calls)
}
