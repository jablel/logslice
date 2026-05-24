// Package fieldfilter filters log lines by matching a specific field value.
// It supports JSON and logfmt structured log formats and can be used to
// include or exclude lines where a named field matches a given pattern.
package fieldfilter

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Filter keeps or drops lines based on whether a named field matches a pattern.
type Filter struct {
	field   string
	re      *regexp.Regexp
	invert  bool
}

// New creates a Filter that matches lines where the given field value matches
// pattern. If invert is true, lines that do NOT match are kept instead.
// pattern is treated as a regular expression.
func New(field, pattern string, invert bool) (*Filter, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldfilter: field name must not be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("fieldfilter: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("fieldfilter: invalid pattern: %w", err)
	}
	return &Filter{field: field, re: re, invert: invert}, nil
}

// Keep returns true if the line should be kept according to the filter rules.
func (f *Filter) Keep(line string) bool {
	val, ok := extractField(line, f.field)
	if !ok {
		return f.invert
	}
	matched := f.re.MatchString(val)
	if f.invert {
		return !matched
	}
	return matched
}

// extractField attempts to extract a field value from a JSON or logfmt line.
func extractField(line, field string) (string, bool) {
	line = strings.TrimSpace(line)
	if len(line) > 0 && line[0] == '{' {
		return extractJSON(line, field)
	}
	return extractLogfmt(line, field)
}

func extractJSON(line, field string) (string, bool) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return "", false
	}
	v, ok := m[field]
	if !ok {
		return "", false
	}
	return fmt.Sprintf("%v", v), true
}

func extractLogfmt(line, field string) (string, bool) {
	prefix := field + "="
	for _, part := range strings.Fields(line) {
		if strings.HasPrefix(part, prefix) {
			val := strings.TrimPrefix(part, prefix)
			val = strings.Trim(val, `"`)
			return val, true
		}
	}
	return "", false
}
