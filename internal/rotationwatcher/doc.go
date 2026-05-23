// Package rotationwatcher monitors a log file for rotation events.
//
// Log rotation typically occurs when a log manager (e.g. logrotate) renames
// the current log file and creates a new one at the original path. This
// package detects such events by periodically comparing the file's inode
// number and size against a previously recorded snapshot.
//
// Usage:
//
//	w, err := rotationwatcher.New("/var/log/app.log", 500*time.Millisecond)
//	if err != nil { ... }
//	w.Start()
//	defer w.Stop()
//
//	// Later, check for rotation:
//	if w.Rotated() {
//		// Reopen the file, then:
//		w.Reset()
//	}
//
// On Windows, inode-based detection is unavailable; only size-shrinkage
// detection is used.
package rotationwatcher
