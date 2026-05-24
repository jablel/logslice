// Package lineannotator prepends a configurable prefix tag to each log line.
// It is useful for labelling lines by source file, host, or stream name
// when merging output from multiple origins.
package lineannotator

import (
	"errors"
	"fmt"
)

// ErrEmptyTag is returned when an empty tag is provided to New.
var ErrEmptyTag = errors.New("lineannotator: tag must not be empty")

// Annotator prepends a bracketed tag to every line it processes.
type Annotator struct {
	tag    string
	prefix string
}

// New creates an Annotator that will prepend [tag] to each line.
// Returns ErrEmptyTag if tag is the empty string.
func New(tag string) (*Annotator, error) {
	if tag == "" {
		return nil, ErrEmptyTag
	}
	return &Annotator{
		tag:    tag,
		prefix: fmt.Sprintf("[%s] ", tag),
	}, nil
}

// Apply returns line with the configured tag prefix prepended.
// Empty lines are returned unchanged.
func (a *Annotator) Apply(line string) string {
	if line == "" {
		return line
	}
	return a.prefix + line
}

// Tag returns the raw tag string (without brackets or trailing space).
func (a *Annotator) Tag() string {
	return a.tag
}
