// Package lineskipper skips the first N lines of a stream before
// passing subsequent lines through. This is useful for bypassing
// file headers or comment blocks at the start of a log file.
package lineskipper

import "errors"

// ErrNegativeSkip is returned when the skip count is negative.
var ErrNegativeSkip = errors.New("lineskipper: skip count must be non-negative")

// Skipper tracks how many lines have been seen and suppresses the
// first N of them.
type Skipper struct {
	skip int
	seen int
}

// New creates a Skipper that discards the first skip lines.
// A skip value of 0 is valid and means every line is kept.
func New(skip int) (*Skipper, error) {
	if skip < 0 {
		return nil, ErrNegativeSkip
	}
	return &Skipper{skip: skip}, nil
}

// Keep returns true once the first skip lines have been seen.
// It increments the internal counter on every call.
func (s *Skipper) Keep(line string) bool {
	if s.seen < s.skip {
		s.seen++
		return false
	}
	return true
}

// Skipped returns the number of lines that have been discarded so far.
func (s *Skipper) Skipped() int {
	if s.seen < s.skip {
		return s.seen
	}
	return s.skip
}

// Reset resets the internal counter so the Skipper can be reused
// against a new stream.
func (s *Skipper) Reset() {
	s.seen = 0
}
