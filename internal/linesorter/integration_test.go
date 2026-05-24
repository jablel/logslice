package linesorter_test

import (
	"strings"
	"testing"

	"logslice/internal/linesorter"
)

func feedAll(s *linesorter.Sorter, lines []string) {
	for _, l := range lines {
		s.Feed(l)
	}
}

func TestIntegration_SortLogLines_Ascending(t *testing.T) {
	input := []string{
		"2024-01-03 ERROR disk full",
		"2024-01-01 INFO  started",
		"2024-01-02 WARN  high memory",
	}

	s, err := linesorter.New(linesorter.Ascending)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	feedAll(s, input)
	got := s.Flush()

	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if !strings.HasPrefix(got[0], "2024-01-01") {
		t.Errorf("first line should start with 2024-01-01, got %q", got[0])
	}
	if !strings.HasPrefix(got[2], "2024-01-03") {
		t.Errorf("last line should start with 2024-01-03, got %q", got[2])
	}
}

func TestIntegration_SortLogLines_Descending(t *testing.T) {
	input := []string{
		"2024-01-01 INFO  started",
		"2024-01-03 ERROR disk full",
		"2024-01-02 WARN  high memory",
	}

	s, err := linesorter.New(linesorter.Descending)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	feedAll(s, input)
	got := s.Flush()

	if !strings.HasPrefix(got[0], "2024-01-03") {
		t.Errorf("first line should start with 2024-01-03, got %q", got[0])
	}
}

func TestIntegration_MultipleFlushCycles(t *testing.T) {
	s, _ := linesorter.New(linesorter.Ascending)

	s.Feed("c")
	s.Feed("a")
	first := s.Flush()
	if first[0] != "a" {
		t.Errorf("cycle 1: want \"a\", got %q", first[0])
	}

	s.Feed("z")
	s.Feed("m")
	second := s.Flush()
	if second[0] != "m" {
		t.Errorf("cycle 2: want \"m\", got %q", second[0])
	}
	if s.Len() != 0 {
		t.Errorf("buffer should be empty after second flush")
	}
}
