package headtailreader

import (
	"testing"
)

func TestNew_InvalidCount(t *testing.T) {
	_, err := New("head", 0)
	if err == nil {
		t.Fatal("expected error for count=0")
	}
	_, err = New("tail", -1)
	if err == nil {
		t.Fatal("expected error for count=-1")
	}
}

func TestNew_InvalidMode(t *testing.T) {
	_, err := New("middle", 5)
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNew_Valid(t *testing.T) {
	for _, mode := range []string{"head", "tail"} {
		r, err := New(mode, 3)
		if err != nil {
			t.Fatalf("mode %q: unexpected error: %v", mode, err)
		}
		if r == nil {
			t.Fatalf("mode %q: expected non-nil reader", mode)
		}
	}
}

func TestHead_KeepsFirstN(t *testing.T) {
	r, _ := New("head", 3)
	for _, l := range []string{"a", "b", "c", "d", "e"} {
		r.Feed(l)
	}
	got := r.Lines()
	if len(got) != 3 {
		t.Fatalf("want 3 lines, got %d", len(got))
	}
	for i, want := range []string{"a", "b", "c"} {
		if got[i] != want {
			t.Errorf("line[%d]: want %q got %q", i, want, got[i])
		}
	}
}

func TestTail_KeepsLastN(t *testing.T) {
	r, _ := New("tail", 3)
	for _, l := range []string{"a", "b", "c", "d", "e"} {
		r.Feed(l)
	}
	got := r.Lines()
	if len(got) != 3 {
		t.Fatalf("want 3 lines, got %d", len(got))
	}
	for i, want := range []string{"c", "d", "e"} {
		if got[i] != want {
			t.Errorf("line[%d]: want %q got %q", i, want, got[i])
		}
	}
}

func TestSeen(t *testing.T) {
	r, _ := New("head", 2)
	for _, l := range []string{"x", "y", "z"} {
		r.Feed(l)
	}
	if r.Seen() != 3 {
		t.Errorf("want Seen()=3, got %d", r.Seen())
	}
}

func TestReset(t *testing.T) {
	r, _ := New("tail", 3)
	r.Feed("a")
	r.Feed("b")
	r.Reset()
	if r.Seen() != 0 {
		t.Errorf("want Seen()=0 after reset, got %d", r.Seen())
	}
	if len(r.Lines()) != 0 {
		t.Errorf("want empty Lines() after reset")
	}
}

func TestFewerLinesThanN(t *testing.T) {
	r, _ := New("head", 10)
	r.Feed("only")
	if len(r.Lines()) != 1 {
		t.Errorf("want 1 line, got %d", len(r.Lines()))
	}
}
