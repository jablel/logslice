// Package fieldextractor provides utilities for extracting named fields
// from structured log lines (JSON, logfmt, or key=value formats).
package fieldextractor

import (
	"encoding/json"
	"strings"
)

// Format represents the detected or specified log line format.
type Format int

const (
	FormatUnknown Format = iota
	FormatJSON
	FormatLogfmt
)

// Extractor extracts field values from log lines.
type Extractor struct {
	format Format
	field  string
}

// New creates a new Extractor for the given field name.
// If format is FormatUnknown, it will be auto-detected per line.
func New(field string, format Format) (*Extractor, error) {
	if field == "" {
		return nil, ErrEmptyField
	}
	return &Extractor{field: field, format: format}, nil
}

// Extract returns the value of the configured field from the given log line.
// Returns empty string if the field is not found.
func (e *Extractor) Extract(line string) string {
	fmt := e.format
	if fmt == FormatUnknown {
		fmt = detect(line)
	}
	switch fmt {
	case FormatJSON:
		return extractJSON(line, e.field)
	case FormatLogfmt:
		return extractLogfmt(line, e.field)
	}
	return ""
}

// detect guesses the format of a log line.
func detect(line string) Format {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) > 0 && trimmed[0] == '{' {
		return FormatJSON
	}
	return FormatLogfmt
}

// extractJSON extracts a top-level string field from a JSON log line.
func extractJSON(line, field string) string {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return ""
	}
	if v, ok := m[field]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// extractLogfmt extracts a field value from a key=value log line.
func extractLogfmt(line, field string) string {
	prefix := field + "="
	for _, part := range strings.Fields(line) {
		if strings.HasPrefix(part, prefix) {
			val := part[len(prefix):]
			return strings.Trim(val, `"`)
		}
	}
	return ""
}
