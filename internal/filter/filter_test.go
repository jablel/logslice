package filter

import (
	"testing"
)

func TestNew_AutoDetect(t *testing.T) {
	f, err := New("2024-01-15T10:00:00Z", "2024-01-15T12:00:00Z", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
	if f.From.IsZero() || f.To.IsZero() {
		t.Error("expected non-zero time bounds")
	}
}

func TestNew_WithFormat(t *testing.T) {
	f, err := New("2024-01-15 10:00:00", "2024-01-15 12:00:00", "2006-01-02 15:04:05")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Format != "2006-01-02 15:04:05" {
		t.Errorf("expected format to be preserved, got %q", f.Format)
	}
}

func TestNew_InvalidFrom(t *testing.T) {
	_, err := New("not-a-time", "2024-01-15T12:00:00Z", "")
	if err == nil {
		t.Fatal("expected error for invalid from timestamp")
	}
}

func TestNew_InvalidTo(t *testing.T) {
	_, err := New("2024-01-15T10:00:00Z", "not-a-time", "")
	if err == nil {
		t.Fatal("expected error for invalid to timestamp")
	}
}

func TestMatch_InRange(t *testing.T) {
	f, err := New("2024-01-15T10:00:00Z", "2024-01-15T12:00:00Z", "")
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	match, err := f.Match("2024-01-15T11:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Error("expected line to match")
	}
}

func TestMatch_OutOfRange(t *testing.T) {
	f, err := New("2024-01-15T10:00:00Z", "2024-01-15T12:00:00Z", "")
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	match, err := f.Match("2024-01-15T13:00:00Z")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if match {
		t.Error("expected line not to match")
	}
}

func TestMatch_InvalidLine(t *testing.T) {
	f, err := New("2024-01-15T10:00:00Z", "2024-01-15T12:00:00Z", "")
	if err != nil {
		t.Fatalf("setup error: %v", err)
	}

	_, err = f.Match("this line has no timestamp")
	if err == nil {
		t.Error("expected error for line without timestamp")
	}
}

func TestIsEmpty(t *testing.T) {
	f := &Filter{}
	if !f.IsEmpty() {
		t.Error("expected empty filter")
	}
}
