// Package contextreader provides a line reader that captures N lines of
// context (before/after) around each matched log line, similar to grep -C.
package contextreader

import "errors"

// errInvalidContext is returned when before or after values are negative.
var errInvalidContext = errors.New("contextreader: before and after must be >= 0")

// ContextReader buffers a sliding window of lines and emits matched lines
// together with their surrounding context lines.
type ContextReader struct {
	before     int
	after      int
	buf        []string
	afterCount int
}

// New creates a ContextReader that will include `before` lines preceding a
// match and `after` lines following a match. Both values must be >= 0.
func New(before, after int) (*ContextReader, error) {
	if before < 0 || after < 0 {
		return nil, errInvalidContext
	}
	return &ContextReader{
		before: before,
		after:  after,
		buf:    make([]string, 0, before+1),
	}, nil
}

// Feed adds a line to the reader. matched indicates whether the line itself
// satisfies the caller's filter. Feed returns the set of lines that should be
// emitted (context + matched lines), which may be empty.
func (c *ContextReader) Feed(line string, matched bool) []string {
	var out []string

	if matched {
		// Emit buffered before-context then the matched line itself.
		out = append(out, c.buf...)
		out = append(out, line)
		c.afterCount = c.after
		// Clear before-buffer so context lines aren't re-emitted.
		c.buf = c.buf[:0]
	} else {
		if c.afterCount > 0 {
			out = append(out, line)
			c.afterCount--
			// Don't buffer after-context lines as before-context.
			return out
		}
		// Maintain sliding before-context window.
		if c.before > 0 {
			if len(c.buf) == c.before {
				c.buf = c.buf[1:]
			}
			c.buf = append(c.buf, line)
		}
	}

	return out
}

// Reset clears internal state so the reader can be reused.
func (c *ContextReader) Reset() {
	c.buf = c.buf[:0]
	c.afterCount = 0
}
