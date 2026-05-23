// Package invertfilter wraps any line-keep predicate and negates its result,
// allowing users to exclude lines that match a given condition.
package invertfilter

// Keeper is the interface satisfied by any filter whose Keep method
// determines whether a line should be retained.
type Keeper interface {
	Keep(line string) bool
}

// InvertFilter negates the decision of an underlying Keeper.
type InvertFilter struct {
	inner   Keeper
	enabled bool
}

// New returns an InvertFilter that wraps inner.
// When enabled is false the filter is a no-op and always returns true.
func New(inner Keeper, enabled bool) (*InvertFilter, error) {
	if inner == nil {
		return nil, ErrNilKeeper
	}
	return &InvertFilter{inner: inner, enabled: enabled}, nil
}

// Keep returns true when the underlying keeper returns false (i.e. the line
// does NOT match), effectively inverting the filter logic.
// If the filter is disabled every line is kept.
func (f *InvertFilter) Keep(line string) bool {
	if !f.enabled {
		return true
	}
	return !f.inner.Keep(line)
}

// Enabled reports whether inversion is active.
func (f *InvertFilter) Enabled() bool { return f.enabled }
