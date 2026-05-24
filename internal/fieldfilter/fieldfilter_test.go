package fieldfilter

import (
	"testing"
)

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", "error", false)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New("level", "", false)
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidRegex(t *testing.T) {
	_, err := New("level", "[invalid", false)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_Valid(t *testing.T) {
	f, err := New("level", "error", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestKeep_JSON_Match(t *testing.T) {
	f, _ := New("level", "error", false)
	line := `{"time":"2024-01-01T00:00:00Z","level":"error","msg":"oops"}`
	if !f.Keep(line) {
		t.Error("expected line to be kept")
	}
}

func TestKeep_JSON_NoMatch(t *testing.T) {
	f, _ := New("level", "error", false)
	line := `{"time":"2024-01-01T00:00:00Z","level":"info","msg":"ok"}`
	if f.Keep(line) {
		t.Error("expected line to be dropped")
	}
}

func TestKeep_Logfmt_Match(t *testing.T) {
	f, _ := New("level", "warn", false)
	line := `time=2024-01-01T00:00:00Z level=warn msg="disk full"`
	if !f.Keep(line) {
		t.Error("expected line to be kept")
	}
}

func TestKeep_Logfmt_NoMatch(t *testing.T) {
	f, _ := New("level", "warn", false)
	line := `time=2024-01-01T00:00:00Z level=info msg="all good"`
	if f.Keep(line) {
		t.Error("expected line to be dropped")
	}
}

func TestKeep_Invert(t *testing.T) {
	f, _ := New("level", "debug", true)
	keepLine := `{"level":"info","msg":"hello"}`
	dropLine := `{"level":"debug","msg":"verbose"}`
	if !f.Keep(keepLine) {
		t.Error("expected non-debug line to be kept with invert")
	}
	if f.Keep(dropLine) {
		t.Error("expected debug line to be dropped with invert")
	}
}

func TestKeep_MissingField_NotInvert(t *testing.T) {
	f, _ := New("level", "error", false)
	line := `{"msg":"no level field"}`
	if f.Keep(line) {
		t.Error("expected line without field to be dropped")
	}
}

func TestKeep_MissingField_Invert(t *testing.T) {
	f, _ := New("level", "error", true)
	line := `{"msg":"no level field"}`
	if !f.Keep(line) {
		t.Error("expected line without field to be kept when inverted")
	}
}
