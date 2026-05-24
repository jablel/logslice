// Package linesorter buffers log lines and emits them in lexicographic order.
//
// This is useful when log sources write out-of-order entries (e.g. merged
// streams from multiple goroutines) and a deterministic output order is
// required for downstream processing or display.
//
// Usage:
//
//	s, _ := linesorter.New(linesorter.Ascending)
//	for _, line := range lines {
//		s.Feed(line)
//	}
//	sorted := s.Flush()
//
// Order options:
//
//	linesorter.Ascending  – A → Z
//	linesorter.Descending – Z → A
package linesorter
