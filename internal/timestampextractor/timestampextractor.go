// Package timestampextractor locates and extracts the timestamp substring
// from a raw log line by trying each known time format in order.
package timestampextractor

import (
	"errors"
	"strings"
	"time"

	"github.com/user/logslice/internal/timeparser"
)

// ErrNoTimestamp is returned when no known timestamp can be found in a line.
var ErrNoTimestamp = errors.New("timestampextractor: no timestamp found in line")

// Result holds the extracted timestamp and its position within the line.
type Result struct {
	Time   time.Time
	Raw    string // the matched substring
	Offset int    // byte offset of Raw within the line
	Format string // the format that matched
}

// Extractor tries to find a timestamp in log lines.
type Extractor struct {
	formats []string
}

// New creates an Extractor that will try the provided formats in order.
// If formats is empty, timeparser.KnownFormats is used.
func New(formats []string) *Extractor {
	if len(formats) == 0 {
		formats = timeparser.KnownFormats
	}
	return &Extractor{formats: formats}
}

// Extract scans line for a recognisable timestamp and returns a Result.
// It tries every registered format and returns the first successful match.
func (e *Extractor) Extract(line string) (Result, error) {
	for _, fmt := range e.formats {
		fmtLen := len(fmt)
		if fmtLen > len(line) {
			continue
		}
		// Slide a window of the format length across the line.
		for offset := 0; offset <= len(line)-fmtLen; offset++ {
			candidate := line[offset : offset+fmtLen]
			t, err := timeparser.ParseWithFormat(candidate, fmt)
			if err != nil {
				// Skip whitespace-leading positions quickly.
				if offset < len(line)-1 && line[offset] != ' ' && !strings.ContainsRune("[(", rune(line[offset])) {
					// advance to next space boundary
					next := strings.IndexAny(line[offset+1:], " \t[(")
					if next >= 0 {
						offset += next
					}
				}
				continue
			}
			return Result{
				Time:   t,
				Raw:    candidate,
				Offset: offset,
				Format: fmt,
			}, nil
		}
	}
	return Result{}, ErrNoTimestamp
}
