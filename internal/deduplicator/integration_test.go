package deduplicator_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/deduplicator"
)

func simulateStream(lines []string, window int) ([]string, error) {
	d, err := deduplicator.New(window)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, l := range lines {
		if d.Keep(l) {
			out = append(out, l)
		}
	}
	return out, nil
}

func TestIntegration_NoDuplicatesInStream(t *testing.T) {
	input := []string{"a", "b", "c", "d"}
	got, err := simulateStream(input, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(input) {
		t.Errorf("expected %d lines, got %d", len(input), len(got))
	}
}

func TestIntegration_AllDuplicates(t *testing.T) {
	input := []string{"x", "x", "x", "x"}
	got, err := simulateStream(input, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 unique line, got %d", len(got))
	}
}

func TestIntegration_WindowEviction(t *testing.T) {
	// window=2: "a" is evicted after 2 unique lines, so second "a" should pass
	input := []string{"a", "b", "c", "a"}
	got, err := simulateStream(input, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Expect: a, b, c, a (second a re-admitted after eviction)
	if len(got) != 4 {
		t.Errorf("expected 4 lines after window eviction, got %d: %v", len(got), got)
	}
}

func TestIntegration_MixedBurstAndUnique(t *testing.T) {
	input := []string{"err", "err", "err", "info", "debug", "err"}
	got, err := simulateStream(input, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// err(1), info, debug, err suppressed (still in window)
	expected := []string{"err", "info", "debug"}
	if len(got) != len(expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}
