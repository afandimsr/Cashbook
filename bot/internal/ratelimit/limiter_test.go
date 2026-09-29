package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimiter_AllowsUpToMax(t *testing.T) {
	l := New(3, time.Minute)

	assert.True(t, l.Allow(1))
	assert.True(t, l.Allow(1))
	assert.True(t, l.Allow(1))
	assert.False(t, l.Allow(1), "4th call within the window must be rejected")
}

func TestLimiter_TracksKeysIndependently(t *testing.T) {
	l := New(1, time.Minute)

	assert.True(t, l.Allow(1))
	assert.False(t, l.Allow(1))
	assert.True(t, l.Allow(2), "a different key must have its own budget")
}

func TestLimiter_ResetsAfterWindow(t *testing.T) {
	current := time.Now()
	l := New(1, time.Minute)
	l.now = func() time.Time { return current }

	assert.True(t, l.Allow(1))
	assert.False(t, l.Allow(1))

	current = current.Add(2 * time.Minute)
	assert.True(t, l.Allow(1), "old hits outside the window must be pruned")
}
