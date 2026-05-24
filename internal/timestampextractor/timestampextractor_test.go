package timestampextractor_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/timestampextractor"
)

func TestNew_DefaultFormats(t *testing.T) {
	ex := timestampextractor.New(nil)
	if ex == nil {
		t.Fatal("expected non-nil extractor")
	}
}

func TestNew_CustomFormats(t *testing.T) {
	ex := timestampextractor.New([]string{"2006-01-02"})
	if ex == nil {
		t.Fatal("expected non-nil extractor")
	}
}

func TestExtract_RFC3339AtStart(t *testing.T) {
	ex := timestampextractor.New([]string{time.RFC3339})
	line := "2024-06-15T08:23:01Z INFO server started"
	res, err := ex.Extract(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Offset != 0 {
		t.Errorf("expected offset 0, got %d", res.Offset)
	}
	if res.Raw != "2024-06-15T08:23:01Z" {
		t.Errorf("unexpected raw: %q", res.Raw)
	}
	if res.Format != time.RFC3339 {
		t.Errorf("unexpected format: %q", res.Format)
	}
	if res.Time.IsZero() {
		t.Error("expected non-zero time")
	}
}

func TestExtract_TimestampAfterPrefix(t *testing.T) {
	ex := timestampextractor.New([]string{"2006-01-02"})
	line := "[INFO] 2024-06-15 something happened"
	res, err := ex.Extract(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Raw != "2024-06-15" {
		t.Errorf("unexpected raw: %q", res.Raw)
	}
}

func TestExtract_NoTimestamp(t *testing.T) {
	ex := timestampextractor.New([]string{time.RFC3339})
	_, err := ex.Extract("this line has no timestamp at all")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err != timestampextractor.ErrNoTimestamp {
		t.Errorf("expected ErrNoTimestamp, got %v", err)
	}
}

func TestExtract_EmptyLine(t *testing.T) {
	ex := timestampextractor.New(nil)
	_, err := ex.Extract("")
	if err != timestampextractor.ErrNoTimestamp {
		t.Errorf("expected ErrNoTimestamp for empty line, got %v", err)
	}
}

func TestExtract_MultipleFormats_FirstWins(t *testing.T) {
	formats := []string{"2006-01-02", time.RFC3339}
	ex := timestampextractor.New(formats)
	line := "2024-06-15T08:23:01Z some event"
	res, err := ex.Extract(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// "2006-01-02" is shorter and will match the date prefix first.
	if res.Format != "2006-01-02" {
		t.Errorf("expected first format to win, got %q", res.Format)
	}
}
