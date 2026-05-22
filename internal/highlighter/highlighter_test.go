package highlighter

import (
	"strings"
	"testing"
)

func TestNew_NoKeywords(t *testing.T) {
	h, err := New(nil, "red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.Enabled() {
		t.Error("expected highlighter to be disabled when no keywords given")
	}
}

func TestNew_UnknownColor(t *testing.T) {
	_, err := New([]string{"ERROR"}, "purple")
	if err == nil {
		t.Fatal("expected error for unknown color")
	}
	if !strings.Contains(err.Error(), "purple") {
		t.Errorf("error should mention the bad color, got: %v", err)
	}
}

func TestNew_ValidColors(t *testing.T) {
	colors := []string{"red", "yellow", "cyan", "bold", "RED", "Yellow"}
	for _, c := range colors {
		_, err := New([]string{"kw"}, c)
		if err != nil {
			t.Errorf("color %q should be valid, got: %v", c, err)
		}
	}
}

func TestApply_NoOp_WhenDisabled(t *testing.T) {
	h, _ := New(nil, "red")
	line := "some ERROR line"
	if got := h.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_HighlightsKeyword(t *testing.T) {
	h, _ := New([]string{"ERROR"}, "red")
	line := "2024-01-01 ERROR something failed"
	result := h.Apply(line)
	if !strings.Contains(result, colorRed+"ERROR"+colorReset) {
		t.Errorf("expected ANSI-wrapped keyword in output, got: %q", result)
	}
	if strings.Count(result, "ERROR") != 1 {
		t.Error("keyword should appear exactly once (wrapped)")
	}
}

func TestApply_MultipleKeywords(t *testing.T) {
	h, _ := New([]string{"ERROR", "WARN"}, "yellow")
	line := "ERROR and WARN in one line"
	result := h.Apply(line)
	if !strings.Contains(result, colorYellow+"ERROR"+colorReset) {
		t.Error("ERROR not highlighted")
	}
	if !strings.Contains(result, colorYellow+"WARN"+colorReset) {
		t.Error("WARN not highlighted")
	}
}

func TestApply_EmptyKeywordSkipped(t *testing.T) {
	h, _ := New([]string{"", "ERROR"}, "cyan")
	line := "ERROR occurred"
	result := h.Apply(line)
	if !strings.Contains(result, "ERROR") {
		t.Error("expected ERROR to appear in output")
	}
}

func TestApply_LineWithoutKeyword(t *testing.T) {
	h, _ := New([]string{"ERROR"}, "bold")
	line := "INFO everything is fine"
	if got := h.Apply(line); got != line {
		t.Errorf("line without keyword should be unchanged, got %q", got)
	}
}
