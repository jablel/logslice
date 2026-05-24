// Package linebuffer provides a fixed-capacity ring buffer for log lines.
// It is useful for keeping the last N lines seen during streaming, enabling
// sliding-window operations without unbounded memory growth.
package linebuffer

import "errors"

// ErrInvalidCapacity is returned when capacity is less than 1.
var ErrInvalidCapacity = errors.New("linebuffer: capacity must be at least 1")

// Buffer is a fixed-capacity ring buffer that retains the last N lines.
type Buffer struct {
	buf  []string
	head int
	size int
	cap  int
}

// New creates a Buffer with the given capacity.
// Returns ErrInvalidCapacity if cap < 1.
func New(capacity int) (*Buffer, error) {
	if capacity < 1 {
		return nil, ErrInvalidCapacity
	}
	return &Buffer{
		buf: make([]string, capacity),
		cap: capacity,
	}, nil
}

// Push adds a line to the buffer, evicting the oldest entry when full.
func (b *Buffer) Push(line string) {
	if b.size < b.cap {
		b.buf[(b.head+b.size)%b.cap] = line
		b.size++
	} else {
		b.buf[b.head] = line
		b.head = (b.head + 1) % b.cap
	}
}

// Lines returns all buffered lines in insertion order (oldest first).
func (b *Buffer) Lines() []string {
	out := make([]string, b.size)
	for i := 0; i < b.size; i++ {
		out[i] = b.buf[(b.head+i)%b.cap]
	}
	return out
}

// Len returns the number of lines currently held.
func (b *Buffer) Len() int { return b.size }

// Cap returns the maximum capacity of the buffer.
func (b *Buffer) Cap() int { return b.cap }

// Reset clears all buffered lines without reallocating.
func (b *Buffer) Reset() {
	b.head = 0
	b.size = 0
}
