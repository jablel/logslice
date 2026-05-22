package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"logslice/internal/stats"
)

func TestNew(t *testing.T) {
	s := stats.New()
	if s == nil {
		t.Fatal("New() returned nil")
	}
	if s.StartTime.IsZero() {
		t.Error("StartTime should be set on New()")
	}
}

func TestFinishAndDuration(t *testing.T) {
	s := stats.New()
	time.Sleep(10 * time.Millisecond)
	s.Finish()

	d := s.Duration()
	if d < 10*time.Millisecond {
		t.Errorf("Duration() = %s, want >= 10ms", d)
	}
	// calling Duration again after Finish should be stable
	if s.Duration() != d {
		t.Error("Duration() should be stable after Finish()")
	}
}

func TestMatchRate(t *testing.T) {
	tests := []struct {
		scanned, matched int64
		want             float64
	}{
		{0, 0, 0},
		{100, 50, 0.5},
		{200, 200, 1.0},
		{1000, 1, 0.001},
	}
	for _, tc := range tests {
		s := stats.New()
		s.LinesScanned = tc.scanned
		s.LinesMatched = tc.matched
		got := s.MatchRate()
		if got != tc.want {
			t.Errorf("MatchRate(%d/%d) = %f, want %f", tc.matched, tc.scanned, got, tc.want)
		}
	}
}

func TestPrint(t *testing.T) {
	s := stats.New()
	s.LinesScanned = 1000
	s.LinesMatched = 42
	s.LinesSkipped = 958
	s.BytesRead = 204800
	s.Finish()

	var buf bytes.Buffer
	s.Print(&buf)
	out := buf.String()

	for _, want := range []string{"1000", "42", "958", "204800", "logslice stats"} {
		if !strings.Contains(out, want) {
			t.Errorf("Print() output missing %q\ngot:\n%s", want, out)
		}
	}
}

func TestPrint_NilWriterUsesStderr(t *testing.T) {
	// Just ensure it doesn't panic when w is nil.
	s := stats.New()
	s.Finish()
	// We can't easily capture stderr here, so we only verify no panic.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Print(nil) panicked: %v", r)
		}
	}()
	s.Print(nil)
}
