package sampler_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/sampler"
)

func TestNew_ValidStep(t *testing.T) {
	s, err := sampler.New(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Step() != 3 {
		t.Errorf("expected step 3, got %d", s.Step())
	}
}

func TestNew_StepOne(t *testing.T) {
	s, err := sampler.New(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// step=1 should keep every line
	for i := 0; i < 5; i++ {
		if !s.Keep() {
			t.Errorf("step=1: expected Keep()=true on call %d", i+1)
		}
	}
}

func TestNew_InvalidStep(t *testing.T) {
	for _, bad := range []int{0, -1, -100} {
		_, err := sampler.New(bad)
		if err == nil {
			t.Errorf("expected error for step=%d", bad)
		}
	}
}

func TestKeep_EveryNth(t *testing.T) {
	s, _ := sampler.New(3)
	results := make([]bool, 9)
	for i := range results {
		results[i] = s.Keep()
	}
	// With step=3, indices 2,5,8 (0-based) should be true
	expected := []bool{false, false, true, false, false, true, false, false, true}
	for i, got := range results {
		if got != expected[i] {
			t.Errorf("index %d: expected %v, got %v", i, expected[i], got)
		}
	}
}

func TestReset(t *testing.T) {
	s, _ := sampler.New(3)
	s.Keep() // counter=1
	s.Keep() // counter=2
	s.Reset()
	// After reset, behaviour should restart
	if s.Keep() { // counter=1, not yet at step
		t.Error("expected false after reset on first Keep()")
	}
	if s.Keep() { // counter=2
		t.Error("expected false on second Keep() after reset")
	}
	if !s.Keep() { // counter=3 => keep
		t.Error("expected true on third Keep() after reset")
	}
}
