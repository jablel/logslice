package linepadder_test

import (
	"testing"

	"logslice/internal/linepadder"
)

var logLines = []string{
	"INFO  starting server",
	"WARN  disk usage high",
	"ERROR connection refused",
	"DEBUG received 42 bytes",
}

func applyAll(p *linepadder.Padder, lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = p.Apply(l)
	}
	return out
}

func TestIntegration_FixedWidth_AllSameLength(t *testing.T) {
	p, err := linepadder.New(30, linepadder.AlignLeft, ' ')
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	results := applyAll(p, logLines)
	for i, r := range results {
		runes := []rune(r)
		if len(runes) != 30 {
			t.Errorf("line %d: expected length 30, got %d (%q)", i, len(runes), r)
		}
	}
}

func TestIntegration_RightAlign_PaddedPrefix(t *testing.T) {
	p, err := linepadder.New(30, linepadder.AlignRight, '.')
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	results := applyAll(p, logLines)
	for i, r := range results {
		runes := []rune(r)
		if len(runes) != 30 {
			t.Errorf("line %d: expected length 30, got %d", i, len(runes))
		}
		// original line must appear as suffix
		orig := logLines[i]
		if len([]rune(orig)) <= 30 && r[len(r)-len(orig):] != orig {
			t.Errorf("line %d: expected suffix %q in %q", i, orig, r)
		}
	}
}

func TestIntegration_Disabled_NoChange(t *testing.T) {
	p, _ := linepadder.New(30, linepadder.AlignLeft, ' ')
	p.SetEnabled(false)
	results := applyAll(p, logLines)
	for i, r := range results {
		if r != logLines[i] {
			t.Errorf("line %d: expected unchanged %q, got %q", i, logLines[i], r)
		}
	}
}
