// Package offsettracker records the byte offsets of log lines that pass
// through the processing pipeline.
//
// # Overview
//
// When slicing large log files it is sometimes useful to know exactly where
// in the source file each matched line lives. offsettracker maintains a
// running byte cursor that callers advance after consuming each raw line,
// and records (lineNumber, byteOffset, line) triples for every line that
// should be tracked.
//
// # Usage
//
//	tr := offsettracker.New(true)
//	for scanner.Scan() {
//		raw := scanner.Text() + "\n"
//		if filter.Match(raw) {
//			tr.Record(lineNum, raw)
//		}
//		tr.Advance(len(raw))
//		lineNum++
//	}
//	for _, e := range tr.Entries() {
//		fmt.Printf("line %d @ byte %d\n", e.LineNumber, e.ByteOffset)
//	}
package offsettracker
