package lineformatter_test

import (
	"encoding/json"
	"strings"
	"testing"

	"logslice/internal/lineformatter"
)

func TestNew_ValidFormats(t *testing.T) {
	for _, f := range []lineformatter.Format{lineformatter.FormatRaw, lineformatter.FormatJSON, lineformatter.FormatText} {
		_, err := lineformatter.New(f, nil, true)
		if err != nil {
			t.Errorf("expected no error for format %d, got %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := lineformatter.New(lineformatter.Format(99), nil, true)
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestApply_Raw_ReturnsUnchanged(t *testing.T) {
	f, _ := lineformatter.New(lineformatter.FormatRaw, nil, true)
	line := `{"level":"info","msg":"hello"}`
	if got := f.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_JSON_AllFields(t *testing.T) {
	f, _ := lineformatter.New(lineformatter.FormatJSON, nil, true)
	line := `{"level":"info","msg":"hello"}`
	got := f.Apply(line)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["level"] != "info" || m["msg"] != "hello" {
		t.Errorf("unexpected fields in output: %v", m)
	}
}

func TestApply_JSON_IncludeFields(t *testing.T) {
	f, _ := lineformatter.New(lineformatter.FormatJSON, []string{"level"}, true)
	line := `{"level":"warn","msg":"oops","ts":"2024-01-01"}`
	got := f.Apply(line)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := m["level"]; !ok {
		t.Error("expected 'level' field to be present")
	}
	if _, ok := m["msg"]; ok {
		t.Error("expected 'msg' field to be excluded")
	}
}

func TestApply_JSON_ExcludeFields(t *testing.T) {
	f, _ := lineformatter.New(lineformatter.FormatJSON, []string{"ts"}, false)
	line := `{"level":"error","msg":"fail","ts":"2024-01-01"}`
	got := f.Apply(line)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if _, ok := m["ts"]; ok {
		t.Error("expected 'ts' field to be excluded")
	}
}

func TestApply_Text_ContainsKeyValue(t *testing.T) {
	f, _ := lineformatter.New(lineformatter.FormatText, []string{"level"}, true)
	line := `{"level":"debug","msg":"trace"}`
	got := f.Apply(line)
	if !strings.Contains(got, "level=debug") {
		t.Errorf("expected key=value pair in %q", got)
	}
}

func TestApply_NonJSON_WrapsAsMessage(t *testing.T) {
	f, _ := lineformatter.New(lineformatter.FormatJSON, nil, true)
	line := "plain text log line"
	got := f.Apply(line)
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(got), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["message"] != line {
		t.Errorf("expected message=%q, got %v", line, m["message"])
	}
}
