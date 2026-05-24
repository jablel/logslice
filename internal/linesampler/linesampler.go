// Package linesampler provides probabilistic line sampling based on a
// configurable keep-rate between 0.0 and 1.0. Unlike the deterministic
// sampler (which keeps every Nth line), linesampler uses a hash of the
// line content to decide whether to keep it, ensuring reproducible results
// for the same input regardless of stream position.
package linesampler

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math"
)

// ErrInvalidRate is returned when the keep rate is outside [0.0, 1.0].
var ErrInvalidRate = fmt.Errorf("linesampler: rate must be between 0.0 and 1.0 inclusive")

// LineSampler decides whether to keep a line based on a probabilistic rate.
type LineSampler struct {
	rate      float64
	threshold uint64
}

// New creates a LineSampler with the given keep rate.
// A rate of 1.0 keeps all lines; 0.0 drops all lines.
func New(rate float64) (*LineSampler, error) {
	if rate < 0.0 || rate > 1.0 {
		return nil, ErrInvalidRate
	}
	threshold := uint64(math.Round(rate * float64(math.MaxUint64)))
	return &LineSampler{rate: rate, threshold: threshold}, nil
}

// Keep returns true if the line should be kept based on its content hash.
// The decision is deterministic: the same line always produces the same result.
func (s *LineSampler) Keep(line string) bool {
	if s.rate == 1.0 {
		return true
	}
	if s.rate == 0.0 {
		return false
	}
	h := md5.Sum([]byte(line))
	v := binary.LittleEndian.Uint64(h[:8])
	return v <= s.threshold
}

// Rate returns the configured keep rate.
func (s *LineSampler) Rate() float64 {
	return s.rate
}
