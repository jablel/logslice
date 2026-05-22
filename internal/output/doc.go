// Package output provides utilities for writing filtered log lines
// to an output destination with optional formatting.
//
// Supported formats:
//
//   - FormatRaw      — write lines verbatim, one per output line.
//   - FormatNumbered — prefix each line with its original line number
//                      separated by a tab character.
//
// Example usage:
//
//	w := output.New(output.Options{
//		Format:      output.FormatNumbered,
//		Destination: os.Stdout,
//	})
//	defer w.Flush()
//
//	w.WriteLine(101, "2024-01-15T10:00:00Z INFO server started")
//	fmt.Printf("wrote %d lines\n", w.Count())
package output
