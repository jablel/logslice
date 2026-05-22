// Package output handles writing filtered log lines to various destinations.
package output

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Format represents the output format for log lines.
type Format int

const (
	// FormatRaw writes lines as-is.
	FormatRaw Format = iota
	// FormatNumbered prefixes each line with its original line number.
	FormatNumbered
)

// Options configures the Writer behavior.
type Options struct {
	Format     Format
	Destination io.Writer
}

// Writer writes matched log lines to an output destination.
type Writer struct {
	opts Options
	bw   *bufio.Writer
	count int
}

// New creates a new Writer with the given options.
// If opts.Destination is nil, os.Stdout is used.
func New(opts Options) *Writer {
	dst := opts.Destination
	if dst == nil {
		dst = os.Stdout
	}
	return &Writer{
		opts: opts,
		bw:   bufio.NewWriter(dst),
	}
}

// WriteLine writes a single log line, applying the configured format.
func (w *Writer) WriteLine(lineNum int, line string) error {
	w.count++
	var err error
	switch w.opts.Format {
	case FormatNumbered:
		_, err = fmt.Fprintf(w.bw, "%d\t%s\n", lineNum, line)
	default:
		_, err = fmt.Fprintln(w.bw, line)
	}
	return err
}

// WriteLines writes multiple log lines, stopping on the first error.
// It returns the number of lines successfully written and any error encountered.
func (w *Writer) WriteLines(lines map[int]string) (int, error) {
	written := 0
	for lineNum, line := range lines {
		if err := w.WriteLine(lineNum, line); err != nil {
			return written, fmt.Errorf("writing line %d: %w", lineNum, err)
		}
		written++
	}
	return written, nil
}

// Flush flushes any buffered output.
func (w *Writer) Flush() error {
	return w.bw.Flush()
}

// Count returns the number of lines written.
func (w *Writer) Count() int {
	return w.count
}
