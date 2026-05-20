// Package scanner provides line-by-line scanning of log files with
// time-based filtering capabilities.
package scanner

import (
	"bufio"
	"io"
	"time"

	"github.com/user/logslice/internal/timeparser"
)

// LineResult holds a parsed log line along with its extracted timestamp.
type LineResult struct {
	Line      string
	Timestamp time.Time
	LineNum   int
}

// Options configures the behavior of the Scanner.
type Options struct {
	// Format is an optional explicit time format string. If empty, all
	// KnownFormats are tried in order.
	Format string

	// Start is the beginning of the desired time range (inclusive).
	// A zero value means no lower bound.
	Start time.Time

	// End is the end of the desired time range (inclusive).
	// A zero value means no upper bound.
	End time.Time

	// SkipUnparseable controls whether lines whose timestamps cannot be
	// extracted are silently dropped (true) or passed through as-is (false).
	SkipUnparseable bool
}

// Scanner reads log lines from a reader and filters them by time range.
type Scanner struct {
	opts    Options
	reader  *bufio.Scanner
	current LineResult
	err     error
	lineNum int
}

// New creates a new Scanner that reads from r using the provided options.
func New(r io.Reader, opts Options) *Scanner {
	return &Scanner{
		opts:   opts,
		reader: bufio.NewScanner(r),
	}
}

// Scan advances to the next matching line. It returns true if a line is
// available via Line(), and false when the input is exhausted or an error
// occurs.
func (s *Scanner) Scan() bool {
	for s.reader.Scan() {
		s.lineNum++
		raw := s.reader.Text()

		var ts time.Time
		var parseErr error

		if s.opts.Format != "" {
			ts, parseErr = timeparser.ParseWithFormat(raw, s.opts.Format)
		} else {
			ts, parseErr = timeparser.Parse(raw)
		}

		if parseErr != nil {
			// Line has no recognisable timestamp.
			if s.opts.SkipUnparseable {
				continue
			}
			// Pass through without range filtering.
			s.current = LineResult{Line: raw, LineNum: s.lineNum}
			return true
		}

		// Apply time-range filter only when at least one bound is set.
		if !s.opts.Start.IsZero() || !s.opts.End.IsZero() {
			if !timeparser.InRange(ts, s.opts.Start, s.opts.End) {
				continue
			}
		}

		s.current = LineResult{Line: raw, Timestamp: ts, LineNum: s.lineNum}
		return true
	}

	s.err = s.reader.Err()
	return false
}

// Line returns the most recent line produced by a call to Scan.
func (s *Scanner) Line() LineResult {
	return s.current
}

// Err returns the first non-EOF error encountered by the scanner.
func (s *Scanner) Err() error {
	return s.err
}
