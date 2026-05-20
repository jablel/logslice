package timeparser

import (
	"fmt"
	"time"
)

// Common log timestamp formats to attempt parsing.
var knownFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05.000",
	"2006-01-02 15:04:05.000000",
	"02/Jan/2006:15:04:05 -0700",
	"Jan 02 15:04:05",
}

// Parse attempts to parse a timestamp string using a list of known formats.
// It returns the parsed time and the format that succeeded, or an error.
func Parse(s string) (time.Time, string, error) {
	for _, fmt := range knownFormats {
		t, err := time.Parse(fmt, s)
		if err == nil {
			return t, fmt, nil
		}
	}
	return time.Time{}, "", fmt.Errorf("timeparser: unable to parse timestamp %q", s)
}

// ParseWithFormat parses a timestamp using a specific format string.
func ParseWithFormat(s, layout string) (time.Time, error) {
	t, err := time.Parse(layout, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("timeparser: parse with format %q failed: %w", layout, err)
	}
	return t, nil
}

// InRange reports whether t is within [start, end] inclusive.
func InRange(t, start, end time.Time) bool {
	return !t.Before(start) && !t.After(end)
}

// KnownFormats returns a copy of the supported timestamp formats.
func KnownFormats() []string {
	copy := make([]string, len(knownFormats))
	for i, f := range knownFormats {
		copy[i] = f
	}
	return copy
}
