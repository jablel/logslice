package columnextractor

import (
	"fmt"
	"strings"
)

// Extractor splits each log line by a delimiter and returns the value at a
// specific column index (0-based). It is useful for space- or tab-separated
// log formats where a fixed positional field carries the timestamp or level.
type Extractor struct {
	delimiter string
	column    int
}

// New creates an Extractor that splits on delimiter and returns column index.
// Returns an error if delimiter is empty or column is negative.
func New(delimiter string, column int) (*Extractor, error) {
	if delimiter == "" {
		return nil, fmt.Errorf("columnextractor: delimiter must not be empty")
	}
	if column < 0 {
		return nil, fmt.Errorf("columnextractor: column index must be non-negative, got %d", column)
	}
	return &Extractor{delimiter: delimiter, column: column}, nil
}

// Extract splits line by the configured delimiter and returns the value at
// the configured column index. Returns an empty string and false when the
// line has fewer columns than required.
func (e *Extractor) Extract(line string) (string, bool) {
	if line == "" {
		return "", false
	}
	parts := strings.Split(line, e.delimiter)
	if e.column >= len(parts) {
		return "", false
	}
	value := strings.TrimSpace(parts[e.column])
	if value == "" {
		return "", false
	}
	return value, true
}

// ColumnCount returns the number of columns found in line when split by the
// configured delimiter.
func (e *Extractor) ColumnCount(line string) int {
	if line == "" {
		return 0
	}
	return len(strings.Split(line, e.delimiter))
}
