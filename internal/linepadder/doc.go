// Package linepadder provides a fixed-width line padder for log output
// formatting.
//
// Each line fed through Apply is either padded with a configurable fill
// character or truncated so that the result is exactly Width() runes long.
//
// Alignment options:
//
//	AlignLeft  – padding is appended to the right of the line.
//	AlignRight – padding is prepended to the left of the line.
//
// Typical usage:
//
//	p, err := linepadder.New(80, linepadder.AlignLeft, ' ')
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(p.Apply("hello"))
package linepadder
