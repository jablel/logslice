// Package filter provides time-range filtering logic for log lines.
//
// It integrates with the timeparser package to parse timestamps from
// log lines and determine whether they fall within a specified time window.
//
// Basic usage:
//
//	f, err := filter.New("2024-01-15T10:00:00Z", "2024-01-15T12:00:00Z", "")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	match, err := f.Match(logLine)
//	if err != nil {
//		// line has no parseable timestamp
//	}
//	if match {
//		// line is within the time range
//	}
//
// The format parameter is optional. When empty, auto-detection is used
// via timeparser.Parse. When provided, it must be a valid Go time layout
// string (e.g. "2006-01-02 15:04:05").
package filter
