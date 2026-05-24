package severityfilter

import (
	"errors"
	"testing"
)

func TestNew_ValidLevel(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  Level
	}{
		{"DEBUG", DEBUG},
		{"info", INFO},
		{"Warn", WARN},
		{"WARNING", WARN},
		{"ERROR", ERROR},
		{"err", ERROR},
		{"FATAL", FATAL},
		{"crit", FATAL},
	} {
		f, err := New(tc.input)
		if err != nil {
			t.Fatalf("New(%q) unexpected error: %v", tc.input, err)
		}
		if f.Min() != tc.want {
			t.Errorf("New(%q).Min() = %d, want %d", tc.input, f.Min(), tc.want)
		}
	}
}

func TestNew_InvalidLevel(t *testing.T) {
	_, err := New("VERBOSE")
	if err == nil {
		t.Fatal("expected error for unknown level, got nil")
	}
	if !errors.Is(err, ErrUnknownLevel) {
		t.Errorf("expected ErrUnknownLevel, got %v", err)
	}
}

func TestKeep_AboveMinLevel(t *testing.T) {
	f, _ := New("WARN")
	cases := []struct {
		line string
		want bool
	}{
		{"2024-01-01 DEBUG starting up", false},
		{"2024-01-01 INFO  server ready", false},
		{"2024-01-01 WARN  disk usage high", true},
		{"2024-01-01 ERROR connection refused", true},
		{"2024-01-01 FATAL out of memory", true},
	}
	for _, tc := range cases {
		got := f.Keep(tc.line)
		if got != tc.want {
			t.Errorf("Keep(%q) = %v, want %v", tc.line, got, tc.want)
		}
	}
}

func TestKeep_NoLevelPassesThrough(t *testing.T) {
	f, _ := New("ERROR")
	line := "plain log line without any level token"
	if !f.Keep(line) {
		t.Errorf("Keep(%q) = false, want true (pass-through)", line)
	}
}

func TestKeep_CaseInsensitiveDetection(t *testing.T) {
	f, _ := New("INFO")
	for _, line := range []string{
		"level=error msg=\"oops\"",
		"[Error] something went wrong",
		"{\"level\":\"ERROR\"}",
	} {
		if !f.Keep(line) {
			t.Errorf("Keep(%q) = false, want true", line)
		}
	}
}

func TestKeep_DebugFilteredAtErrorMin(t *testing.T) {
	f, _ := New("ERROR")
	if f.Keep("DEBUG initialising cache") {
		t.Error("DEBUG line should be filtered when min=ERROR")
	}
}
