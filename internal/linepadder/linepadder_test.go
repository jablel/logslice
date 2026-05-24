package linepadder

import (
	"strings"
	"testing"
)

func TestNew_ValidLeft(t *testing.T) {
	p, err := New(10, AlignLeft, ' ')
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Width() != 10 {
		t.Errorf("expected width 10, got %d", p.Width())
	}
}

func TestNew_InvalidWidth(t *testing.T) {
	_, err := New(0, AlignLeft, ' ')
	if err == nil {
		t.Fatal("expected error for zero width")
	}
}

func TestNew_ZeroFill(t *testing.T) {
	_, err := New(5, AlignLeft, 0)
	if err == nil {
		t.Fatal("expected error for zero fill rune")
	}
}

func TestNew_UnknownAlignment(t *testing.T) {
	_, err := New(5, Alignment(99), ' ')
	if err == nil {
		t.Fatal("expected error for unknown alignment")
	}
}

func TestApply_Left_ShortLine(t *testing.T) {
	p, _ := New(10, AlignLeft, '-')
	got := p.Apply("hello")
	if got != "hello-----" {
		t.Errorf("got %q", got)
	}
}

func TestApply_Right_ShortLine(t *testing.T) {
	p, _ := New(10, AlignRight, '-')
	got := p.Apply("hello")
	if got != "-----hello" {
		t.Errorf("got %q", got)
	}
}

func TestApply_ExactWidth(t *testing.T) {
	p, _ := New(5, AlignLeft, ' ')
	got := p.Apply("hello")
	if got != "hello" {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestApply_Truncates_LongLine(t *testing.T) {
	p, _ := New(4, AlignLeft, ' ')
	got := p.Apply("toolongline")
	if got != "tool" {
		t.Errorf("expected truncation to 4 chars, got %q", got)
	}
}

func TestApply_Disabled_ReturnsOriginal(t *testing.T) {
	p, _ := New(10, AlignLeft, ' ')
	p.SetEnabled(false)
	line := "unchanged"
	if got := p.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_Unicode(t *testing.T) {
	p, _ := New(5, AlignLeft, '·')
	got := p.Apply("日本")
	runes := []rune(got)
	if len(runes) != 5 {
		t.Errorf("expected 5 runes, got %d: %q", len(runes), got)
	}
	if !strings.HasPrefix(got, "日本") {
		t.Errorf("expected prefix 日本, got %q", got)
	}
}
