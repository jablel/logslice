package timestamprangefilter_test

import (
	"testing"
	"time"

	"logslice/internal/timestamprangefilter"
)

// stubExtractor implements Extractor for testing.
type stubExtractor struct {
	timestamps map[string]time.Time
}

func (s *stubExtractor) Extract(line string) (time.Time, bool) {
	ts, ok := s.timestamps[line]
	return ts, ok
}

var (
	t0 = time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	t1 = time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
	t2 = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	t3 = time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC)
)

func newStub() *stubExtractor {
	return &stubExtractor{
		timestamps: map[string]time.Time{
			"before": t0,
			"start":  t1,
			"middle": t2,
			"end":    t3,
		},
	}
}

func TestNew_NilExtractor(t *testing.T) {
	_, err := timestamprangefilter.New(nil, t1, t3)
	if err == nil {
		t.Fatal("expected error for nil extractor")
	}
}

func TestNew_ToBeforeFrom(t *testing.T) {
	_, err := timestamprangefilter.New(newStub(), t3, t1)
	if err == nil {
		t.Fatal("expected error when to < from")
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := timestamprangefilter.New(newStub(), t1, t3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestKeep_InRange(t *testing.T) {
	f, _ := timestamprangefilter.New(newStub(), t1, t3)
	for _, line := range []string{"start", "middle", "end"} {
		if !f.Keep(line) {
			t.Errorf("expected %q to be kept", line)
		}
	}
}

func TestKeep_BeforeRange(t *testing.T) {
	f, _ := timestamprangefilter.New(newStub(), t1, t3)
	if f.Keep("before") {
		t.Error("expected 'before' to be dropped")
	}
}

func TestKeep_NoTimestamp_PassThrough(t *testing.T) {
	f, _ := timestamprangefilter.New(newStub(), t1, t3)
	if !f.Keep("no-timestamp-line") {
		t.Error("expected line without timestamp to pass through by default")
	}
}

func TestKeep_NoTimestamp_Dropped(t *testing.T) {
	f, _ := timestamprangefilter.New(newStub(), t1, t3,
		timestamprangefilter.PassThroughNoTimestamp(false))
	if f.Keep("no-timestamp-line") {
		t.Error("expected line without timestamp to be dropped")
	}
}

func TestKeep_UnboundedTo(t *testing.T) {
	f, _ := timestamprangefilter.New(newStub(), t1, time.Time{})
	for _, line := range []string{"start", "middle", "end"} {
		if !f.Keep(line) {
			t.Errorf("expected %q to be kept with unbounded upper range", line)
		}
	}
	if f.Keep("before") {
		t.Error("expected 'before' to be dropped even with unbounded upper range")
	}
}
