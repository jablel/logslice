package limiter

import (
	"testing"
)

func TestNew_PositiveLimit(t *testing.T) {
	l := New(3)
	if l.max != 3 {
		t.Fatalf("expected max=3, got %d", l.max)
	}
}

func TestKeep_RespectsLimit(t *testing.T) {
	l := New(3)
	for i := 0; i < 3; i++ {
		if !l.Keep() {
			t.Fatalf("expected Keep()=true on call %d", i+1)
		}
	}
	if l.Keep() {
		t.Fatal("expected Keep()=false after limit reached")
	}
}

func TestKeep_NoLimit(t *testing.T) {
	for _, max := range []int{0, -1, -100} {
		l := New(max)
		for i := 0; i < 1000; i++ {
			if !l.Keep() {
				t.Fatalf("max=%d: Keep() returned false at call %d", max, i+1)
			}
		}
	}
}

func TestCount(t *testing.T) {
	l := New(5)
	for i := 0; i < 5; i++ {
		l.Keep()
	}
	if l.Count() != 5 {
		t.Fatalf("expected count=5, got %d", l.Count())
	}
	// Extra calls beyond limit should not increment count.
	l.Keep()
	if l.Count() != 5 {
		t.Fatalf("expected count still 5 after limit, got %d", l.Count())
	}
}

func TestDone(t *testing.T) {
	l := New(2)
	if l.Done() {
		t.Fatal("expected Done()=false before any calls")
	}
	l.Keep()
	l.Keep()
	if !l.Done() {
		t.Fatal("expected Done()=true after limit reached")
	}
}

func TestDone_NoLimit(t *testing.T) {
	l := New(0)
	for i := 0; i < 100; i++ {
		l.Keep()
	}
	if l.Done() {
		t.Fatal("expected Done()=false when no limit set")
	}
}

func TestReset(t *testing.T) {
	l := New(2)
	l.Keep()
	l.Keep()
	if !l.Done() {
		t.Fatal("expected Done()=true before reset")
	}
	l.Reset()
	if l.Done() {
		t.Fatal("expected Done()=false after reset")
	}
	if l.Count() != 0 {
		t.Fatalf("expected count=0 after reset, got %d", l.Count())
	}
	if !l.Keep() {
		t.Fatal("expected Keep()=true after reset")
	}
}
