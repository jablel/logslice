// Package highlighter provides keyword-based ANSI colour highlighting for
// log lines emitted by logslice.
//
// Usage:
//
//	h, err := highlighter.New([]string{"ERROR", "WARN"}, "red")
//	if err != nil { ... }
//	fmt.Println(h.Apply(line))
//
// Supported colours: red, yellow, cyan, bold.
// When no keywords are supplied the Apply method is a zero-allocation no-op.
//
// Colour codes:
//
//	red    – ANSI 31 (bright red foreground)
//	yellow – ANSI 33 (yellow foreground)
//	cyan   – ANSI 36 (cyan foreground)
//	bold   – ANSI 1  (bold / increased intensity)
//
// All sequences are reset with ANSI 0 after each highlighted keyword so that
// surrounding text is unaffected.
package highlighter
