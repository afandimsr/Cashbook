package ratelimit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLimiter_Allow_UpToMax(t *testing.T) {
	l := New(3, time.Minute)

	assert.True(t, l.Allow("k"))
	assert.True(t, l.Allow("k"))
	assert.True(t, l.Allow("k"))
	assert.False(t, l.Allow("k"))
}

func TestLimiter_Allow_TracksKeysIndependently(t *testing.T) {
	l := New(1, time.Minute)

	assert.True(t, l.Allow("a"))
	assert.False(t, l.Allow("a"))
	assert.True(t, l.Allow("b"))
}

func TestLimiter_Allow_ResetsAfterWindow(t *testing.T) {
	current := time.Now()
	l := New(1, time.Minute)
	l.now = func() time.Time { return current }

	assert.True(t, l.Allow("k"))
	assert.False(t, l.Allow("k"))

	current = current.Add(2 * time.Minute)
	assert.True(t, l.Allow("k"))
}

func TestLimiter_Blocked_DoesNotRecord(t *testing.T) {
	l := New(1, time.Minute)

	assert.False(t, l.Blocked("k"), "no attempts recorded yet")
	assert.False(t, l.Blocked("k"), "Blocked must not itself count as an attempt")
}

func TestLimiter_RecordFailure_ThenBlocked(t *testing.T) {
	l := New(2, time.Minute)

	assert.False(t, l.Blocked("k"))
	l.RecordFailure("k")
	assert.False(t, l.Blocked("k"))
	l.RecordFailure("k")
	assert.True(t, l.Blocked("k"))
}

func TestLimiter_Reset_ClearsFailures(t *testing.T) {
	l := New(1, time.Minute)

	l.RecordFailure("k")
	assert.True(t, l.Blocked("k"))

	l.Reset("k")
	assert.False(t, l.Blocked("k"))
}
