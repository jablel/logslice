package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/logslice/logslice/internal/output"
)

func TestWriter_FormatRaw(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.Options{
		Format:      output.FormatRaw,
		Destination: &buf,
	})

	lines := []string{"line one", "line two", "line three"}
	for i, l := range lines {
		if err := w.WriteLine(i+1, l); err != nil {
			t.Fatalf("WriteLine(%d) error: %v", i+1, err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush error: %v", err)
	}

	got := buf.String()
	for _, l := range lines {
		if !strings.Contains(got, l) {
			t.Errorf("expected output to contain %q, got:\n%s", l, got)
		}
	}
}

func TestWriter_FormatNumbered(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.Options{
		Format:      output.FormatNumbered,
		Destination: &buf,
	})

	if err := w.WriteLine(42, "some log entry"); err != nil {
		t.Fatalf("WriteLine error: %v", err)
	}
	_ = w.Flush()

	got := buf.String()
	if !strings.Contains(got, "42\t") {
		t.Errorf("expected numbered line to contain '42\\t', got: %q", got)
	}
	if !strings.Contains(got, "some log entry") {
		t.Errorf("expected output to contain log entry, got: %q", got)
	}
}

func TestWriter_Count(t *testing.T) {
	var buf bytes.Buffer
	w := output.New(output.Options{Destination: &buf})

	for i := 0; i < 5; i++ {
		_ = w.WriteLine(i, "line")
	}
	if got := w.Count(); got != 5 {
		t.Errorf("Count() = %d, want 5", got)
	}
}

func TestWriter_DefaultsToStdout(t *testing.T) {
	// Just ensure no panic when Destination is nil.
	w := output.New(output.Options{})
	if w == nil {
		t.Fatal("expected non-nil writer")
	}
}
