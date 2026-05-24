package multipatternfilter

import (
	"fmt"
	"regexp"
)

// Filter matches lines against multiple patterns with AND or OR logic.
type Filter struct {
	patterns []*regexp.Regexp
	mode     Mode
}

// Mode controls how multiple patterns are combined.
type Mode int

const (
	// ModeAnd requires all patterns to match.
	ModeAnd Mode = iota
	// ModeOr requires at least one pattern to match.
	ModeOr
)

// New creates a Filter that matches lines using the given regex patterns.
// mode must be ModeAnd or ModeOr. patterns must be non-empty.
func New(patterns []string, mode Mode) (*Filter, error) {
	if len(patterns) == 0 {
		return nil, ErrNoPatterns
	}
	if mode != ModeAnd && mode != ModeOr {
		return nil, fmt.Errorf("multipatternfilter: unknown mode %d", mode)
	}
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("multipatternfilter: invalid pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return &Filter{patterns: compiled, mode: mode}, nil
}

// Keep returns true if the line satisfies the pattern combination.
func (f *Filter) Keep(line string) bool {
	switch f.mode {
	case ModeAnd:
		for _, re := range f.patterns {
			if !re.MatchString(line) {
				return false
			}
		}
		return true
	case ModeOr:
		for _, re := range f.patterns {
			if re.MatchString(line) {
				return true
			}
		}
		return false
	}
	return false
}

// Len returns the number of compiled patterns.
func (f *Filter) Len() int { return len(f.patterns) }
