// Package linebuffer implements a fixed-capacity ring buffer for log lines.
//
// # Overview
//
// Buffer retains the last N lines pushed into it. When the buffer is full,
// the oldest line is silently evicted to make room for the new one. This
// makes it suitable for sliding-window log processing where only a recent
// window of lines is needed (e.g., printing trailing context on a match).
//
// # Usage
//
//	buf, err := linebuffer.New(100)
//	if err != nil { ... }
//
//	for _, line := range logLines {
//	    buf.Push(line)
//	}
//
//	for _, l := range buf.Lines() {
//	    fmt.Println(l)
//	}
package linebuffer
