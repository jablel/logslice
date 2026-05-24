// Package stats provides lightweight counters and reporting for logslice
// processing runs.
//
// Usage:
//
//	s := stats.New()          // records start time
//
//	// ... process lines ...
//	s.LinesScanned++
//	s.BytesRead += int64(len(line))
//	if matched {
//		s.LinesMatched++
//	} else {
//		s.LinesSkipped++
//	}
//
//	s.Finish()                // records end time
//	s.Print(os.Stderr)        // prints summary
//
// The summary includes elapsed time, total lines scanned, lines matched,
// lines skipped, and bytes read. For example:
//
//	scanned 10000 lines (8192 matched, 1808 skipped) in 42ms, 512.3 KB read
//
// Stats is not safe for concurrent use; callers that update counters from
// multiple goroutines should use their own synchronisation.
package stats
