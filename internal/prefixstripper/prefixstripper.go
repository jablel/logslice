// Package prefixstripper removes a fixed string or regex prefix from each log line.
// This is useful when log lines contain a leading label, host, or process name
// that should be stripped before further processing.
package prefixstripper

import (
	"fmt"
	"regexp"
	"strings"
)

// Stripper removes a prefix from log lines.
type Stripper struct {
	fixed  string
	re     *regexp.Regexp
	useRe  bool
}

// New creates a Stripper that removes the given prefix from each line.
// If isRegex is true, prefix is compiled as a regular expression that must
// match at the start of the line (^ is prepended automatically if absent).
// Returns an error if the regex is invalid.
func New(prefix string, isRegex bool) (*Stripper, error) {
	if prefix == "" {
		return nil, fmt.Errorf("prefixstripper: prefix must not be empty")
	}
	if !isRegex {
		return &Stripper{fixed: prefix}, nil
	}
	// Anchor to start if not already anchored.
	pattern := prefix
	if !strings.HasPrefix(pattern, "^") {
		pattern = "^" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("prefixstripper: invalid regex: %w", err)
	}
	return &Stripper{re: re, useRe: true}, nil
}

// Apply removes the configured prefix from line and returns the result.
// If the prefix is not present the original line is returned unchanged.
func (s *Stripper) Apply(line string) string {
	if s.useRe {
		loc := s.re.FindStringIndex(line)
		if loc == nil {
			return line
		}
		return line[loc[1]:]
	}
	after, found := strings.CutPrefix(line, s.fixed)
	if !found {
		return line
	}
	return after
}
