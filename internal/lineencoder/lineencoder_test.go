package lineencoder

import (
	"encoding/base64"
	"net/url"
	"testing"
)

func TestNew_ValidBase64(t *testing.T) {
	e, err := New(Base64, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil encoder")
	}
}

func TestNew_UnknownEncoding(t *testing.T) {
	_, err := New(Encoding(99), true)
	if err == nil {
		t.Fatal("expected error for unknown encoding")
	}
}

func TestApply_Disabled_ReturnsOriginal(t *testing.T) {
	e, _ := New(Base64, false)
	got := e.Apply("hello world")
	if got != "hello world" {
		t.Errorf("expected original line, got %q", got)
	}
}

func TestApply_Base64(t *testing.T) {
	e, _ := New(Base64, true)
	input := "2024-01-01 ERROR something went wrong"
	got := e.Apply(input)
	want := base64.StdEncoding.EncodeToString([]byte(input))
	if got != want {
		t.Errorf("base64: got %q, want %q", got, want)
	}
}

func TestApply_URLEncode(t *testing.T) {
	e, _ := New(URLEncode, true)
	input := "key=value&foo=bar baz"
	got := e.Apply(input)
	want := url.QueryEscape(input)
	if got != want {
		t.Errorf("url: got %q, want %q", got, want)
	}
}

func TestApply_Hex(t *testing.T) {
	e, _ := New(Hex, true)
	got := e.Apply("AB")
	if got != "4142" {
		t.Errorf("hex: got %q, want %q", got, "4142")
	}
}

func TestApply_EmptyLine(t *testing.T) {
	for _, enc := range []Encoding{Base64, URLEncode, Hex} {
		e, _ := New(enc, true)
		got := e.Apply("")
		if got != "" {
			t.Errorf("enc %d: expected empty string for empty input, got %q", enc, got)
		}
	}
}

func TestApply_Hex_RoundTrip(t *testing.T) {
	e, _ := New(Hex, true)
	input := "logslice"
	got := e.Apply(input)
	if len(got) != len(input)*2 {
		t.Errorf("hex length: got %d, want %d", len(got), len(input)*2)
	}
}
