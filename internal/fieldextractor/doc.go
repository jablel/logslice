// Package fieldextractor provides a lightweight field extraction utility
// for structured log lines.
//
// Supported formats:
//
//	- JSON: lines starting with '{' are parsed as JSON objects.
//	- Logfmt: key=value or key="value" pairs separated by whitespace.
//
// Format detection is automatic when FormatUnknown is specified, or can
// be set explicitly to FormatJSON or FormatLogfmt.
//
// Example usage:
//
//	ex, _ := fieldextractor.New("level", fieldextractor.FormatUnknown)
//	level := ex.Extract(`{"time":"2024-01-01","level":"error","msg":"oops"}`)
//	// level == "error"
package fieldextractor
