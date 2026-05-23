package invertfilter_test

import (
	"testing"

	"logslice/internal/invertfilter"
)

// stubKeeper is a simple Keeper that keeps lines containing a fixed substring.
type stubKeeper struct{ sub string }

func (s *stubKeeper) Keep(line string) bool {
	return len(line) > 0 && containsSubstr(line, s.sub)
}

func containsSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func TestNew_NilKeeper(t *testing.T) {
	_, err := invertfilter.New(nil, true)
	if err == nil {
		t.Fatal("expected error for nil keeper, got nil")
	}
}

func TestNew_Valid(t *testing.T) {
	f, err := invertfilter.New(&stubKeeper{"ERROR"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Enabled() {
		t.Fatal("expected filter to be enabled")
	}
}

func TestKeep_InvertMatchingLine(t *testing.T) {
	inner := &stubKeeper{sub: "ERROR"}
	f, _ := invertfilter.New(inner, true)

	if f.Keep("2024-01-01 ERROR something broke") {
		t.Error("expected Keep=false for matching line when inverted")
	}
	if !f.Keep("2024-01-01 INFO all good") {
		t.Error("expected Keep=true for non-matching line when inverted")
	}
}

func TestKeep_DisabledPassesAll(t *testing.T) {
	inner := &stubKeeper{sub: "ERROR"}
	f, _ := invertfilter.New(inner, false)

	lines := []string{
		"ERROR critical failure",
		"INFO startup complete",
		"",
	}
	for _, l := range lines {
		if !f.Keep(l) {
			t.Errorf("disabled filter should keep all lines, but rejected %q", l)
		}
	}
}

func TestKeep_EmptyLine(t *testing.T) {
	inner := &stubKeeper{sub: "ERROR"}
	f, _ := invertfilter.New(inner, true)

	// stubKeeper returns false for empty string; inversion → true
	if !f.Keep("") {
		t.Error("expected Keep=true for empty line (inner returns false, inverted to true)")
	}
}

func TestIntegration_FilterOutDebugLines(t *testing.T) {
	inner := &stubKeeper{sub: "DEBUG"}
	f, _ := invertfilter.New(inner, true)

	input := []string{
		"INFO  server started",
		"DEBUG entering handler",
		"WARN  disk usage high",
		"DEBUG leaving handler",
		"ERROR connection refused",
	}

	var kept []string
	for _, l := range input {
		if f.Keep(l) {
			kept = append(kept, l)
		}
	}

	if len(kept) != 3 {
		t.Fatalf("expected 3 lines after filtering DEBUG, got %d: %v", len(kept), kept)
	}
	for _, l := range kept {
		if containsSubstr(l, "DEBUG") {
			t.Errorf("DEBUG line leaked through: %q", l)
		}
	}
}
