package linebuffer

import (
	"testing"
)

func TestNew_ValidCapacity(t *testing.T) {
	b, err := New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Cap() != 10 {
		t.Fatalf("expected cap 10, got %d", b.Cap())
	}
	if b.Len() != 0 {
		t.Fatalf("expected len 0, got %d", b.Len())
	}
}

func TestNew_InvalidCapacity(t *testing.T) {
	for _, c := range []int{0, -1, -100} {
		_, err := New(c)
		if err != ErrInvalidCapacity {
			t.Errorf("cap=%d: expected ErrInvalidCapacity, got %v", c, err)
		}
	}
}

func TestPush_BelowCapacity(t *testing.T) {
	b, _ := New(5)
	b.Push("a")
	b.Push("b")
	b.Push("c")
	if b.Len() != 3 {
		t.Fatalf("expected len 3, got %d", b.Len())
	}
	lines := b.Lines()
	if lines[0] != "a" || lines[1] != "b" || lines[2] != "c" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

func TestPush_EvictsOldestWhenFull(t *testing.T) {
	b, _ := New(3)
	for _, l := range []string{"1", "2", "3", "4", "5"} {
		b.Push(l)
	}
	lines := b.Lines()
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[0] != "3" || lines[1] != "4" || lines[2] != "5" {
		t.Fatalf("unexpected lines after eviction: %v", lines)
	}
}

func TestLines_OrderPreserved(t *testing.T) {
	b, _ := New(4)
	input := []string{"alpha", "beta", "gamma", "delta"}
	for _, l := range input {
		b.Push(l)
	}
	got := b.Lines()
	for i, want := range input {
		if got[i] != want {
			t.Errorf("index %d: want %q, got %q", i, want, got[i])
		}
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	b, _ := New(5)
	b.Push("x")
	b.Push("y")
	b.Reset()
	if b.Len() != 0 {
		t.Fatalf("expected len 0 after reset, got %d", b.Len())
	}
	if lines := b.Lines(); len(lines) != 0 {
		t.Fatalf("expected no lines after reset, got %v", lines)
	}
}

func TestReset_AllowsReuse(t *testing.T) {
	b, _ := New(3)
	b.Push("a")
	b.Push("b")
	b.Reset()
	b.Push("c")
	lines := b.Lines()
	if len(lines) != 1 || lines[0] != "c" {
		t.Fatalf("expected [c] after reuse, got %v", lines)
	}
}

func TestCapacityOne(t *testing.T) {
	b, _ := New(1)
	b.Push("first")
	b.Push("second")
	lines := b.Lines()
	if len(lines) != 1 || lines[0] != "second" {
		t.Fatalf("expected [second], got %v", lines)
	}
}
