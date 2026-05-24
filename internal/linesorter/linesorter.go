package linesorter

import (
	"fmt"
	"sort"
)

// Order defines the sort direction.
type Order int

const (
	Ascending  Order = iota
	Descending Order = iota
)

// Sorter buffers lines and returns them sorted.
type Sorter struct {
	order  Order
	buffer []string
}

// New creates a Sorter with the given order.
// Returns an error if the order is not Ascending or Descending.
func New(order Order) (*Sorter, error) {
	if order != Ascending && order != Descending {
		return nil, fmt.Errorf("linesorter: unknown order %d", order)
	}
	return &Sorter{order: order}, nil
}

// Feed appends a line to the internal buffer.
func (s *Sorter) Feed(line string) {
	s.buffer = append(s.buffer, line)
}

// Flush returns all buffered lines in sorted order and clears the buffer.
func (s *Sorter) Flush() []string {
	if len(s.buffer) == 0 {
		return nil
	}
	out := make([]string, len(s.buffer))
	copy(out, s.buffer)
	if s.order == Ascending {
		sort.Strings(out)
	} else {
		sort.Sort(sort.Reverse(sort.StringSlice(out)))
	}
	s.buffer = s.buffer[:0]
	return out
}

// Len returns the number of buffered lines.
func (s *Sorter) Len() int {
	return len(s.buffer)
}

// Reset discards all buffered lines.
func (s *Sorter) Reset() {
	s.buffer = s.buffer[:0]
}
