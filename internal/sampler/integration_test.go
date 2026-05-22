package sampler_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/sampler"
)

// simulateFilter mimics the filter loop in cmd/logslice/main.go.
func simulateFilter(lines []string, step int) ([]string, error) {
	s, err := sampler.New(step)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, l := range lines {
		if s.Keep() {
			out = append(out, l)
		}
	}
	return out, nil
}

func TestIntegration_SampleHalfLines(t *testing.T) {
	lines := strings.Split(
		"line1\nline2\nline3\nline4\nline5\nline6", "\n",
	)
	got, err := simulateFilter(lines, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Errorf("expected 3 sampled lines, got %d: %v", len(got), got)
	}
	expected := []string{"line2", "line4", "line6"}
	for i, want := range expected {
		if got[i] != want {
			t.Errorf("index %d: expected %q, got %q", i, want, got[i])
		}
	}
}

func TestIntegration_StepLargerThanInput(t *testing.T) {
	lines := []string{"a", "b", "c"}
	got, err := simulateFilter(lines, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 lines, got %d", len(got))
	}
}
