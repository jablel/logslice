// Package headtailreader implements head/tail line extraction for log streams.
//
// # Overview
//
// headtailreader exposes a single Reader type that can operate in two modes:
//
//   - "head" — retains the first N lines fed to it, discarding the rest.
//   - "tail" — maintains a rolling window of the last N lines seen.
//
// # Usage
//
//	r, err := headtailreader.New("tail", 20)
//	if err != nil { ... }
//	for _, line := range lines {
//	    r.Feed(line)
//	}
//	fmt.Println(r.Lines())
//
// Both modes are allocation-efficient: the internal buffer is pre-allocated
// to the requested capacity and never grows beyond it.
package headtailreader
