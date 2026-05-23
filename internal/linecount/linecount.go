// Package linecount provides a simple counter that tracks the total number
// of lines seen and the number of lines that passed a filter, enabling
// progress reporting during large log file processing.
package linecount

import "sync/atomic"

// Counter tracks total and matched line counts in a thread-safe manner.
type Counter struct {
	total   atomic.Int64
	matched atomic.Int64
}

// New creates a new Counter with all counts at zero.
func New() *Counter {
	return &Counter{}
}

// Inc increments the total line count by one.
func (c *Counter) Inc() {
	c.total.Add(1)
}

// Match increments both the total and matched line counts by one.
func (c *Counter) Match() {
	c.total.Add(1)
	c.matched.Add(1)
}

// Total returns the total number of lines seen.
func (c *Counter) Total() int64 {
	return c.total.Load()
}

// Matched returns the number of lines that passed the filter.
func (c *Counter) Matched() int64 {
	return c.matched.Load()
}

// Skipped returns the number of lines that did not pass the filter.
func (c *Counter) Skipped() int64 {
	return c.total.Load() - c.matched.Load()
}

// MatchRate returns the ratio of matched lines to total lines as a value
// between 0.0 and 1.0. Returns 0 if no lines have been seen.
func (c *Counter) MatchRate() float64 {
	t := c.total.Load()
	if t == 0 {
		return 0
	}
	return float64(c.matched.Load()) / float64(t)
}

// Reset zeroes all counters.
func (c *Counter) Reset() {
	c.total.Store(0)
	c.matched.Store(0)
}
