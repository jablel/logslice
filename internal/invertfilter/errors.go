package invertfilter

import "errors"

// ErrNilKeeper is returned when New is called with a nil Keeper.
var ErrNilKeeper = errors.New("invertfilter: inner keeper must not be nil")
