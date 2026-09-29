// Package ratelimit provides a small in-memory per-key sliding-window
// limiter. It's intentionally not distributed — good enough for throttling
// abuse from a single bot process (e.g. someone spamming /link guesses),
// which is the only thing it's used for here.
package ratelimit

import (
	"sync"
	"time"
)

type Limiter struct {
	mu     sync.Mutex
	max    int
	window time.Duration
	hits   map[int64][]time.Time
	now    func() time.Time
}

func New(max int, window time.Duration) *Limiter {
	return &Limiter{
		max:    max,
		window: window,
		hits:   make(map[int64][]time.Time),
		now:    time.Now,
	}
}

// Allow reports whether a call for the given key is allowed right now, and
// records it if so.
func (l *Limiter) Allow(key int64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	cutoff := now.Add(-l.window)

	pruned := l.hits[key][:0]
	for _, h := range l.hits[key] {
		if h.After(cutoff) {
			pruned = append(pruned, h)
		}
	}

	if len(pruned) >= l.max {
		l.hits[key] = pruned
		return false
	}

	l.hits[key] = append(pruned, now)
	return true
}
