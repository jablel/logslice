// Package lineencoder encodes log lines into alternate text representations
// such as base64 or URL encoding. Useful for safely transporting lines that
// may contain special characters or binary data.
package lineencoder

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// Encoding selects the output encoding format.
type Encoding int

const (
	Base64  Encoding = iota // standard base64 encoding
	URLEncode               // percent-encoded URL encoding
	Hex                     // lowercase hexadecimal encoding
)

// Encoder transforms each log line into the chosen encoding.
type Encoder struct {
	enc     Encoding
	enabled bool
}

// New creates an Encoder for the given encoding format.
// Returns an error for unknown encoding values.
func New(enc Encoding, enabled bool) (*Encoder, error) {
	if enc != Base64 && enc != URLEncode && enc != Hex {
		return nil, fmt.Errorf("lineencoder: unknown encoding %d", enc)
	}
	return &Encoder{enc: enc, enabled: enabled}, nil
}

// Apply encodes line according to the configured encoding.
// If the encoder is disabled, line is returned unchanged.
func (e *Encoder) Apply(line string) string {
	if !e.enabled {
		return line
	}
	switch e.enc {
	case Base64:
		return base64.StdEncoding.EncodeToString([]byte(line))
	case URLEncode:
		return url.QueryEscape(line)
	case Hex:
		var sb strings.Builder
		for _, b := range []byte(line) {
			fmt.Fprintf(&sb, "%02x", b)
		}
		return sb.String()
	}
	return line
}
