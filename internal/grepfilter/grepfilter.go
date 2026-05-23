// Package grepfilter provides line-level substring and regex matching
// for filtering log lines by their raw content.
package grepfilter

import (
	"fmt"
	"regexp"
)

// Filter matches log lines against a fixed substring or a compiled regular
// expression pattern. When Invert is true the sense of the match is flipped,
// mimicking grep -v behaviour.
type Filter struct {
	pattern *regexp.Regexp
	invert  bool
}

// Options controls how the Filter is constructed.
type Options struct {
	// Pattern is a literal string or regular expression to match against.
	Pattern string
	// Regex treats Pattern as a regular expression when true; otherwise a
	// fixed-string (substring) match is performed.
	Regex bool
	// Invert inverts the match: only lines that do NOT match are kept.
	Invert bool
	// CaseInsensitive wraps the pattern with (?i) when true.
	CaseInsensitive bool
}

// New constructs a Filter from the provided Options.
// Returns an error if the pattern cannot be compiled.
func New(opts Options) (*Filter, error) {
	if opts.Pattern == "" {
		return nil, fmt.Errorf("grepfilter: pattern must not be empty")
	}

	raw := opts.Pattern
	if !opts.Regex {
		raw = regexp.QuoteMeta(raw)
	}
	if opts.CaseInsensitive {
		raw = "(?i)" + raw
	}

	re, err := regexp.Compile(raw)
	if err != nil {
		return nil, fmt.Errorf("grepfilter: compile pattern: %w", err)
	}

	return &Filter{pattern: re, invert: opts.Invert}, nil
}

// Keep returns true when the line should be forwarded to the next stage.
func (f *Filter) Keep(line string) bool {
	matched := f.pattern.MatchString(line)
	if f.invert {
		return !matched
	}
	return matched
}

// Pattern returns the compiled regular expression used for matching.
func (f *Filter) Pattern() *regexp.Regexp {
	return f.pattern
}
