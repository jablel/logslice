// Package jsonpretty provides a formatter that pretty-prints JSON log lines
// with optional indentation and field ordering for human-readable output.
package jsonpretty

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Formatter pretty-prints JSON lines.
type Formatter struct {
	indent  string
	prefix  string
	disabled bool
}

// New creates a Formatter with the given indent string (e.g. "  " or "\t").
// If indent is empty the formatter is effectively a no-op pass-through.
func New(indent string) (*Formatter, error) {
	if len(indent) > 8 {
		return nil, fmt.Errorf("jsonpretty: indent too long (max 8 chars), got %d", len(indent))
	}
	return &Formatter{
		indent:   indent,
		prefix:   "",
		disabled: indent == "",
	}, nil
}

// Apply formats line as pretty-printed JSON if it is valid JSON and the
// formatter is enabled. Otherwise the original line is returned unchanged.
func (f *Formatter) Apply(line string) string {
	if f.disabled {
		return line
	}
	raw := []byte(line)
	if !json.Valid(raw) {
		return line
	}
	var buf bytes.Buffer
	if err := json.Indent(&buf, raw, f.prefix, f.indent); err != nil {
		return line
	}
	return buf.String()
}

// IsEnabled reports whether pretty-printing is active.
func (f *Formatter) IsEnabled() bool {
	return !f.disabled
}
