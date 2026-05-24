package lineskipper

import (
	"testing"
)

func TestNew_ZeroSkip(t *testing.T) {
	s, err := New(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Skipper")
	}
}

func TestNew_PositiveSkip(t *testing.T) {
	s, err := New(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.skip != 5 {
		t.Fatalf("expected skip=5, got %d", s.skip)
	}
}

func TestNew_NegativeSkip(t *testing.T) {
	_, err := New(-1)
	if err == nil {
		t.Fatal("expected error for negative skip")
	}
	if err != ErrNegativeSkip {
		t.Fatalf("expected ErrNegativeSkip, got %v", err)
	}
}

func TestKeep_SkipsFirstN(t *testing.T) {
	s, _ := New(3)
	lines := []string{"a", "b", "c", "d", "e"}
	expected := []bool{false, false, false, true, true}
	for i, line := range lines {
		got := s.Keep(line)
		if got != expected[i] {
			t.Errorf("line %d (%q): expected Keep=%v, got %v", i, line, expected[i], got)
		}
	}
}

func TestKeep_ZeroSkip_KeepsAll(t *testing.T) {
	s, _ := New(0)
	for i, line := range []string{"x", "y", "z"} {
		if !s.Keep(line) {
			t.Errorf("line %d: expected Keep=true with skip=0", i)
		}
	}
}

func TestSkipped_Count(t *testing.T) {
	s, _ := New(3)
	s.Keep("a")
	s.Keep("b")
	if got := s.Skipped(); got != 2 {
		t.Fatalf("expected Skipped()=2, got %d", got)
	}
	s.Keep("c")
	s.Keep("d") // past the skip window
	if got := s.Skipped(); got != 3 {
		t.Fatalf("expected Skipped()=3 after window, got %d", got)
	}
}

func TestReset(t *testing.T) {
	s, _ := New(2)
	s.Keep("a")
	s.Keep("b")
	s.Keep("c")
	s.Reset()
	if s.seen != 0 {
		t.Fatalf("expected seen=0 after Reset, got %d", s.seen)
	}
	// After reset the first two lines should be skipped again.
	if s.Keep("x") {
		t.Error("expected first line after Reset to be skipped")
	}
	if s.Keep("y") {
		t.Error("expected second line after Reset to be skipped")
	}
	if !s.Keep("z") {
		t.Error("expected third line after Reset to be kept")
	}
}
