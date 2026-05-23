package ratelimiter

import (
	"testing"
	"time"
)

func TestNew_ValidRate(t *testing.T) {
	rl, err := New(100)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if rl.Rate() != 100 {
		t.Errorf("expected rate 100, got %d", rl.Rate())
	}
}

func TestNew_ZeroRate(t *testing.T) {
	rl, err := New(0)
	if err != nil {
		t.Fatalf("expected no error for zero rate, got %v", err)
	}
	if rl.Rate() != 0 {
		t.Errorf("expected rate 0, got %d", rl.Rate())
	}
}

func TestNew_NegativeRate(t *testing.T) {
	_, err := New(-1)
	if err == nil {
		t.Fatal("expected error for negative rate, got nil")
	}
	if err != ErrInvalidRate {
		t.Errorf("expected ErrInvalidRate, got %v", err)
	}
}

func TestWait_Unlimited(t *testing.T) {
	rl, _ := New(0)
	slept := false
	rl.sleep = func(d time.Duration) { slept = true }
	for i := 0; i < 10000; i++ {
		rl.Wait()
	}
	if slept {
		t.Error("expected no sleep for unlimited rate")
	}
}

func TestWait_ThrottlesWhenLimitExceeded(t *testing.T) {
	rl, _ := New(5)

	var sleptDuration time.Duration
	rl.sleep = func(d time.Duration) { sleptDuration = d }

	fixedTime := time.Now()
	rl.now = func() time.Time { return fixedTime }
	rl.windowStart = fixedTime

	// First 5 calls should not sleep
	for i := 0; i < 5; i++ {
		rl.Wait()
	}
	if sleptDuration != 0 {
		t.Errorf("expected no sleep for first 5 calls, got %v", sleptDuration)
	}

	// 6th call should trigger sleep
	rl.Wait()
	if sleptDuration == 0 {
		t.Error("expected sleep on 6th call exceeding rate limit")
	}
}

func TestWait_ResetsAfterWindow(t *testing.T) {
	rl, _ := New(3)
	slept := false
	rl.sleep = func(d time.Duration) { slept = true }

	base := time.Now()
	call := 0
	rl.now = func() time.Time {
		call++
		// After 3 calls, simulate 1s has passed
		if call > 3 {
			return base.Add(time.Second + time.Millisecond)
		}
		return base
	}
	rl.windowStart = base

	for i := 0; i < 3; i++ {
		rl.Wait()
	}
	// Next call is in a new window — should not sleep
	rl.Wait()
	if slept {
		t.Error("expected no sleep after window reset")
	}
}

func TestCount_TracksCurrentWindow(t *testing.T) {
	rl, _ := New(100)
	rl.sleep = func(d time.Duration) {}

	for i := 0; i < 10; i++ {
		rl.Wait()
	}
	if rl.Count() != 10 {
		t.Errorf("expected count 10, got %d", rl.Count())
	}
}
