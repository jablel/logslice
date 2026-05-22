// Package deduplicator provides a sliding-window deduplication filter
// for log lines. It removes repeated lines that appear within a
// configurable window of recently seen entries, helping to reduce
// noise from burst-repeated log messages.
//
// Usage:
//
//	d, err := deduplicator.New(50)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, line := range lines {
//		if d.Keep(line) {
//			fmt.Println(line)
//		}
//	}
//
// A window size of 1 only suppresses immediately consecutive duplicates.
// Larger windows catch duplicates that reappear within N unique lines.
package deduplicator
