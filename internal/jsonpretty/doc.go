// Package jsonpretty formats JSON log lines for human-readable output.
//
// Usage:
//
//	f, err := jsonpretty.New("  ")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(f.Apply(line))
//
// Non-JSON lines are passed through unchanged, making the formatter safe
// to insert into any pipeline that may contain mixed-format log lines.
//
// An empty indent string disables formatting entirely so the formatter can
// be wired in unconditionally and toggled via configuration.
package jsonpretty
