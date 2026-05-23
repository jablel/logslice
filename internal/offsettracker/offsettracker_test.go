package offsettracker_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/offsettracker"
)

func TestNew_Enabled(t *testing.T) {
	tr := offsettracker.New(true)
	if tr == nil {
		t.Fatal("expected non-nil tracker")
	}
}

func TestNew_Disabled(t *testing.T) {
	tr := offsettracker.New(false)
	tr.Record(1, "some line")
	if len(tr.Entries()) != 0 {
		t.Errorf("disabled tracker should not record entries")
	}
}

func TestRecord_StoresEntry(t *testing.T) {
	tr := offsettracker.New(true)
	tr.Record(1, "first line")
	entries := tr.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Line != "first line" {
		t.Errorf("unexpected line: %q", entries[0].Line)
	}
	if entries[0].LineNumber != 1 {
		t.Errorf("unexpected line number: %d", entries[0].LineNumber)
	}
}

func TestAdvance_UpdatesByteOffset(t *testing.T) {
	tr := offsettracker.New(true)
	tr.Advance(20) // simulate consuming a 20-byte line
	tr.Record(2, "second line")
	entries := tr.Entries()
	if entries[0].ByteOffset != 20 {
		t.Errorf("expected byte offset 20, got %d", entries[0].ByteOffset)
	}
}

func TestAdvance_AccumulatesOffset(t *testing.T) {
	tr := offsettracker.New(true)
	tr.Record(1, "line one")
	tr.Advance(10)
	tr.Record(2, "line two")
	tr.Advance(15)
	tr.Record(3, "line three")

	entries := tr.Entries()
	if entries[1].ByteOffset != 10 {
		t.Errorf("expected 10, got %d", entries[1].ByteOffset)
	}
	if entries[2].ByteOffset != 25 {
		t.Errorf("expected 25, got %d", entries[2].ByteOffset)
	}
}

func TestReset_ClearsState(t *testing.T) {
	tr := offsettracker.New(true)
	tr.Advance(50)
	tr.Record(1, "some line")
	tr.Reset()

	if len(tr.Entries()) != 0 {
		t.Errorf("expected 0 entries after reset")
	}
	// After reset, next record should have offset 0
	tr.Record(1, "fresh line")
	if tr.Entries()[0].ByteOffset != 0 {
		t.Errorf("expected byte offset 0 after reset")
	}
}

func TestSummary_Enabled(t *testing.T) {
	tr := offsettracker.New(true)
	tr.Advance(100)
	tr.Record(1, "line")
	s := tr.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestSummary_Disabled(t *testing.T) {
	tr := offsettracker.New(false)
	s := tr.Summary()
	if s != "offset tracking disabled" {
		t.Errorf("unexpected summary: %q", s)
	}
}
