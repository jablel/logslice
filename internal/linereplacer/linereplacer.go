// Package linereplacer provides a transformer that replaces substrings or
// regex patterns within log lines with a fixed replacement string.
package linereplacer

import (
	"fmt"
	"regexp"
	"strings"
)

// Replacer replaces occurrences of a pattern within each line.
type Replacer struct {
	re          *regexp.Regexp
	fixed       string
	replacement string
	useRegex    bool
}

// New creates a Replacer that substitutes all occurrences of pattern with
// replacement. If isRegex is true the pattern is compiled as a regular
// expression; otherwise a plain substring replacement is performed.
//
// Returns an error when isRegex is true and the pattern fails to compile, or
// when pattern is empty.
func New(pattern, replacement string, isRegex bool) (*Replacer, error) {
	if pattern == "" {
		return nil, fmt.Errorf("linereplacer: pattern must not be empty")
	}
	if isRegex {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("linereplacer: invalid regex %q: %w", pattern, err)
		}
		return &Replacer{re: re, replacement: replacement, useRegex: true}, nil
	}
	return &Replacer{fixed: pattern, replacement: replacement, useRegex: false}, nil
}

// Apply returns line with all occurrences of the configured pattern replaced
// by the replacement string. Lines that do not contain the pattern are
// returned unchanged.
func (r *Replacer) Apply(line string) string {
	if r.useRegex {
		return r.re.ReplaceAllString(line, r.replacement)
	}
	return strings.ReplaceAll(line, r.fixed, r.replacement)
}
