package linetruncator

import (
	"strings"
	"testing"
)

func TestNew_ValidNoSuffix(t *testing.T) {
	tr, err := New(80, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Truncator")
	}
}

func TestNew_ValidWithSuffix(t *testing.T) {
	_, err := New(80, "...")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_InvalidMaxLen(t *testing.T) {
	_, err := New(0, "")
	if err != ErrInvalidMaxLen {
		t.Fatalf("expected ErrInvalidMaxLen, got %v", err)
	}
}

func TestNew_SuffixTooLong(t *testing.T) {
	_, err := New(3, "...")
	if err != ErrSuffixTooLong {
		t.Fatalf("expected ErrSuffixTooLong, got %v", err)
	}
}

func TestApply_ShortLine_Unchanged(t *testing.T) {
	tr, _ := New(20, "...")
	line := "hello world"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_ExactLength_Unchanged(t *testing.T) {
	tr, _ := New(5, "...")
	line := "hello"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_TruncatesWithSuffix(t *testing.T) {
	tr, _ := New(10, "...")
	line := "hello world this is long"
	got := tr.Apply(line)
	if len(got) != 10 {
		t.Errorf("expected length 10, got %d (%q)", len(got), got)
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected suffix '...', got %q", got)
	}
}

func TestApply_TruncatesNoSuffix(t *testing.T) {
	tr, _ := New(5, "")
	line := "abcdefghij"
	got := tr.Apply(line)
	if got != "abcde" {
		t.Errorf("expected %q, got %q", "abcde", got)
	}
}

func TestApply_EmptyLine(t *testing.T) {
	tr, _ := New(10, "...")
	if got := tr.Apply(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestApply_MultipleLines(t *testing.T) {
	tr, _ := New(8, ">>")
	cases := []struct {
		input string
		want  string
	}{
		{"short", "short"},
		{"exactly8", "exactly8"},
		{"toolongline", "toolong>>"},
	}
	for _, c := range cases {
		if got := tr.Apply(c.input); got != c.want {
			t.Errorf("Apply(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
