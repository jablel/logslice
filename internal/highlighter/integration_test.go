package highlighter_test

import (
	"strings"
	"testing"

	"logslice/internal/highlighter"
)

var sampleLines = []string{
	"2024-06-01T10:00:00Z INFO  server started",
	"2024-06-01T10:01:00Z WARN  high memory usage",
	"2024-06-01T10:02:00Z ERROR disk write failed",
	"2024-06-01T10:03:00Z INFO  request handled",
	"2024-06-01T10:04:00Z ERROR connection timeout",
}

func applyAll(h *highlighter.Highlighter, lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = h.Apply(l)
	}
	return out
}

func TestIntegration_HighlightErrorLines(t *testing.T) {
	h, err := highlighter.New([]string{"ERROR"}, "red")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	results := applyAll(h, sampleLines)

	for _, r := range results {
		if strings.Contains(r, "ERROR") {
			// must be wrapped with ANSI codes
			if !strings.Contains(r, "\033[") {
				t.Errorf("ERROR line not highlighted: %q", r)
			}
		} else {
			// non-ERROR lines must not gain escape codes
			if strings.Contains(r, "\033[") {
				t.Errorf("non-ERROR line was unexpectedly modified: %q", r)
			}
		}
	}
}

func TestIntegration_MultiKeyword_WARNandERROR(t *testing.T) {
	h, err := highlighter.New([]string{"WARN", "ERROR"}, "yellow")
	if err != nil {
		t.Fatalf("setup: %v", err)
	}

	results := applyAll(h, sampleLines)
	highlightedCount := 0
	for _, r := range results {
		if strings.Contains(r, "\033[") {
			highlightedCount++
		}
	}

	// 1 WARN + 2 ERROR lines = 3
	if highlightedCount != 3 {
		t.Errorf("expected 3 highlighted lines, got %d", highlightedCount)
	}
}
