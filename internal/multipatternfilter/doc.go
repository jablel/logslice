// Package multipatternfilter provides a line filter that evaluates multiple
// regular-expression patterns combined with AND or OR logic.
//
// # Usage
//
//	f, err := multipatternfilter.New([]string{"ERROR", "timeout"}, multipatternfilter.ModeAnd)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, line := range lines {
//		if f.Keep(line) {
//			fmt.Println(line)
//		}
//	}
//
// ModeAnd keeps only lines that match every pattern.
// ModeOr keeps lines that match at least one pattern.
package multipatternfilter
