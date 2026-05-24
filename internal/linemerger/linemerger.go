// Package linemerger merges continuation lines into a single logical log entry.
// Some log formats (e.g. Java stack traces, multi-line JSON) emit a single
// logical event across several physical lines. LineMerger joins those lines
// using a configurable continuation predicate.
package linemerger

import (
	"errors"
	"regexp"
	"strings"
)

// ErrNilPredicate is returned when no continuation predicate is provided.
var ErrNilPredicate = errors.New("linemerger: continuation predicate must not be nil")

// Merger accumulates physical lines and emits logical records.
type Merger struct {
	isContinuation func(line string) bool
	delimiter      string
	pending        []string
}

// New creates a Merger. isContinuation returns true when a line is a
// continuation of the previous record (e.g. starts with whitespace or matches
// a pattern). delimiter is inserted between joined lines (default "\n").
func New(isContinuation func(line string) bool, delimiter string) (*Merger, error) {
	if isContinuation == nil {
		return nil, ErrNilPredicate
	}
	if delimiter == "" {
		delimiter = "\n"
	}
	return &Merger{
		isContinuation: isContinuation,
		delimiter:      delimiter,
	}, nil
}

// NewIndentMerger returns a Merger that treats lines beginning with whitespace
// as continuations — the most common heuristic for stack traces.
func NewIndentMerger() (*Merger, error) {
	return New(func(line string) bool {
		return len(line) > 0 && (line[0] == ' ' || line[0] == '\t')
	}, "\n")
}

// NewPatternMerger returns a Merger that treats lines matching pattern as
// continuations.
func NewPatternMerger(pattern string) (*Merger, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return New(func(line string) bool {
		return re.MatchString(line)
	}, "\n")
}

// Feed accepts the next physical line. It returns a completed logical record
// and true when one is ready, or "", false when the line was buffered.
func (m *Merger) Feed(line string) (string, bool) {
	if m.isContinuation(line) {
		m.pending = append(m.pending, line)
		return "", false
	}
	// Flush previous record, start a new one.
	if len(m.pending) == 0 {
		m.pending = append(m.pending, line)
		return "", false
	}
	record := strings.Join(m.pending, m.delimiter)
	m.pending = []string{line}
	return record, true
}

// Flush returns any buffered record that has not yet been emitted. Call this
// after the input stream is exhausted.
func (m *Merger) Flush() (string, bool) {
	if len(m.pending) == 0 {
		return "", false
	}
	record := strings.Join(m.pending, m.delimiter)
	m.pending = m.pending[:0]
	return record, true
}
