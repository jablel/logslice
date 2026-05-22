// Package sampler provides line-sampling functionality for logslice.
// It allows extracting every Nth line from a matched result set,
// which is useful for reducing output volume on high-frequency logs.
package sampler

import "errors"

// Sampler holds the configuration for line sampling.
type Sampler struct {
	step    int
	counter int
}

// New creates a new Sampler that keeps every nth line.
// step must be >= 1; a step of 1 keeps every line (no-op sampling).
func New(step int) (*Sampler, error) {
	if step < 1 {
		return nil, errors.New("sampler: step must be >= 1")
	}
	return &Sampler{step: step}, nil
}

// Keep returns true if the current line should be kept,
// advancing the internal counter regardless.
func (s *Sampler) Keep() bool {
	s.counter++
	if s.counter >= s.step {
		s.counter = 0
		return true
	}
	return false
}

// Reset resets the internal counter to its initial state.
func (s *Sampler) Reset() {
	s.counter = 0
}

// Step returns the configured sampling step.
func (s *Sampler) Step() int {
	return s.step
}
