package columnextractor_test

import (
	"testing"

	"github.com/user/logslice/internal/columnextractor"
)

var syslogLines = []string{
	"Jan 15 12:00:01 myhost sshd[1234]: Accepted publickey for user",
	"Jan 15 12:00:02 myhost kernel: eth0: renamed from veth3a2b",
	"Jan 15 12:00:03 myhost systemd[1]: Started Daily apt upgrade.",
}

func TestIntegration_ExtractHostname(t *testing.T) {
	ex, err := columnextractor.New(" ", 3)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	expected := []string{"myhost", "myhost", "myhost"}
	for i, line := range syslogLines {
		val, ok := ex.Extract(line)
		if !ok {
			t.Errorf("line %d: expected ok=true", i)
			continue
		}
		if val != expected[i] {
			t.Errorf("line %d: expected %q, got %q", i, expected[i], val)
		}
	}
}

func TestIntegration_ExtractProcess(t *testing.T) {
	ex, err := columnextractor.New(" ", 4)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	expected := []string{"sshd[1234]:", "kernel:", "systemd[1]:"}
	for i, line := range syslogLines {
		val, ok := ex.Extract(line)
		if !ok {
			t.Errorf("line %d: expected ok=true", i)
			continue
		}
		if val != expected[i] {
			t.Errorf("line %d: expected %q, got %q", i, expected[i], val)
		}
	}
}

func TestIntegration_SkipShortLines(t *testing.T) {
	lines := []string{
		"full line with many columns here",
		"short",
		"also a full line with columns",
	}

	ex, _ := columnextractor.New(" ", 5)
	var kept []string
	for _, line := range lines {
		if _, ok := ex.Extract(line); ok {
			kept = append(kept, line)
		}
	}
	if len(kept) != 2 {
		t.Errorf("expected 2 kept lines, got %d", len(kept))
	}
}
