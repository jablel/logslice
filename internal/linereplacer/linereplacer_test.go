package linereplacer

import (
	"testing"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New("", "x", false)
	if err == nil {
		t.Fatal("expected error for empty pattern, got nil")
	}
}

func TestNew_InvalidRegex(t *testing.T) {
	_, err := New("[invalid", "x", true)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestNew_ValidFixed(t *testing.T) {
	r, err := New("foo", "bar", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Replacer")
	}
}

func TestNew_ValidRegex(t *testing.T) {
	r, err := New(`\d+`, "NUM", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil Replacer")
	}
}

func TestApply_Fixed_ReplacesSubstring(t *testing.T) {
	r, _ := New("ERROR", "WARN", false)
	got := r.Apply("2024-01-01 ERROR something went wrong ERROR")
	want := "2024-01-01 WARN something went wrong WARN"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApply_Fixed_NoMatch_Unchanged(t *testing.T) {
	r, _ := New("ERROR", "WARN", false)
	line := "2024-01-01 INFO all good"
	if got := r.Apply(line); got != line {
		t.Errorf("got %q, want %q", got, line)
	}
}

func TestApply_Regex_ReplacesMatches(t *testing.T) {
	r, _ := New(`\d{4}-\d{2}-\d{2}`, "DATE", true)
	got := r.Apply("2024-01-01 INFO event on 2024-06-15")
	want := "DATE INFO event on DATE"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApply_Regex_CaptureGroup(t *testing.T) {
	r, _ := New(`user=(\w+)`, "user=REDACTED", true)
	got := r.Apply("login user=alice host=web1")
	want := "login user=REDACTED host=web1"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApply_EmptyLine(t *testing.T) {
	r, _ := New("foo", "bar", false)
	if got := r.Apply(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestApply_Regex_EmptyReplacement(t *testing.T) {
	r, _ := New(`\s+`, "", true)
	got := r.Apply("hello world  test")
	want := "helloworldtest"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
