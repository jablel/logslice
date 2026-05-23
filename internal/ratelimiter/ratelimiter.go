// Package ratelimiter provides a token-bucket style rate limiter that
// throttles the number of log lines emitted per second during processing.
package ratelimiter

import (
	"time"
)

// RateLimiter controls the throughput of log lines by enforcing a maximum
// number of lines allowed per second. A rate of zero means unlimited.
type RateLimiter struct {
	rate     int
	count    int
	windowStart time.Time
	now      func() time.Time
	sleep    func(time.Duration)
}

// New creates a RateLimiter that allows at most ratePerSec lines per second.
// A ratePerSec of 0 disables rate limiting entirely.
// Returns an error if ratePerSec is negative.
func New(ratePerSec int) (*RateLimiter, error) {
	if ratePerSec < 0 {
		return nil, ErrInvalidRate
	}
	return &RateLimiter{
		rate:        ratePerSec,
		windowStart: time.Now(),
		now:         time.Now,
		sleep:       time.Sleep,
	}, nil
}

// Wait blocks if the current rate would exceed the configured limit.
// It is a no-op when the rate is zero (unlimited).
func (r *RateLimiter) Wait() {
	if r.rate == 0 {
		return
	}

	now := r.now()
	elapsed := now.Sub(r.windowStart)

	if elapsed >= time.Second {
		r.windowStart = now
		r.count = 0
	}

	r.count++

	if r.count > r.rate {
		sleepUntil := r.windowStart.Add(time.Second)
		wait := sleepUntil.Sub(r.now())
		if wait > 0 {
			r.sleep(wait)
		}
		r.windowStart = r.now()
		r.count = 1
	}
}

// Rate returns the configured rate limit (lines per second).
func (r *RateLimiter) Rate() int {
	return r.rate
}

// Count returns the number of lines processed in the current window.
func (r *RateLimiter) Count() int {
	return r.count
}
