// Package tailreader provides a reader that tracks and resumes reading
// from a file starting at a given byte offset, enabling tail-like behaviour
// for incremental log processing.
package tailreader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Reader reads lines from a file starting at a specified byte offset.
type Reader struct {
	path   string
	offset int64
	file   *os.File
	scanner *bufio.Scanner
}

// New opens the file at path and seeks to startOffset before reading.
// Returns an error if the file cannot be opened or seeked.
func New(path string, startOffset int64) (*Reader, error) {
	if path == "" {
		return nil, fmt.Errorf("tailreader: path must not be empty")
	}
	if startOffset < 0 {
		return nil, fmt.Errorf("tailreader: startOffset must be >= 0, got %d", startOffset)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("tailreader: open %q: %w", path, err)
	}

	if startOffset > 0 {
		if _, err := f.Seek(startOffset, io.SeekStart); err != nil {
			f.Close()
			return nil, fmt.Errorf("tailreader: seek %q to %d: %w", path, startOffset, err)
		}
	}

	return &Reader{
		path:    path,
		offset:  startOffset,
		file:    f,
		scanner: bufio.NewScanner(f),
	}, nil
}

// ReadLine returns the next line (without trailing newline) and the new byte
// offset after reading it. Returns "", offset, io.EOF when no more lines exist.
func (r *Reader) ReadLine() (string, int64, error) {
	if !r.scanner.Scan() {
		if err := r.scanner.Err(); err != nil {
			return "", r.offset, fmt.Errorf("tailreader: scan: %w", err)
		}
		return "", r.offset, io.EOF
	}
	line := r.scanner.Text()
	r.offset += int64(len(line)) + 1 // +1 for newline
	return line, r.offset, nil
}

// Offset returns the current byte offset within the file.
func (r *Reader) Offset() int64 {
	return r.offset
}

// Close releases the underlying file handle.
func (r *Reader) Close() error {
	return r.file.Close()
}
