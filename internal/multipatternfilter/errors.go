package multipatternfilter

import "errors"

// ErrNoPatterns is returned when no patterns are provided to New.
var ErrNoPatterns = errors.New("multipatternfilter: at least one pattern is required")
