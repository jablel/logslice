// Package invertfilter provides a decorator that negates any Keeper predicate.
//
// It is useful when a user wants to exclude lines that match a pattern rather
// than include them — for example, dropping all lines containing "DEBUG" while
// keeping everything else.
//
// Usage:
//
//	grep, _ := grepfilter.New("DEBUG", false)
//	inv, _ := invertfilter.New(grep, true)
//
//	for _, line := range lines {
//		if inv.Keep(line) {
//			// line does NOT contain "DEBUG"
//		}
//	}
package invertfilter
