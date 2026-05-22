package limiter

// Limiter caps the number of lines passed through to a maximum count.
// Once the limit is reached, Keep returns false for all subsequent lines.
type Limiter struct {
	max   int
	count int
}

// New creates a Limiter that allows at most max lines through.
// If max is zero or negative, no limit is applied (Keep always returns true).
func New(max int) *Limiter {
	return &Limiter{max: max}
}

// Keep returns true if the line should be passed through, false once the
// limit has been reached. It is safe to call Keep after the limit is hit;
// it will continue to return false.
func (l *Limiter) Keep() bool {
	if l.max <= 0 {
		return true
	}
	if l.count >= l.max {
		return false
	}
	l.count++
	return true
}

// Count returns the number of lines that have been allowed through so far.
func (l *Limiter) Count() int {
	return l.count
}

// Reset resets the internal counter, allowing another max lines through.
func (l *Limiter) Reset() {
	l.count = 0
}

// Done reports whether the limit has been reached.
func (l *Limiter) Done() bool {
	if l.max <= 0 {
		return false
	}
	return l.count >= l.max
}
