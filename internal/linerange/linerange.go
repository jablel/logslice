// Package linerange provides line-number based range filtering,
// allowing extraction of specific line intervals from a log stream.
package linerange

import (
	"errors"
	"fmt"
)

// ErrInvalidRange is returned when the specified range is invalid.
var ErrInvalidRange = errors.New("linerange: invalid range")

// Filter keeps only lines whose 1-based index falls within [From, To].
// A zero value for To means "no upper bound".
type Filter struct {
	from    int
	to      int
	current int
}

// New creates a Filter that passes lines numbered from..to (inclusive, 1-based).
// to == 0 means unlimited upper bound.
func New(from, to int) (*Filter, error) {
	if from < 1 {
		return nil, fmt.Errorf("%w: from must be >= 1, got %d", ErrInvalidRange, from)
	}
	if to != 0 && to < from {
		return nil, fmt.Errorf("%w: to (%d) must be >= from (%d)", ErrInvalidRange, to, from)
	}
	return &Filter{from: from, to: to}, nil
}

// Keep increments the internal line counter and reports whether the current
// line falls within the configured range. Once past the upper bound the
// filter is done and will never return true again.
func (f *Filter) Keep(_ string) bool {
	f.current++
	if f.current < f.from {
		return false
	}
	if f.to != 0 && f.current > f.to {
		return false
	}
	return true
}

// Done reports whether the upper bound has been reached. Callers may use this
// to short-circuit scanning when no further lines can match.
func (f *Filter) Done() bool {
	return f.to != 0 && f.current >= f.to
}

// Reset restores the filter to its initial state so it can be reused.
func (f *Filter) Reset() {
	f.current = 0
}

// Current returns the number of lines seen so far.
func (f *Filter) Current() int {
	return f.current
}
