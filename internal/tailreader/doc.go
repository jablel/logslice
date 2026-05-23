// Package tailreader implements offset-aware log file reading.
//
// It opens a file at a given byte offset and exposes a line-by-line iterator
// that returns each line together with the updated byte offset. This allows
// callers to persist the offset (e.g. via offsettracker) and resume reading
// from exactly where they left off on subsequent runs — without re-processing
// already-seen log entries.
//
// Typical usage:
//
//	r, err := tailreader.New("/var/log/app.log", savedOffset)
//	if err != nil { ... }
//	defer r.Close()
//
//	for {
//		line, newOffset, err := r.ReadLine()
//		if err == io.EOF { break }
//		if err != nil { ... }
//		process(line)
//		savedOffset = newOffset
//	}
package tailreader
