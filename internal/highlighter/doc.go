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
package highlighter
