package linehasher

import (
	"strings"
	"testing"
)

func TestNew_ValidMD5Prepend(t *testing.T) {
	h, err := New(MD5, Prepend, 8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil Hasher")
	}
}

func TestNew_ValidSHA256Append(t *testing.T) {
	h, err := New(SHA256, Append, 16)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil Hasher")
	}
}

func TestNew_UnknownAlgorithm(t *testing.T) {
	_, err := New("crc32", Prepend, 8)
	if err == nil {
		t.Fatal("expected error for unknown algorithm")
	}
}

func TestNew_UnknownPosition(t *testing.T) {
	_, err := New(MD5, "middle", 8)
	if err == nil {
		t.Fatal("expected error for unknown position")
	}
}

func TestNew_DigestLenClamped(t *testing.T) {
	h, err := New(MD5, Prepend, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := h.Apply("hello")
	// digest of length 1 should still produce a bracket annotation
	if !strings.HasPrefix(out, "[") {
		t.Errorf("expected prepend annotation, got: %q", out)
	}
}

func TestApply_Prepend_Format(t *testing.T) {
	h, _ := New(MD5, Prepend, 8)
	line := "2024-01-15T10:00:00Z level=info msg=\"started\""
	out := h.Apply(line)
	if !strings.HasPrefix(out, "[") {
		t.Errorf("expected leading '[', got: %q", out)
	}
	if !strings.Contains(out, line) {
		t.Errorf("original line not found in output: %q", out)
	}
}

func TestApply_Append_Format(t *testing.T) {
	h, _ := New(SHA256, Append, 12)
	line := "2024-01-15T10:00:00Z level=error msg=\"disk full\""
	out := h.Apply(line)
	if !strings.HasSuffix(out, "]") {
		t.Errorf("expected trailing ']', got: %q", out)
	}
	if !strings.HasPrefix(out, line) {
		t.Errorf("original line should be prefix in append mode: %q", out)
	}
}

func TestApply_Deterministic(t *testing.T) {
	h, _ := New(SHA256, Prepend, 16)
	line := "some log line content"
	out1 := h.Apply(line)
	out2 := h.Apply(line)
	if out1 != out2 {
		t.Errorf("Apply is not deterministic: %q vs %q", out1, out2)
	}
}

func TestApply_DifferentLines_DifferentDigests(t *testing.T) {
	h, _ := New(MD5, Prepend, 32)
	out1 := h.Apply("line one")
	out2 := h.Apply("line two")
	if out1 == out2 {
		t.Error("different lines produced the same digest annotation")
	}
}

func TestApply_DigestLength(t *testing.T) {
	for _, n := range []int{4, 8, 16} {
		h, _ := New(MD5, Prepend, n)
		out := h.Apply("test")
		// format: "[<n hex chars>] test"
		start := strings.Index(out, "[")
		end := strings.Index(out, "]")
		if start < 0 || end < 0 {
			t.Fatalf("n=%d: bracket not found in %q", n, out)
		}
		digest := out[start+1 : end]
		if len(digest) != n {
			t.Errorf("n=%d: expected digest length %d, got %d (%q)", n, n, len(digest), digest)
		}
	}
}
