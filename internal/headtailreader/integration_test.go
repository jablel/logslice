package headtailreader_test

import (
	"fmt"
	"testing"

	"logslice/internal/headtailreader"
)

func generateLines(n int) []string {
	lines := make([]string, n)
	for i := range lines {
		lines[i] = fmt.Sprintf("line-%04d", i+1)
	}
	return lines
}

func TestIntegration_Head100From1000(t *testing.T) {
	r, err := headtailreader.New("head", 100)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range generateLines(1000) {
		r.Feed(l)
	}
	got := r.Lines()
	if len(got) != 100 {
		t.Fatalf("want 100 lines, got %d", len(got))
	}
	if got[0] != "line-0001" {
		t.Errorf("first line: want line-0001, got %q", got[0])
	}
	if got[99] != "line-0100" {
		t.Errorf("last line: want line-0100, got %q", got[99])
	}
}

func TestIntegration_Tail100From1000(t *testing.T) {
	r, err := headtailreader.New("tail", 100)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range generateLines(1000) {
		r.Feed(l)
	}
	got := r.Lines()
	if len(got) != 100 {
		t.Fatalf("want 100 lines, got %d", len(got))
	}
	if got[0] != "line-0901" {
		t.Errorf("first tail line: want line-0901, got %q", got[0])
	}
	if got[99] != "line-1000" {
		t.Errorf("last tail line: want line-1000, got %q", got[99])
	}
}

func TestIntegration_ReuseAfterReset(t *testing.T) {
	r, _ := headtailreader.New("head", 3)
	for _, l := range []string{"a", "b", "c", "d"} {
		r.Feed(l)
	}
	r.Reset()
	for _, l := range []string{"x", "y"} {
		r.Feed(l)
	}
	got := r.Lines()
	if len(got) != 2 {
		t.Fatalf("after reset want 2 lines, got %d", len(got))
	}
	if got[0] != "x" || got[1] != "y" {
		t.Errorf("unexpected lines after reset: %v", got)
	}
}
