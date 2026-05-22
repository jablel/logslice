package limiter_test

import (
	"strings"
	"testing"

	"github.com/logslice/logslice/internal/limiter"
)

// simulatePipeline mimics the scanner → limiter → collect loop used in main.
func simulatePipeline(lines []string, max int) []string {
	lim := limiter.New(max)
	var out []string
	for _, line := range lines {
		if !lim.Keep() {
			break
		}
		out = append(out, line)
	}
	return out
}

func TestIntegration_LimitFirstN(t *testing.T) {
	input := strings.Split("a\nb\nc\nd\ne", "\n")
	got := simulatePipeline(input, 3)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d: %v", len(got), got)
	}
	if got[0] != "a" || got[1] != "b" || got[2] != "c" {
		t.Fatalf("unexpected lines: %v", got)
	}
}

func TestIntegration_LimitLargerThanInput(t *testing.T) {
	input := strings.Split("x\ny", "\n")
	got := simulatePipeline(input, 100)
	if len(got) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(got))
	}
}

func TestIntegration_LimitZeroMeansUnlimited(t *testing.T) {
	input := make([]string, 50)
	for i := range input {
		input[i] = "line"
	}
	got := simulatePipeline(input, 0)
	if len(got) != 50 {
		t.Fatalf("expected 50 lines, got %d", len(got))
	}
}
