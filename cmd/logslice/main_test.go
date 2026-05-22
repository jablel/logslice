package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/output"
	"github.com/logslice/logslice/internal/scanner"
)

func TestIntegration_SliceByTime(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-15T09:00:00Z INFO startup",
		"2024-01-15T10:00:00Z INFO request received",
		"2024-01-15T11:00:00Z WARN slow query",
		"2024-01-15T12:00:00Z ERROR timeout",
		"2024-01-15T13:00:00Z INFO shutdown",
	}, "\n")

	f, err := filter.New("2024-01-15T10:00:00Z", "2024-01-15T12:00:00Z", "")
	if err != nil {
		t.Fatalf("filter.New: %v", err)
	}

	var buf bytes.Buffer
	w := output.New(output.Options{Destination: &buf})

	sc := scanner.New(strings.NewReader(input))
	for sc.Scan() {
		line := sc.Text()
		if f.Match(line) {
			_ = w.WriteLine(sc.LineNumber(), line)
		}
	}
	_ = w.Flush()

	got := buf.String()
	expected := []string{
		"2024-01-15T10:00:00Z INFO request received",
		"2024-01-15T11:00:00Z WARN slow query",
		"2024-01-15T12:00:00Z ERROR timeout",
	}
	for _, e := range expected {
		if !strings.Contains(got, e) {
			t.Errorf("expected output to contain %q\ngot:\n%s", e, got)
		}
	}
	if strings.Contains(got, "startup") || strings.Contains(got, "shutdown") {
		t.Errorf("output should not contain out-of-range lines:\n%s", got)
	}
	if w.Count() != 3 {
		t.Errorf("Count() = %d, want 3", w.Count())
	}
}
