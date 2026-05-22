package highlighter

import (
	"fmt"
	"strings"
)

// ANSI color codes.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// Highlighter marks occurrences of keywords within log lines using ANSI
// escape codes. It is safe to use with a nil keyword list (no-op).
type Highlighter struct {
	keywords []string
	color    string
	enabled  bool
}

// New creates a Highlighter for the given keywords.
// color must be one of: red, yellow, cyan, bold.
// If keywords is empty or color is unrecognised, highlighting is disabled.
func New(keywords []string, color string) (*Highlighter, error) {
	if len(keywords) == 0 {
		return &Highlighter{enabled: false}, nil
	}

	ansi, ok := map[string]string{
		"red":    colorRed,
		"yellow": colorYellow,
		"cyan":   colorCyan,
		"bold":   colorBold,
	}[strings.ToLower(color)]
	if !ok {
		return nil, fmt.Errorf("highlighter: unknown color %q (want red|yellow|cyan|bold)", color)
	}

	return &Highlighter{
		keywords: keywords,
		color:    ansi,
		enabled:  true,
	}, nil
}

// Apply returns the line with all keyword occurrences wrapped in ANSI codes.
// If highlighting is disabled the original line is returned unchanged.
func (h *Highlighter) Apply(line string) string {
	if !h.enabled {
		return line
	}
	for _, kw := range h.keywords {
		if kw == "" {
			continue
		}
		replacement := h.color + kw + colorReset
		line = strings.ReplaceAll(line, kw, replacement)
	}
	return line
}

// Enabled reports whether the highlighter will modify lines.
func (h *Highlighter) Enabled() bool { return h.enabled }
