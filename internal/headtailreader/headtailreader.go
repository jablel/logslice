// Package headtailreader provides utilities to read the first N or last N
// lines from a log stream, enabling quick previews without full scans.
package headtailreader

import (
	"errors"
)

// ErrInvalidCount is returned when n is not a positive integer.
var ErrInvalidCount = errors.New("headtailreader: count must be greater than zero")

// Reader holds configuration for head/tail extraction.
type Reader struct {
	n    int
	mode string // "head" or "tail"
	buf  []string
	seen int
}

// New creates a Reader in the given mode ("head" or "tail") keeping n lines.
func New(mode string, n int) (*Reader, error) {
	if n <= 0 {
		return nil, ErrInvalidCount
	}
	if mode != "head" && mode != "tail" {
		return nil, errors.New("headtailreader: mode must be \"head\" or \"tail\"")
	}
	return &Reader{n: n, mode: mode, buf: make([]string, 0, n)}, nil
}

// Feed accepts a line. For head mode it keeps the first n lines; for tail mode
// it maintains a rolling window of the last n lines.
func (r *Reader) Feed(line string) {
	r.seen++
	switch r.mode {
	case "head":
		if len(r.buf) < r.n {
			r.buf = append(r.buf, line)
		}
	case "tail":
		if len(r.buf) < r.n {
			r.buf = append(r.buf, line)
		} else {
			// Shift left and append.
			copy(r.buf, r.buf[1:])
			r.buf[r.n-1] = line
		}
	}
}

// Lines returns the collected lines in order.
func (r *Reader) Lines() []string {
	out := make([]string, len(r.buf))
	copy(out, r.buf)
	return out
}

// Seen returns the total number of lines fed.
func (r *Reader) Seen() int { return r.seen }

// Reset clears internal state so the Reader can be reused.
func (r *Reader) Reset() {
	r.buf = r.buf[:0]
	r.seen = 0
}
