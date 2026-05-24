package timestampinjector_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/timestampextractor"
	"github.com/user/logslice/internal/timestampinjector"
)

func newExtractor(t *testing.T) *timestampextractor.Extractor {
	t.Helper()
	ex, err := timestampextractor.New(nil)
	if err != nil {
		t.Fatalf("extractor: %v", err)
	}
	return ex
}

func TestNew_NilExtractor(t *testing.T) {
	_, err := timestampinjector.New(nil, "", timestampinjector.Prepend)
	if err == nil {
		t.Fatal("expected error for nil extractor")
	}
}

func TestNew_DefaultFormat(t *testing.T) {
	ex := newExtractor(t)
	inj, err := timestampinjector.New(ex, "", timestampinjector.Prepend)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inj == nil {
		t.Fatal("expected non-nil injector")
	}
}

func TestApply_NoTimestamp_ReturnsOriginal(t *testing.T) {
	ex := newExtractor(t)
	inj, _ := timestampinjector.New(ex, time.RFC3339, timestampinjector.Prepend)
	line := "no timestamp here at all"
	got := inj.Apply(line)
	if got != line {
		t.Errorf("want %q, got %q", line, got)
	}
}

func TestApply_Prepend(t *testing.T) {
	ex := newExtractor(t)
	inj, _ := timestampinjector.New(ex, "2006-01-02", timestampinjector.Prepend)
	line := "2024-03-15T10:00:00Z INFO server started"
	got := inj.Apply(line)
	const wantPrefix = "2024-03-15 "
	if len(got) < len(wantPrefix) || got[:len(wantPrefix)] != wantPrefix {
		t.Errorf("want prefix %q, got %q", wantPrefix, got)
	}
}

func TestApply_Append(t *testing.T) {
	ex := newExtractor(t)
	inj, _ := timestampinjector.New(ex, "2006-01-02", timestampinjector.Append)
	line := "2024-03-15T10:00:00Z INFO server started"
	got := inj.Apply(line)
	const wantSuffix = " 2024-03-15"
	if len(got) < len(wantSuffix) || got[len(got)-len(wantSuffix):] != wantSuffix {
		t.Errorf("want suffix %q, got %q", wantSuffix, got)
	}
}

func TestApply_StripNewline(t *testing.T) {
	ex := newExtractor(t)
	inj, _ := timestampinjector.New(ex, time.RFC3339, timestampinjector.Prepend)
	line := "2024-03-15T10:00:00Z INFO done\n"
	got := inj.Apply(line)
	for _, c := range got {
		if c == '\n' {
			t.Errorf("result should not contain newline, got %q", got)
		}
	}
}
