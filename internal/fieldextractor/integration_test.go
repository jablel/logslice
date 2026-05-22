package fieldextractor_test

import (
	"testing"

	"logslice/internal/fieldextractor"
)

var mixedLines = []struct {
	line  string
	wantLevel string
}{
	{`{"time":"2024-01-01T00:00:00Z","level":"error","msg":"crash"}`, "error"},
	{`{"time":"2024-01-01T00:01:00Z","level":"info","msg":"started"}`, "info"},
	{`time=2024-01-01T00:02:00Z level=warn msg=slowdown`, "warn"},
	{`time=2024-01-01T00:03:00Z level=debug msg=tick`, "debug"},
	{`not a structured line at all`, ""},
}

func TestIntegration_AutoDetectMixedFormats(t *testing.T) {
	ex, err := fieldextractor.New("level", fieldextractor.FormatUnknown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, tc := range mixedLines {
		got := ex.Extract(tc.line)
		if got != tc.wantLevel {
			t.Errorf("line %q: expected %q, got %q", tc.line, tc.wantLevel, got)
		}
	}
}

func TestIntegration_FilterByField(t *testing.T) {
	ex, _ := fieldextractor.New("level", fieldextractor.FormatUnknown)

	var errors []string
	for _, tc := range mixedLines {
		if ex.Extract(tc.line) == "error" {
			errors = append(errors, tc.line)
		}
	}

	if len(errors) != 1 {
		t.Errorf("expected 1 error line, got %d", len(errors))
	}
}
