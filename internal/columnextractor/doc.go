// Package columnextractor provides a positional field extractor for
// delimiter-separated log lines.
//
// It is designed to work alongside timestampextractor and fieldextractor when
// log lines use a fixed columnar format (e.g. syslog, Apache access logs,
// or any whitespace-delimited output) rather than JSON or logfmt.
//
// Usage:
//
//	ex, err := columnextractor.New(" ", 2)
//	if err != nil { ... }
//	value, ok := ex.Extract("2024-01-15 12:00:00 INFO starting server")
//	// value == "INFO", ok == true
//
// The column index is zero-based. If the line contains fewer columns than
// the requested index, Extract returns ("", false).
package columnextractor
