package linemerger

import (
	"strings"
	"testing"
)

func TestNew_NilPredicate(t *testing.T) {
	_, err := New(nil, "\n")
	if err == nil {
		t.Fatal("expected error for nil predicate")
	}
}

func TestNew_DefaultDelimiter(t *testing.T) {
	m, err := New(func(string) bool { return false }, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.delimiter != "\n" {
		t.Errorf("expected delimiter '\\n', got %q", m.delimiter)
	}
}

func TestNewIndentMerger_Valid(t *testing.T) {
	m, err := NewIndentMerger()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil merger")
	}
}

func TestNewPatternMerger_InvalidRegex(t *testing.T) {
	_, err := NewPatternMerger("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func collectRecords(m *Merger, lines []string) []string {
	var records []string
	for _, l := range lines {
		if rec, ok := m.Feed(l); ok {
			records = append(records, rec)
		}
	}
	if rec, ok := m.Flush(); ok {
		records = append(records, rec)
	}
	return records
}

func TestFeed_NoMerge(t *testing.T) {
	m, _ := NewIndentMerger()
	lines := []string{"line1", "line2", "line3"}
	records := collectRecords(m, lines)
	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d: %v", len(records), records)
	}
}

func TestFeed_MergesIndentedContinuations(t *testing.T) {
	m, _ := NewIndentMerger()
	lines := []string{
		"ERROR exception occurred",
		"\tat com.example.Foo.bar(Foo.java:42)",
		"\tat com.example.Main.main(Main.java:10)",
		"INFO next event",
	}
	records := collectRecords(m, lines)
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if !strings.Contains(records[0], "Foo.java") {
		t.Errorf("expected first record to contain stack trace, got: %s", records[0])
	}
	if records[1] != "INFO next event" {
		t.Errorf("unexpected second record: %s", records[1])
	}
}

func TestFlush_EmptyBuffer(t *testing.T) {
	m, _ := NewIndentMerger()
	_, ok := m.Flush()
	if ok {
		t.Error("expected Flush to return false on empty buffer")
	}
}

func TestNewPatternMerger_MergesByPattern(t *testing.T) {
	m, err := NewPatternMerger(`^\+`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := []string{"base line", "+ continuation", "+ more", "next base"}
	records := collectRecords(m, lines)
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d: %v", len(records), records)
	}
	if !strings.Contains(records[0], "continuation") {
		t.Errorf("expected merged record, got: %s", records[0])
	}
}
