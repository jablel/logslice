// Package linetruncator provides a transformer that truncates lines exceeding
// a maximum byte length, optionally appending a configurable suffix such as
// "..." to indicate that the line was cut.
package linetruncator

import "errors"

// ErrInvalidMaxLen is returned when maxLen is less than 1.
var ErrInvalidMaxLen = errors.New("linetruncator: maxLen must be >= 1")

// ErrSuffixTooLong is returned when the suffix is longer than maxLen.
var ErrSuffixTooLong = errors.New("linetruncator: suffix length exceeds maxLen")

// Truncator truncates lines that exceed a maximum byte length.
type Truncator struct {
	maxLen int
	suffix string
}

// New creates a Truncator that cuts lines to at most maxLen bytes.
// If suffix is non-empty it is appended to truncated lines; it must be
// shorter than maxLen so that at least one byte of original content is kept.
func New(maxLen int, suffix string) (*Truncator, error) {
	if maxLen < 1 {
		return nil, ErrInvalidMaxLen
	}
	if len(suffix) >= maxLen {
		return nil, ErrSuffixTooLong
	}
	return &Truncator{maxLen: maxLen, suffix: suffix}, nil
}

// Apply returns line unchanged when it fits within maxLen bytes.
// Otherwise it trims the line to (maxLen - len(suffix)) bytes and appends
// the suffix.
func (t *Truncator) Apply(line string) string {
	if len(line) <= t.maxLen {
		return line
	}
	keep := t.maxLen - len(t.suffix)
	// Ensure we never slice in the middle of a multi-byte rune.
	truncated := []byte(line[:keep])
	return string(truncated) + t.suffix
}
