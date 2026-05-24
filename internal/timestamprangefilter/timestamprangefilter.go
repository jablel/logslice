// Package timestamprangefilter filters log lines by extracting timestamps
// and checking whether they fall within a specified time range.
package timestamprangefilter

import (
	"fmt"
	"time"
)

// Extractor extracts a timestamp from a raw log line.
type Extractor interface {
	Extract(line string) (time.Time, bool)
}

// Filter keeps only lines whose extracted timestamp falls within [From, To].
// Lines with no extractable timestamp are passed through unconditionally.
type Filter struct {
	extractor Extractor
	from      time.Time
	to        time.Time
	passNoTS  bool
}

// Option configures a Filter.
type Option func(*Filter)

// PassThroughNoTimestamp controls whether lines without a detectable
// timestamp are kept (true) or dropped (false). Default: keep them.
func PassThroughNoTimestamp(pass bool) Option {
	return func(f *Filter) { f.passNoTS = pass }
}

// New creates a Filter that retains lines whose timestamp is in [from, to].
// An empty (zero) to means "no upper bound".
func New(extractor Extractor, from, to time.Time, opts ...Option) (*Filter, error) {
	if extractor == nil {
		return nil, fmt.Errorf("timestamprangefilter: extractor must not be nil")
	}
	if !to.IsZero() && to.Before(from) {
		return nil, fmt.Errorf("timestamprangefilter: 'to' (%s) is before 'from' (%s)", to, from)
	}
	f := &Filter{
		extractor: extractor,
		from:      from,
		to:        to,
		passNoTS:  true,
	}
	for _, o := range opts {
		o(f)
	}
	return f, nil
}

// Keep returns true when the line should be included in the output.
func (f *Filter) Keep(line string) bool {
	ts, ok := f.extractor.Extract(line)
	if !ok {
		return f.passNoTS
	}
	if ts.Before(f.from) {
		return false
	}
	if !f.to.IsZero() && ts.After(f.to) {
		return false
	}
	return true
}
