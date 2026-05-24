package columnextractor

import (
	"testing"
)

func TestNew_EmptyDelimiter(t *testing.T) {
	_, err := New("", 0)
	if err == nil {
		t.Fatal("expected error for empty delimiter")
	}
}

func TestNew_NegativeColumn(t *testing.T) {
	_, err := New(" ", -1)
	if err == nil {
		t.Fatal("expected error for negative column")
	}
}

func TestNew_Valid(t *testing.T) {
	ex, err := New(" ", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex == nil {
		t.Fatal("expected non-nil extractor")
	}
}

func TestExtract_FirstColumn(t *testing.T) {
	ex, _ := New(" ", 0)
	val, ok := ex.Extract("2024-01-15 12:00:00 INFO message")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if val != "2024-01-15" {
		t.Errorf("expected '2024-01-15', got %q", val)
	}
}

func TestExtract_MiddleColumn(t *testing.T) {
	ex, _ := New(" ", 2)
	val, ok := ex.Extract("2024-01-15 12:00:00 INFO starting server")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if val != "INFO" {
		t.Errorf("expected 'INFO', got %q", val)
	}
}

func TestExtract_TabDelimiter(t *testing.T) {
	ex, _ := New("\t", 1)
	val, ok := ex.Extract("host\t192.168.1.1\tGET")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if val != "192.168.1.1" {
		t.Errorf("expected '192.168.1.1', got %q", val)
	}
}

func TestExtract_ColumnOutOfRange(t *testing.T) {
	ex, _ := New(" ", 10)
	_, ok := ex.Extract("only three words")
	if ok {
		t.Fatal("expected ok=false when column out of range")
	}
}

func TestExtract_EmptyLine(t *testing.T) {
	ex, _ := New(" ", 0)
	_, ok := ex.Extract("")
	if ok {
		t.Fatal("expected ok=false for empty line")
	}
}

func TestExtract_WhitespaceOnlyColumn(t *testing.T) {
	ex, _ := New("|", 1)
	_, ok := ex.Extract("foo|   |bar")
	if ok {
		t.Fatal("expected ok=false for whitespace-only column")
	}
}

func TestColumnCount(t *testing.T) {
	ex, _ := New(" ", 0)
	count := ex.ColumnCount("a b c d")
	if count != 4 {
		t.Errorf("expected 4, got %d", count)
	}
}

func TestColumnCount_EmptyLine(t *testing.T) {
	ex, _ := New(" ", 0)
	count := ex.ColumnCount("")
	if count != 0 {
		t.Errorf("expected 0, got %d", count)
	}
}
