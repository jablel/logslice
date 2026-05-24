// Package linehasher annotates log lines with a short hash digest.
//
// Each line is hashed using the chosen algorithm (md5 or sha256) and a
// configurable number of hex characters from the digest are prepended or
// appended to the line. This is useful for:
//
//   - Generating stable deduplication keys from log content.
//   - Correlating the same log event across multiple files or streams.
//   - Lightweight integrity verification of pipeline output.
//
// Example usage:
//
//	h, err := linehasher.New(linehasher.SHA256, linehasher.Prepend, 8)
//	if err != nil { ... }
//	annotated := h.Apply("2024-01-15T10:00:00Z level=error msg=\"disk full\"")
//	// => "[a3f1b2c4] 2024-01-15T10:00:00Z level=error msg=\"disk full\""
package linehasher
