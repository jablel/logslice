package grepfilter_test

import (
	"testing"

	"logslice/internal/grepfilter"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := grepfilter.New(grepfilter.Options{Pattern: ""})
	if err == nil {
		t.Fatal("expected error for empty pattern, got nil")
	}
}

func TestNew_InvalidRegex(t *testing.T) {
	_, err := grepfilter.New(grepfilter.Options{Pattern: "[invalid", Regex: true})
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestNew_ValidFixed(t *testing.T) {
	f, err := grepfilter.New(grepfilter.Options{Pattern: "ERROR"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil filter")
	}
}

func TestKeep_FixedSubstring(t *testing.T) {
	f, _ := grepfilter.New(grepfilter.Options{Pattern: "ERROR"})

	if !f.Keep("2024-01-01 ERROR something went wrong") {
		t.Error("expected match for line containing ERROR")
	}
	if f.Keep("2024-01-01 INFO all good") {
		t.Error("expected no match for line without ERROR")
	}
}

func TestKeep_Regex(t *testing.T) {
	f, _ := grepfilter.New(grepfilter.Options{Pattern: `ERR(OR)?`, Regex: true})

	if !f.Keep("ERR: disk full") {
		t.Error("expected match for ERR")
	}
	if !f.Keep("ERROR: disk full") {
		t.Error("expected match for ERROR")
	}
	if f.Keep("INFO: all fine") {
		t.Error("expected no match for INFO line")
	}
}

func TestKeep_Invert(t *testing.T) {
	f, _ := grepfilter.New(grepfilter.Options{Pattern: "DEBUG", Invert: true})

	if !f.Keep("INFO: starting") {
		t.Error("expected non-DEBUG line to be kept")
	}
	if f.Keep("DEBUG: verbose output") {
		t.Error("expected DEBUG line to be dropped when inverted")
	}
}

func TestKeep_CaseInsensitive(t *testing.T) {
	f, _ := grepfilter.New(grepfilter.Options{Pattern: "error", CaseInsensitive: true})

	if !f.Keep("ERROR: something") {
		t.Error("expected case-insensitive match for ERROR")
	}
	if !f.Keep("Error: something") {
		t.Error("expected case-insensitive match for Error")
	}
}

func TestKeep_FixedPatternNotTreatedAsRegex(t *testing.T) {
	// Dot should be literal, not match any character.
	f, _ := grepfilter.New(grepfilter.Options{Pattern: "1.0", Regex: false})

	if !f.Keep("version 1.0 released") {
		t.Error("expected literal match for '1.0'")
	}
	if f.Keep("version 1X0 released") {
		t.Error("dot must not act as wildcard in fixed mode")
	}
}
