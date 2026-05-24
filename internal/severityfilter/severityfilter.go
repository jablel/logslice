// Package severityfilter filters log lines by severity/level.
// It supports common severity labels (DEBUG, INFO, WARN, ERROR, FATAL)
// and allows keeping only lines at or above a minimum severity threshold.
package severityfilter

import (
	"fmt"
	"strings"
)

// Level represents a log severity level.
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelNames = map[string]Level{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"WARNING": WARN,
	"ERROR": ERROR,
	"ERR":   ERROR,
	"FATAL": FATAL,
	"CRIT":  FATAL,
}

// ErrUnknownLevel is returned when the severity string is not recognised.
var ErrUnknownLevel = fmt.Errorf("unknown severity level")

// Filter keeps log lines whose detected severity is >= the minimum level.
type Filter struct {
	min Level
}

// New creates a Filter that passes lines at or above minLevel.
// minLevel must be one of: DEBUG, INFO, WARN, WARNING, ERROR, ERR, FATAL, CRIT.
func New(minLevel string) (*Filter, error) {
	lvl, ok := levelNames[strings.ToUpper(strings.TrimSpace(minLevel))]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknownLevel, minLevel)
	}
	return &Filter{min: lvl}, nil
}

// Keep returns true when the line contains a severity token >= the minimum level.
// If no known severity token is found the line is kept (pass-through).
func (f *Filter) Keep(line string) bool {
	upper := strings.ToUpper(line)
	best := -1
	for token, lvl := range levelNames {
		if strings.Contains(upper, token) {
			if int(lvl) > best {
				best = int(lvl)
			}
		}
	}
	if best == -1 {
		return true // no level detected — pass through
	}
	return Level(best) >= f.min
}

// Min returns the minimum severity level configured for this filter.
func (f *Filter) Min() Level { return f.min }
