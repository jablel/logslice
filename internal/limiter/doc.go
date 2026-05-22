// Package limiter provides a simple line-count cap for log slicing pipelines.
//
// When processing large log files it is often useful to retrieve only the
// first N matching lines rather than every line in a time range.  Limiter
// sits between the filter and the output writer and stops forwarding lines
// once the requested maximum has been reached.
//
// Usage:
//
//	lim := limiter.New(100)        // allow at most 100 lines
//	for scanner.Scan() {
//		if !lim.Keep() {
//			break
//		}
//		writer.Write(scanner.Text())
//	}
package limiter
