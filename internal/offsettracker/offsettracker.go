// Package offsettracker tracks byte offsets of matched lines within a log file,
// enabling fast re-seeking and resumable log slicing operations.
package offsettracker

import "fmt"

// Entry records a matched line's position within the source file.
type Entry struct {
	LineNumber int64
	ByteOffset int64
	Line       string
}

// Tracker accumulates byte-offset entries for matched log lines.
type Tracker struct {
	entries  []Entry
	current  int64 // running byte offset
	enabled  bool
}

// New creates a new Tracker. If enabled is false, Record is a no-op and
// the tracker acts as a transparent pass-through.
func New(enabled bool) *Tracker {
	return &Tracker{enabled: enabled}
}

// Advance moves the internal byte cursor forward by n bytes (the length of
// a consumed line including its newline delimiter).
func (t *Tracker) Advance(n int) {
	t.current += int64(n)
}

// Record stores the current byte offset together with the line content and
// its logical line number. It is a no-op when the tracker is disabled.
func (t *Tracker) Record(lineNumber int64, line string) {
	if !t.enabled {
		return
	}
	t.entries = append(t.entries, Entry{
		LineNumber: lineNumber,
		ByteOffset: t.current,
		Line:       line,
	})
}

// Entries returns a copy of all recorded offset entries.
func (t *Tracker) Entries() []Entry {
	out := make([]Entry, len(t.entries))
	copy(out, t.entries)
	return out
}

// Reset clears all recorded entries and resets the byte cursor to zero.
func (t *Tracker) Reset() {
	t.entries = t.entries[:0]
	t.current = 0
}

// Summary returns a human-readable summary of tracked offsets.
func (t *Tracker) Summary() string {
	if !t.enabled {
		return "offset tracking disabled"
	}
	return fmt.Sprintf("tracked %d entries, last offset %d bytes", len(t.entries), t.current)
}
