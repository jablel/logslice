package prefixstripper

import (
	"testing"
)

func TestNew_EmptyPrefix(t *testing.T) {
	_, err := New("", false)
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestNew_InvalidRegex(t *testing.T) {
	_, err := New("[invalid", true)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_ValidFixed(t *testing.T) {
	s, err := New("INFO ", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Stripper")
	}
}

func TestNew_ValidRegex(t *testing.T) {
	s, err := New(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z? `, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Stripper")
	}
}

func TestApply_Fixed_RemovesPrefix(t *testing.T) {
	s, _ := New("INFO ", false)
	got := s.Apply("INFO hello world")
	if got != "hello world" {
		t.Errorf("got %q, want %q", got, "hello world")
	}
}

func TestApply_Fixed_NoMatch_ReturnsOriginal(t *testing.T) {
	s, _ := New("INFO ", false)
	line := "DEBUG something"
	got := s.Apply(line)
	if got != line {
		t.Errorf("got %q, want %q", got, line)
	}
}

func TestApply_Regex_RemovesPrefix(t *testing.T) {
	s, _ := New(`\[\w+\] `, true)
	got := s.Apply("[ERROR] disk full")
	if got != "disk full" {
		t.Errorf("got %q, want %q", got, "disk full")
	}
}

func TestApply_Regex_NoMatch_ReturnsOriginal(t *testing.T) {
	s, _ := New(`\[\w+\] `, true)
	line := "no bracket prefix here"
	got := s.Apply(line)
	if got != line {
		t.Errorf("got %q, want %q", got, line)
	}
}

func TestApply_Regex_AutoAnchored(t *testing.T) {
	// Pattern without ^ should still only strip from the start.
	s, _ := New(`INFO `, true)
	got := s.Apply("INFO booting")
	if got != "booting" {
		t.Errorf("got %q, want %q", got, "booting")
	}
	// Prefix appearing in the middle must NOT be stripped.
	mid := "DEBUG INFO booting"
	got = s.Apply(mid)
	if got != mid {
		t.Errorf("middle match should not strip; got %q", got)
	}
}

func TestApply_EmptyLine(t *testing.T) {
	s, _ := New("INFO ", false)
	got := s.Apply("")
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}
