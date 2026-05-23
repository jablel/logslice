package linecount_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/linecount"
)

func TestNew(t *testing.T) {
	c := linecount.New()
	if c.Total() != 0 {
		t.Fatalf("expected Total=0, got %d", c.Total())
	}
	if c.Matched() != 0 {
		t.Fatalf("expected Matched=0, got %d", c.Matched())
	}
}

func TestInc(t *testing.T) {
	c := linecount.New()
	c.Inc()
	c.Inc()
	if c.Total() != 2 {
		t.Fatalf("expected Total=2, got %d", c.Total())
	}
	if c.Matched() != 0 {
		t.Fatalf("expected Matched=0, got %d", c.Matched())
	}
}

func TestMatch(t *testing.T) {
	c := linecount.New()
	c.Inc()
	c.Match()
	c.Match()
	if c.Total() != 3 {
		t.Fatalf("expected Total=3, got %d", c.Total())
	}
	if c.Matched() != 2 {
		t.Fatalf("expected Matched=2, got %d", c.Matched())
	}
}

func TestSkipped(t *testing.T) {
	c := linecount.New()
	c.Inc()
	c.Inc()
	c.Match()
	if c.Skipped() != 2 {
		t.Fatalf("expected Skipped=2, got %d", c.Skipped())
	}
}

func TestMatchRate_Zero(t *testing.T) {
	c := linecount.New()
	if c.MatchRate() != 0 {
		t.Fatalf("expected MatchRate=0 with no lines, got %f", c.MatchRate())
	}
}

func TestMatchRate(t *testing.T) {
	c := linecount.New()
	c.Match()
	c.Match()
	c.Inc()
	c.Inc()
	rate := c.MatchRate()
	if rate < 0.49 || rate > 0.51 {
		t.Fatalf("expected MatchRate≈0.5, got %f", rate)
	}
}

func TestReset(t *testing.T) {
	c := linecount.New()
	c.Match()
	c.Inc()
	c.Reset()
	if c.Total() != 0 || c.Matched() != 0 {
		t.Fatalf("expected all zeros after Reset, got total=%d matched=%d",
			c.Total(), c.Matched())
	}
}
