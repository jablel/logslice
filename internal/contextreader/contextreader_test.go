package contextreader

import (
	"testing"
)

func TestNew_ValidParams(t *testing.T) {
	cr, err := New(2, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cr == nil {
		t.Fatal("expected non-nil ContextReader")
	}
}

func TestNew_NegativeBefore(t *testing.T) {
	_, err := New(-1, 0)
	if err == nil {
		t.Fatal("expected error for negative before")
	}
}

func TestNew_NegativeAfter(t *testing.T) {
	_, err := New(0, -1)
	if err == nil {
		t.Fatal("expected error for negative after")
	}
}

func TestFeed_NoContext(t *testing.T) {
	cr, _ := New(0, 0)
	lines := []string{"a", "b", "c"}
	matched := []bool{false, true, false}

	var got []string
	for i, l := range lines {
		got = append(got, cr.Feed(l, matched[i])...)
	}
	if len(got) != 1 || got[0] != "b" {
		t.Fatalf("expected [b], got %v", got)
	}
}

func TestFeed_BeforeContext(t *testing.T) {
	cr, _ := New(2, 0)
	input := []string{"x1", "x2", "x3", "MATCH", "x4"}
	matched := []bool{false, false, false, true, false}

	var got []string
	for i, l := range input {
		got = append(got, cr.Feed(l, matched[i])...)
	}
	// expect x2, x3, MATCH
	want := []string{"x2", "x3", "MATCH"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: want %q got %q", i, w, got[i])
		}
	}
}

func TestFeed_AfterContext(t *testing.T) {
	cr, _ := New(0, 2)
	input := []string{"MATCH", "a1", "a2", "a3"}
	matched := []bool{true, false, false, false}

	var got []string
	for i, l := range input {
		got = append(got, cr.Feed(l, matched[i])...)
	}
	want := []string{"MATCH", "a1", "a2"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: want %q got %q", i, w, got[i])
		}
	}
}

func TestFeed_BeforeAndAfterContext(t *testing.T) {
	cr, _ := New(1, 1)
	input := []string{"pre", "MATCH", "post", "other"}
	matched := []bool{false, true, false, false}

	var got []string
	for i, l := range input {
		got = append(got, cr.Feed(l, matched[i])...)
	}
	want := []string{"pre", "MATCH", "post"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i, w := range want {
		if got[i] != w {
			t.Errorf("index %d: want %q got %q", i, w, got[i])
		}
	}
}

func TestReset(t *testing.T) {
	cr, _ := New(2, 2)
	cr.Feed("buffered", false)
	cr.Feed("MATCH", true)
	cr.Reset()

	got := cr.Feed("solo", true)
	if len(got) != 1 || got[0] != "solo" {
		t.Fatalf("after reset expected [solo], got %v", got)
	}
}
