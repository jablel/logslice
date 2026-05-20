package timeparser_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/timeparser"
)

func TestParse(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
	}{
		{"2024-03-15T10:30:00Z", false},
		{"2024-03-15T10:30:00.123456789Z", false},
		{"2024-03-15T10:30:00", false},
		{"2024-03-15 10:30:00", false},
		{"2024-03-15 10:30:00.123", false},
		{"not-a-timestamp", true},
		{"", true},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, _, err := timeparser.Parse(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error for input %q, got time %v", tc.input, got)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %q: %v", tc.input, err)
				}
			}
		})
	}
}

func TestParseWithFormat(t *testing.T) {
	_, err := timeparser.ParseWithFormat("2024-03-15 10:30:00", "2006-01-02 15:04:05")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = timeparser.ParseWithFormat("garbage", "2006-01-02")
	if err == nil {
		t.Fatal("expected error for invalid input")
	}
}

func TestInRange(t *testing.T) {
	base := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	start := base
	end := base.Add(1 * time.Hour)

	cases := []struct {
		name string
		t    time.Time
		want bool
	}{
		{"at start", start, true},
		{"at end", end, true},
		{"in middle", base.Add(30 * time.Minute), true},
		{"before start", base.Add(-1 * time.Second), false},
		{"after end", end.Add(1 * time.Second), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := timeparser.InRange(tc.t, start, end); got != tc.want {
				t.Errorf("InRange(%v) = %v, want %v", tc.t, got, tc.want)
			}
		})
	}
}

func TestKnownFormats(t *testing.T) {
	fmts := timeparser.KnownFormats()
	if len(fmts) == 0 {
		t.Fatal("expected at least one known format")
	}
}
