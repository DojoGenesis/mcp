package authz

import (
	"sync"
	"time"
)

// Limiter is a per-label token bucket for dispatch-class tool calls.
// Refill rate is perMin tokens per minute with burst capacity of perMin
// (minimum 1). A nil *Limiter allows everything.
type Limiter struct {
	mu      sync.Mutex
	perMin  float64
	burst   float64
	buckets map[string]*bucket
	now     func() time.Time
}

type bucket struct {
	tokens float64
	last   time.Time
}

// NewLimiter builds a limiter allowing perMin calls per minute per label.
// perMin values <= 0 fall back to 6.
func NewLimiter(perMin int) *Limiter {
	if perMin <= 0 {
		perMin = 6
	}
	return &Limiter{
		perMin:  float64(perMin),
		burst:   float64(perMin),
		buckets: make(map[string]*bucket),
		now:     time.Now,
	}
}

// Allow takes one token from label's bucket, reporting whether the call may
// proceed.
func (l *Limiter) Allow(label string) bool {
	if l == nil {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	b, ok := l.buckets[label]
	if !ok {
		b = &bucket{tokens: l.burst, last: now}
		l.buckets[label] = b
	}

	elapsed := now.Sub(b.last).Seconds()
	if elapsed > 0 {
		b.tokens += elapsed * l.perMin / 60.0
		if b.tokens > l.burst {
			b.tokens = l.burst
		}
		b.last = now
	}

	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}
