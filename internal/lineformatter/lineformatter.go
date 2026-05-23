// Package lineformatter provides utilities for reformatting log lines
// by extracting and reordering fields into a normalized output format.
package lineformatter

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Format defines the output format for reformatted lines.
type Format int

const (
	// FormatRaw returns the line unchanged.
	FormatRaw Format = iota
	// FormatJSON reformats the line as compact JSON.
	FormatJSON
	// FormatText reformats the line as key=value pairs.
	FormatText
)

// Formatter reformats log lines into a specified output format.
type Formatter struct {
	format  Format
	fields  []string
	include bool
}

// New creates a new Formatter. fields is a list of field names to include
// (include=true) or exclude (include=false). Pass nil fields to select all.
func New(format Format, fields []string, include bool) (*Formatter, error) {
	if format != FormatRaw && format != FormatJSON && format != FormatText {
		return nil, fmt.Errorf("lineformatter: unknown format %d", format)
	}
	return &Formatter{format: format, fields: fields, include: include}, nil
}

// Apply reformats the given log line according to the formatter's configuration.
// If the line cannot be parsed as JSON it is returned unchanged for FormatRaw,
// or wrapped in a single-field map for other formats.
func (f *Formatter) Apply(line string) string {
	if f.format == FormatRaw {
		return line
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(line), &data); err != nil {
		// Not JSON — wrap as raw message
		data = map[string]interface{}{"message": line}
	}

	filtered := f.filterFields(data)

	switch f.format {
	case FormatJSON:
		b, err := json.Marshal(filtered)
		if err != nil {
			return line
		}
		return string(b)
	case FormatText:
		parts := make([]string, 0, len(filtered))
		for k, v := range filtered {
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
		}
		return strings.Join(parts, " ")
	}
	return line
}

func (f *Formatter) filterFields(data map[string]interface{}) map[string]interface{} {
	if len(f.fields) == 0 {
		return data
	}
	set := make(map[string]struct{}, len(f.fields))
	for _, field := range f.fields {
		set[field] = struct{}{}
	}
	out := make(map[string]interface{})
	for k, v := range data {
		_, inSet := set[k]
		if (f.include && inSet) || (!f.include && !inSet) {
			out[k] = v
		}
	}
	return out
}
