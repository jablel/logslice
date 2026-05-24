// Package timestampinjector prepends or appends a parsed timestamp to each
// log line in a normalised format, making heterogeneous log streams easier
// to compare and sort downstream.
package timestampinjector

import (
	"fmt"
	"strings"
	"time"

	"github.com/user/logslice/internal/timestampextractor"
)

// Position controls where the normalised timestamp is injected.
type Position int

const (
	// Prepend places the timestamp at the beginning of the line.
	Prepend Position = iota
	// Append places the timestamp at the end of the line.
	Append
)

// Injector rewrites each line by injecting a normalised timestamp.
type Injector struct {
	extractor  *timestampextractor.Extractor
	outFormat string
	position  Position
}

// New creates an Injector that uses the provided extractor to locate
// timestamps and re-emits them in outFormat (a Go time layout string).
// If outFormat is empty, time.RFC3339 is used.
func New(ex *timestampextractor.Extractor, outFormat string, pos Position) (*Injector, error) {
	if ex == nil {
		return nil, fmt.Errorf("timestampinjector: extractor must not be nil")
	}
	if outFormat == "" {
		outFormat = time.RFC3339
	}
	return &Injector{
		extractor:  ex,
		outFormat: outFormat,
		position:  pos,
	}, nil
}

// Apply rewrites line by injecting the normalised timestamp.
// If no timestamp is found the original line is returned unchanged.
func (inj *Injector) Apply(line string) string {
	t, ok := inj.extractor.Extract(line)
	if !ok {
		return line
	}
	stamp := t.Format(inj.outFormat)
	line = strings.TrimRight(line, "\n")
	switch inj.position {
	case Append:
		return line + " " + stamp
	default: // Prepend
		return stamp + " " + line
	}
}
