// Package stats collects and reports processing statistics for a log slice run.
package stats

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Stats holds counters accumulated during a log-slicing operation.
type Stats struct {
	LinesScanned  int64
	LinesMatched  int64
	LinesSkipped  int64
	BytesRead     int64
	StartTime     time.Time
	EndTime       time.Time
}

// New returns a new Stats instance with the start time set to now.
func New() *Stats {
	return &Stats{StartTime: time.Now()}
}

// Finish records the end time of the operation.
func (s *Stats) Finish() {
	s.EndTime = time.Now()
}

// Duration returns the elapsed time between Start and Finish.
func (s *Stats) Duration() time.Duration {
	if s.EndTime.IsZero() {
		return time.Since(s.StartTime)
	}
	return s.EndTime.Sub(s.StartTime)
}

// MatchRate returns the fraction of scanned lines that matched, in [0,1].
func (s *Stats) MatchRate() float64 {
	if s.LinesScanned == 0 {
		return 0
	}
	return float64(s.LinesMatched) / float64(s.LinesScanned)
}

// Print writes a human-readable summary to w.
// If w is nil, os.Stderr is used.
func (s *Stats) Print(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	fmt.Fprintf(w, "--- logslice stats ---\n")
	fmt.Fprintf(w, "  scanned : %d lines (%d bytes)\n", s.LinesScanned, s.BytesRead)
	fmt.Fprintf(w, "  matched : %d lines (%.1f%%)\n", s.LinesMatched, s.MatchRate()*100)
	fmt.Fprintf(w, "  skipped : %d lines\n", s.LinesSkipped)
	fmt.Fprintf(w, "  elapsed : %s\n", s.Duration().Round(time.Millisecond))
}
