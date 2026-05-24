package linesampler

import (
	"fmt"
	"testing"
)

func TestNew_ValidRate(t *testing.T) {
	for _, rate := range []float64{0.0, 0.1, 0.5, 0.99, 1.0} {
		s, err := New(rate)
		if err != nil {
			t.Fatalf("New(%v) unexpected error: %v", rate, err)
		}
		if s.Rate() != rate {
			t.Errorf("Rate() = %v, want %v", s.Rate(), rate)
		}
	}
}

func TestNew_InvalidRate(t *testing.T) {
	for _, rate := range []float64{-0.1, 1.1, -1.0, 2.0} {
		_, err := New(rate)
		if err == nil {
			t.Errorf("New(%v) expected error, got nil", rate)
		}
	}
}

func TestKeep_RateOne_KeepsAll(t *testing.T) {
	s, _ := New(1.0)
	for i := 0; i < 100; i++ {
		if !s.Keep(fmt.Sprintf("line %d", i)) {
			t.Fatalf("rate=1.0 should keep all lines, dropped line %d", i)
		}
	}
}

func TestKeep_RateZero_DropsAll(t *testing.T) {
	s, _ := New(0.0)
	for i := 0; i < 100; i++ {
		if s.Keep(fmt.Sprintf("line %d", i)) {
			t.Fatalf("rate=0.0 should drop all lines, kept line %d", i)
		}
	}
}

func TestKeep_Deterministic(t *testing.T) {
	s, _ := New(0.5)
	line := "2024-01-01T00:00:00Z INFO some log message"
	first := s.Keep(line)
	for i := 0; i < 10; i++ {
		if s.Keep(line) != first {
			t.Error("Keep() is not deterministic for the same input")
		}
	}
}

func TestKeep_ApproximateRate(t *testing.T) {
	const n = 10000
	const targetRate = 0.3
	const tolerance = 0.05

	s, _ := New(targetRate)
	kept := 0
	for i := 0; i < n; i++ {
		if s.Keep(fmt.Sprintf("log line number %d with some content", i)) {
			kept++
		}
	}
	actual := float64(kept) / float64(n)
	if actual < targetRate-tolerance || actual > targetRate+tolerance {
		t.Errorf("keep rate %.3f outside expected range [%.3f, %.3f]",
			actual, targetRate-tolerance, targetRate+tolerance)
	}
}

func TestKeep_DifferentLinesGiveDifferentResults(t *testing.T) {
	s, _ := New(0.5)
	kept, dropped := 0, 0
	for i := 0; i < 200; i++ {
		if s.Keep(fmt.Sprintf("unique-line-%d", i)) {
			kept++
		} else {
			dropped++
		}
	}
	if kept == 0 || dropped == 0 {
		t.Errorf("expected mix of kept/dropped, got kept=%d dropped=%d", kept, dropped)
	}
}
