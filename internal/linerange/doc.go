// Package linerange implements line-number range filtering for logslice.
//
// It allows callers to extract a contiguous slice of lines from a log stream
// by specifying a 1-based [from, to] interval. When the upper bound is reached
// the filter signals Done() so the scanner can stop reading early.
//
// Typical usage:
//
//	f, err := linerange.New(100, 200)
//	if err != nil { ... }
//	for scanner.Scan() {
//	    line := scanner.Text()
//	    if f.Keep(line) {
//	        fmt.Println(line)
//	    }
//	    if f.Done() { break }
//	}
package linerange
