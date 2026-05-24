package multipatternfilter_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/multipatternfilter"
)

func TestNew_NoPatterns(t *testing.T) {
	_, err := multipatternfilter.New(nil, multipatternfilter.ModeOr)
	if err != multipatternfilter.ErrNoPatterns {
		t.Fatalf("expected ErrNoPatterns, got %v", err)
	}
}

func TestNew_InvalidRegex(t *testing.T) {
	_, err := multipatternfilter.New([]string{"[invalid"}, multipatternfilter.ModeOr)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_UnknownMode(t *testing.T) {
	_, err := multipatternfilter.New([]string{"ok"}, multipatternfilter.Mode(99))
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestNew_Valid(t *testing.T) {
	f, err := multipatternfilter.New([]string{"ERROR", "WARN"}, multipatternfilter.ModeOr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Len() != 2 {
		t.Fatalf("expected 2 patterns, got %d", f.Len())
	}
}

func TestKeep_ModeOr_Matches(t *testing.T) {
	f, _ := multipatternfilter.New([]string{"ERROR", "WARN"}, multipatternfilter.ModeOr)
	if !f.Keep("2024-01-01 ERROR disk full") {
		t.Error("expected match on ERROR line")
	}
	if !f.Keep("2024-01-01 WARN low memory") {
		t.Error("expected match on WARN line")
	}
	if f.Keep("2024-01-01 INFO all good") {
		t.Error("expected no match on INFO line")
	}
}

func TestKeep_ModeAnd_RequiresAll(t *testing.T) {
	f, _ := multipatternfilter.New([]string{"ERROR", "disk"}, multipatternfilter.ModeAnd)
	if !f.Keep("ERROR: disk full") {
		t.Error("expected match when both patterns present")
	}
	if f.Keep("ERROR: network timeout") {
		t.Error("expected no match when only one pattern present")
	}
	if f.Keep("INFO: disk usage high") {
		t.Error("expected no match when only second pattern present")
	}
}

func TestKeep_EmptyLine(t *testing.T) {
	f, _ := multipatternfilter.New([]string{"ERROR"}, multipatternfilter.ModeOr)
	if f.Keep("") {
		t.Error("expected empty line not to match")
	}
}

func TestKeep_ModeAnd_SinglePattern(t *testing.T) {
	f, _ := multipatternfilter.New([]string{"timeout"}, multipatternfilter.ModeAnd)
	if !f.Keep("connection timeout after 30s") {
		t.Error("expected match")
	}
	if f.Keep("connection established") {
		t.Error("expected no match")
	}
}
