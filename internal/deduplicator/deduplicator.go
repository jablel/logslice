// Package deduplicator provides line-level deduplication for log output.
// It tracks seen lines within a configurable window and skips exact duplicates.
package deduplicator

import "errors"

// Deduplicator filters out consecutive or windowed duplicate log lines.
type Deduplicator struct {
	window  int
	seen    []string
	head    int
	size    int
	skipped int
}

// New creates a Deduplicator with the given look-back window size.
// window must be >= 1; use 1 to deduplicate only consecutive identical lines.
func New(window int) (*Deduplicator, error) {
	if window < 1 {
		return nil, errors.New("deduplicator: window must be >= 1")
	}
	return &Deduplicator{
		window: window,
		seen:   make([]string, window),
	}, nil
}

// Keep returns true if the line should be kept (i.e. not a duplicate within the window).
func (d *Deduplicator) Keep(line string) bool {
	for i := 0; i < d.size; i++ {
		idx := (d.head - 1 - i + d.window) % d.window
		if d.seen[idx] == line {
			d.skipped++
			return false
		}
	}
	d.seen[d.head] = line
	d.head = (d.head + 1) % d.window
	if d.size < d.window {
		d.size++
	}
	return true
}

// Skipped returns the number of duplicate lines dropped so far.
func (d *Deduplicator) Skipped() int {
	return d.skipped
}

// Reset clears the internal window, allowing the deduplicator to be reused.
func (d *Deduplicator) Reset() {
	d.seen = make([]string, d.window)
	d.head = 0
	d.size = 0
	d.skipped = 0
}
