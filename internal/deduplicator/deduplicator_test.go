package deduplicator

import (
	"testing"
)

func TestNew_ValidWindow(t *testing.T) {
	d, err := New(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Deduplicator")
	}
}

func TestNew_InvalidWindow(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for window=0")
	}
	_, err = New(-3)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestKeep_NoDuplicates(t *testing.T) {
	d, _ := New(3)
	lines := []string{"a", "b", "c", "d"}
	for _, l := range lines {
		if !d.Keep(l) {
			t.Errorf("expected Keep(%q)=true", l)
		}
	}
	if d.Skipped() != 0 {
		t.Errorf("expected 0 skipped, got %d", d.Skipped())
	}
}

func TestKeep_ConsecutiveDuplicate(t *testing.T) {
	d, _ := New(1)
	if !d.Keep("hello") {
		t.Fatal("first occurrence should be kept")
	}
	if d.Keep("hello") {
		t.Fatal("consecutive duplicate should be dropped")
	}
	if d.Skipped() != 1 {
		t.Errorf("expected 1 skipped, got %d", d.Skipped())
	}
}

func TestKeep_WindowedDuplicate(t *testing.T) {
	d, _ := New(3)
	d.Keep("x")
	d.Keep("y")
	if d.Keep("x") {
		t.Error("'x' is within window, should be dropped")
	}
}

func TestKeep_OutsideWindow(t *testing.T) {
	d, _ := New(2)
	d.Keep("x")
	d.Keep("y")
	d.Keep("z") // 'x' is now outside the window of 2
	if !d.Keep("x") {
		t.Error("'x' is outside window, should be kept")
	}
}

func TestReset(t *testing.T) {
	d, _ := New(3)
	d.Keep("a")
	d.Keep("b")
	d.Reset()
	if d.Skipped() != 0 {
		t.Errorf("expected 0 skipped after reset, got %d", d.Skipped())
	}
	if !d.Keep("a") {
		t.Error("after reset 'a' should be treated as new")
	}
}
