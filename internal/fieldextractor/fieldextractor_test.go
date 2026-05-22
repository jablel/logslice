package fieldextractor

import (
	"testing"
)

func TestNew_ValidField(t *testing.T) {
	ex, err := New("level", FormatUnknown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ex == nil {
		t.Fatal("expected non-nil extractor")
	}
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", FormatUnknown)
	if err != ErrEmptyField {
		t.Fatalf("expected ErrEmptyField, got %v", err)
	}
}

func TestExtract_JSON_KnownFormat(t *testing.T) {
	ex, _ := New("level", FormatJSON)
	got := ex.Extract(`{"time":"2024-01-01T00:00:00Z","level":"error","msg":"fail"}`)
	if got != "error" {
		t.Errorf("expected 'error', got %q", got)
	}
}

func TestExtract_JSON_AutoDetect(t *testing.T) {
	ex, _ := New("msg", FormatUnknown)
	got := ex.Extract(`{"msg":"hello world"}`)
	if got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestExtract_JSON_MissingField(t *testing.T) {
	ex, _ := New("missing", FormatJSON)
	got := ex.Extract(`{"level":"info"}`)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestExtract_Logfmt_Unquoted(t *testing.T) {
	ex, _ := New("level", FormatLogfmt)
	got := ex.Extract(`time=2024-01-01 level=warn msg=something`)
	if got != "warn" {
		t.Errorf("expected 'warn', got %q", got)
	}
}

func TestExtract_Logfmt_Quoted(t *testing.T) {
	ex, _ := New("msg", FormatLogfmt)
	got := ex.Extract(`level=info msg="disk full" host=srv1`)
	if got != "disk full" {
		t.Errorf("expected 'disk full', got %q", got)
	}
}

func TestExtract_Logfmt_AutoDetect(t *testing.T) {
	ex, _ := New("host", FormatUnknown)
	got := ex.Extract(`level=debug host=myserver msg=ok`)
	if got != "myserver" {
		t.Errorf("expected 'myserver', got %q", got)
	}
}

func TestExtract_Logfmt_MissingField(t *testing.T) {
	ex, _ := New("nothere", FormatLogfmt)
	got := ex.Extract(`level=info msg=ok`)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestExtract_InvalidJSON(t *testing.T) {
	ex, _ := New("level", FormatJSON)
	got := ex.Extract(`not valid json`)
	if got != "" {
		t.Errorf("expected empty string for invalid JSON, got %q", got)
	}
}
