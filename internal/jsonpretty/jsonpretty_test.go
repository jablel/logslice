package jsonpretty

import (
	"strings"
	"testing"
)

func TestNew_ValidIndent(t *testing.T) {
	f, err := New("  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.IsEnabled() {
		t.Error("expected formatter to be enabled")
	}
}

func TestNew_EmptyIndentDisabled(t *testing.T) {
	f, err := New("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.IsEnabled() {
		t.Error("expected formatter to be disabled when indent is empty")
	}
}

func TestNew_IndentTooLong(t *testing.T) {
	_, err := New("         ") // 9 spaces
	if err == nil {
		t.Fatal("expected error for indent > 8 chars")
	}
}

func TestApply_ValidJSON(t *testing.T) {
	f, _ := New("  ")
	input := `{"level":"info","msg":"hello"}`
	out := f.Apply(input)
	if !strings.Contains(out, "\n") {
		t.Error("expected multi-line pretty output")
	}
	if !strings.Contains(out, `"level"`) {
		t.Error("expected field 'level' in output")
	}
}

func TestApply_InvalidJSON_PassThrough(t *testing.T) {
	f, _ := New("  ")
	input := "not json at all"
	out := f.Apply(input)
	if out != input {
		t.Errorf("expected passthrough, got %q", out)
	}
}

func TestApply_Disabled_PassThrough(t *testing.T) {
	f, _ := New("")
	input := `{"level":"warn"}`
	out := f.Apply(input)
	if out != input {
		t.Errorf("expected passthrough when disabled, got %q", out)
	}
}

func TestApply_TabIndent(t *testing.T) {
	f, _ := New("\t")
	input := `{"a":1}`
	out := f.Apply(input)
	if !strings.Contains(out, "\t") {
		t.Error("expected tab indentation in output")
	}
}

func TestApply_NestedJSON(t *testing.T) {
	f, _ := New("  ")
	input := `{"meta":{"host":"srv1"},"code":200}`
	out := f.Apply(input)
	if !strings.Contains(out, "meta") || !strings.Contains(out, "host") {
		t.Error("expected nested fields in pretty output")
	}
}
