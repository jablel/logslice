package lineannotator

import (
	"strings"
	"testing"
)

func TestNew_EmptyTag(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty tag, got nil")
	}
	if err != ErrEmptyTag {
		t.Fatalf("expected ErrEmptyTag, got %v", err)
	}
}

func TestNew_ValidTag(t *testing.T) {
	a, err := New("web-01")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Tag() != "web-01" {
		t.Fatalf("expected tag 'web-01', got %q", a.Tag())
	}
}

func TestApply_PrependsTag(t *testing.T) {
	a, _ := New("srv")
	got := a.Apply("hello world")
	if got != "[srv] hello world" {
		t.Fatalf("expected '[srv] hello world', got %q", got)
	}
}

func TestApply_EmptyLine_Unchanged(t *testing.T) {
	a, _ := New("srv")
	got := a.Apply("")
	if got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestApply_PrefixFormat(t *testing.T) {
	a, _ := New("my-service")
	line := "2024-01-01T00:00:00Z INFO started"
	got := a.Apply(line)
	if !strings.HasPrefix(got, "[my-service] ") {
		t.Fatalf("line does not start with expected prefix: %q", got)
	}
	if !strings.HasSuffix(got, line) {
		t.Fatalf("original line not preserved as suffix: %q", got)
	}
}

func TestApply_MultipleLines(t *testing.T) {
	a, _ := New("app")
	lines := []string{
		"line one",
		"line two",
		"line three",
	}
	for _, l := range lines {
		got := a.Apply(l)
		want := "[app] " + l
		if got != want {
			t.Errorf("Apply(%q) = %q, want %q", l, got, want)
		}
	}
}

func TestTag_ReturnsRawTag(t *testing.T) {
	tag := "production"
	a, _ := New(tag)
	if a.Tag() != tag {
		t.Fatalf("Tag() = %q, want %q", a.Tag(), tag)
	}
}
