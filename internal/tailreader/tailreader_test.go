package tailreader_test

import (
	"io"
	"os"
	"testing"

	"github.com/yourorg/logslice/internal/tailreader"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "tailreader-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestNew_EmptyPath(t *testing.T) {
	_, err := tailreader.New("", 0)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNew_NegativeOffset(t *testing.T) {
	_, err := tailreader.New("/dev/null", -1)
	if err == nil {
		t.Fatal("expected error for negative offset")
	}
}

func TestNew_FileNotFound(t *testing.T) {
	_, err := tailreader.New("/nonexistent/path/file.log", 0)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestReadLine_AllLines(t *testing.T) {
	path := writeTempFile(t, "line1\nline2\nline3\n")
	r, err := tailreader.New(path, 0)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	expected := []string{"line1", "line2", "line3"}
	for _, want := range expected {
		got, _, err := r.ReadLine()
		if err != nil {
			t.Fatalf("ReadLine: %v", err)
		}
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
	_, _, err = r.ReadLine()
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}
}

func TestReadLine_FromOffset(t *testing.T) {
	// "line1\n" is 6 bytes; start reading from byte 6 to skip it
	path := writeTempFile(t, "line1\nline2\nline3\n")
	r, err := tailreader.New(path, 6)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	got, _, err := r.ReadLine()
	if err != nil {
		t.Fatalf("ReadLine: %v", err)
	}
	if got != "line2" {
		t.Errorf("got %q, want %q", got, "line2")
	}
}

func TestOffset_AdvancesCorrectly(t *testing.T) {
	path := writeTempFile(t, "hello\nworld\n")
	r, err := tailreader.New(path, 0)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer r.Close()

	if r.Offset() != 0 {
		t.Errorf("initial offset: got %d, want 0", r.Offset())
	}
	_, off, _ := r.ReadLine() // "hello\n" = 6 bytes
	if off != 6 {
		t.Errorf("offset after first line: got %d, want 6", off)
	}
}
