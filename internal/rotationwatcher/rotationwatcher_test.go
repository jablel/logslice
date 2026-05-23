package rotationwatcher_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/rotationwatcher"
)

func writeTmp(t *testing.T, content string) *os.File {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "log-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	_ = f.Close()
	return f
}

func TestNew_EmptyPath(t *testing.T) {
	_, err := rotationwatcher.New("", time.Second)
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNew_ZeroPollInterval(t *testing.T) {
	f := writeTmp(t, "hello\n")
	_, err := rotationwatcher.New(f.Name(), 0)
	if err == nil {
		t.Fatal("expected error for zero poll interval")
	}
}

func TestNew_FileNotFound(t *testing.T) {
	_, err := rotationwatcher.New("/nonexistent/path/logfile.log", time.Second)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNew_Valid(t *testing.T) {
	f := writeTmp(t, "line1\n")
	w, err := rotationwatcher.New(f.Name(), 50*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.Stop()
	if w.Rotated() {
		t.Fatal("expected no rotation on fresh watcher")
	}
}

func TestRotated_DetectedOnTruncation(t *testing.T) {
	f := writeTmp(t, "some initial content\n")
	w, err := rotationwatcher.New(f.Name(), 20*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.Start()
	defer w.Stop()

	// Truncate the file to simulate rotation.
	if err := os.WriteFile(f.Name(), []byte(""), 0644); err != nil {
		t.Fatalf("truncate: %v", err)
	}

	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		if w.Rotated() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !w.Rotated() {
		t.Fatal("expected rotation to be detected after truncation")
	}
}

func TestReset_ClearsRotationFlag(t *testing.T) {
	f := writeTmp(t, "data\n")
	w, err := rotationwatcher.New(f.Name(), 20*time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	w.Start()
	defer w.Stop()

	_ = os.WriteFile(f.Name(), []byte(""), 0644)
	time.Sleep(100 * time.Millisecond)

	if !w.Rotated() {
		t.Skip("rotation not detected; skipping reset test")
	}
	if err := w.Reset(); err != nil {
		t.Fatalf("reset error: %v", err)
	}
	if w.Rotated() {
		t.Fatal("expected rotation flag to be cleared after Reset")
	}
}
