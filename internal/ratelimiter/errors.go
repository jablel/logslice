package ratelimiter

import "errors"

// ErrInvalidRate is returned when a negative rate is provided to New.
var ErrInvalidRate = errors.New("ratelimiter: rate must be zero or positive")
