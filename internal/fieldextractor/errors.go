package fieldextractor

import "errors"

// ErrEmptyField is returned when an empty field name is provided.
var ErrEmptyField = errors.New("fieldextractor: field name must not be empty")
