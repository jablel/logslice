// Package sampler implements periodic line sampling for logslice.
//
// When dealing with extremely high-frequency log files, even a filtered
// time-range slice can produce thousands of lines. The sampler reduces
// output by retaining only every Nth matched line, giving a representative
// view without overwhelming the consumer.
//
// Usage:
//
//	s, err := sampler.New(10)   // keep every 10th line
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, line := range matchedLines {
//		if s.Keep() {
//			fmt.Println(line)
//		}
//	}
//
// A step value of 1 is a no-op and keeps every line.
package sampler
