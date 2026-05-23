package linerange

import (
	"fmt"
	"testing"
)

func TestNew_ValidRange(t *testing.T) {
	f, err := New(1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestNew_UnboundedTo(t *testing.T) {
	_, err := New(5, 0)
	if err != nil {
		t.Fatalf("unexpected error for to=0: %v", err)
	}
}

func TestNew_InvalidFrom(t *testing.T) {
	_, err := New(0, 10)
	if err == nil {
		t.Fatal("expected error for from=0")
	}
}

func TestNew_ToLessThanFrom(t *testing.T) {
	_, err := New(10, 5)
	if err == nil {
		t.Fatal("expected error when to < from")
	}
}

func TestKeep_InRange(t *testing.T) {
	f, _ := New(3, 5)
	results := make([]bool, 7)
	for i := range results {
		results[i] = f.Keep(fmt.Sprintf("line %d", i+1))
	}
	// lines 3,4,5 (index 2,3,4) should be kept
	expected := []bool{false, false, true, true, true, false, false}
	for i, got := range results {
		if got != expected[i] {
			t.Errorf("line %d: got %v, want %v", i+1, got, expected[i])
		}
	}
}

func TestKeep_Unbounded(t *testing.T) {
	f, _ := New(2, 0)
	if f.Keep("line1") {
		t.Error("line 1 should not be kept")
	}
	for i := 2; i <= 100; i++ {
		if !f.Keep(fmt.Sprintf("line%d", i)) {
			t.Errorf("line %d should be kept with unbounded to", i)
		}
	}
}

func TestDone_WhenPastUpperBound(t *testing.T) {
	f, _ := New(1, 3)
	for i := 0; i < 3; i++ {
		f.Keep("x")
	}
	if !f.Done() {
		t.Error("expected Done() after reaching upper bound")
	}
}

func TestDone_Unbounded(t *testing.T) {
	f, _ := New(1, 0)
	for i := 0; i < 1000; i++ {
		f.Keep("x")
	}
	if f.Done() {
		t.Error("unbounded filter should never be Done")
	}
}

func TestReset(t *testing.T) {
	f, _ := New(1, 2)
	f.Keep("a")
	f.Keep("b")
	if !f.Done() {
		t.Fatal("expected Done before reset")
	}
	f.Reset()
	if f.Done() {
		t.Error("expected not Done after reset")
	}
	if f.Current() != 0 {
		t.Errorf("expected current=0 after reset, got %d", f.Current())
	}
}

func TestCurrent(t *testing.T) {
	f, _ := New(1, 10)
	for i := 0; i < 5; i++ {
		f.Keep("x")
	}
	if f.Current() != 5 {
		t.Errorf("expected current=5, got %d", f.Current())
	}
}
