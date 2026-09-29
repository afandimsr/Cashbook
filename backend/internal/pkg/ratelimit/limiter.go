// Package ratelimit provides a small in-memory per-key sliding-window
// limiter. It only coordinates within a single process — fine for this app's
// current single-instance deployment, but a horizontally-scaled deployment
// would need a shared store (e.g. Redis) instead.
package ratelimit

import (
	"sync"
	"time"
)

type Limiter struct {
	mu     sync.Mutex
	max    int
	window time.Duration
	hits   map[string][]time.Time
	now    func() time.Time
}

func New(max int, window time.Duration) *Limiter {
	return &Limiter{
		max:    max,
		window: window,
		hits:   make(map[string][]time.Time),
		now:    time.Now,
	}
}

// Allow reports whether a call for key is allowed right now, and records it
// (whether it "succeeded" or not) if so. Use this for a plain per-request
// throttle, e.g. keyed by client IP.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.isBlockedLocked(key) {
		return false
	}
	l.hits[key] = append(l.pruneLocked(key), l.now())
	return true
}

// Blocked reports whether key is currently at/over the limit, without
// recording a new attempt. Pairs with RecordFailure for "N consecutive
// failures" lockouts (e.g. login), where only failed attempts should count.
func (l *Limiter) Blocked(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.isBlockedLocked(key)
}

// RecordFailure records one failed attempt for key.
func (l *Limiter) RecordFailure(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.hits[key] = append(l.pruneLocked(key), l.now())
}

// Reset clears key's recorded attempts, e.g. after a successful login.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.hits, key)
}

// isBlockedLocked prunes expired hits for key and reports whether it's at or
// over the limit. Caller must hold l.mu.
func (l *Limiter) isBlockedLocked(key string) bool {
	pruned := l.pruneLocked(key)
	l.hits[key] = pruned
	return len(pruned) >= l.max
}

// pruneLocked returns key's hits that are still within the window. Caller
// must hold l.mu.
func (l *Limiter) pruneLocked(key string) []time.Time {
	cutoff := l.now().Add(-l.window)
	hits := l.hits[key]
	out := hits[:0]
	for _, h := range hits {
		if h.After(cutoff) {
			out = append(out, h)
		}
	}
	return out
}
