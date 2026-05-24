package linesorter

import (
	"testing"
)

func TestNew_Ascending(t *testing.T) {
	s, err := New(Ascending)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Sorter")
	}
}

func TestNew_Descending(t *testing.T) {
	s, err := New(Descending)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Sorter")
	}
}

func TestNew_InvalidOrder(t *testing.T) {
	_, err := New(Order(99))
	if err == nil {
		t.Fatal("expected error for unknown order")
	}
}

func TestFlush_Ascending(t *testing.T) {
	s, _ := New(Ascending)
	s.Feed("zebra")
	s.Feed("apple")
	s.Feed("mango")

	got := s.Flush()
	want := []string{"apple", "mango", "zebra"}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("index %d: got %q, want %q", i, got[i], v)
		}
	}
}

func TestFlush_Descending(t *testing.T) {
	s, _ := New(Descending)
	s.Feed("apple")
	s.Feed("zebra")
	s.Feed("mango")

	got := s.Flush()
	want := []string{"zebra", "mango", "apple"}
	for i, v := range want {
		if got[i] != v {
			t.Errorf("index %d: got %q, want %q", i, got[i], v)
		}
	}
}

func TestFlush_Empty(t *testing.T) {
	s, _ := New(Ascending)
	if got := s.Flush(); got != nil {
		t.Errorf("expected nil for empty flush, got %v", got)
	}
}

func TestFlush_ClearsBuffer(t *testing.T) {
	s, _ := New(Ascending)
	s.Feed("line1")
	s.Flush()
	if s.Len() != 0 {
		t.Errorf("expected buffer to be empty after flush, got %d", s.Len())
	}
}

func TestLen(t *testing.T) {
	s, _ := New(Ascending)
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
	s.Feed("a")
	s.Feed("b")
	if s.Len() != 2 {
		t.Errorf("expected 2, got %d", s.Len())
	}
}

func TestReset(t *testing.T) {
	s, _ := New(Ascending)
	s.Feed("x")
	s.Feed("y")
	s.Reset()
	if s.Len() != 0 {
		t.Errorf("expected 0 after reset, got %d", s.Len())
	}
}
