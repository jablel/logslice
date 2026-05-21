package filter

import (
	"time"

	"github.com/logslice/logslice/internal/timeparser"
)

// Filter holds the time range configuration for log filtering.
type Filter struct {
	From   time.Time
	To     time.Time
	Format string
}

// New creates a new Filter from string representations of from/to timestamps.
// If format is empty, it will attempt to auto-detect the format.
func New(from, to, format string) (*Filter, error) {
	var parsedFrom, parsedTo time.Time
	var err error

	if format != "" {
		parsedFrom, err = timeparser.ParseWithFormat(from, format)
		if err != nil {
			return nil, err
		}
		parsedTo, err = timeparser.ParseWithFormat(to, format)
		if err != nil {
			return nil, err
		}
	} else {
		parsedFrom, err = timeparser.Parse(from)
		if err != nil {
			return nil, err
		}
		parsedTo, err = timeparser.Parse(to)
		if err != nil {
			return nil, err
		}
	}

	return &Filter{
		From:   parsedFrom,
		To:     parsedTo,
		Format: format,
	}, nil
}

// Match checks whether a log line's timestamp falls within the filter range.
// It extracts the timestamp from the line using the configured or auto-detected format.
func (f *Filter) Match(line string) (bool, error) {
	var t time.Time
	var err error

	if f.Format != "" {
		t, err = timeparser.ParseWithFormat(line, f.Format)
	} else {
		t, err = timeparser.Parse(line)
	}
	if err != nil {
		return false, err
	}

	return timeparser.InRange(t, f.From, f.To), nil
}

// IsEmpty returns true if the filter has no meaningful time bounds set.
func (f *Filter) IsEmpty() bool {
	return f.From.IsZero() && f.To.IsZero()
}
