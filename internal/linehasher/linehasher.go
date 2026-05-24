// Package linehasher computes a short hash digest for each log line,
// optionally prepending or appending it as an annotation. This is useful
// for deduplication keys, integrity checks, or log correlation.
package linehasher

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
)

// Algorithm selects the hash function to use.
type Algorithm string

const (
	MD5    Algorithm = "md5"
	SHA256 Algorithm = "sha256"
)

// Position controls where the hash is placed relative to the line.
type Position string

const (
	Prepend Position = "prepend"
	Append  Position = "append"
)

// Hasher annotates each line with a truncated hash digest.
type Hasher struct {
	algo   Algorithm
	pos    Position
	len    int // number of hex chars to keep
	newFn  func() hash.Hash
}

// New creates a Hasher. algo must be "md5" or "sha256". pos must be
// "prepend" or "append". digestLen is the number of hex characters to
// include (1–64); values outside that range are clamped.
func New(algo Algorithm, pos Position, digestLen int) (*Hasher, error) {
	var newFn func() hash.Hash
	switch algo {
	case MD5:
		newFn = func() hash.Hash { return md5.New() }
	case SHA256:
		newFn = func() hash.Hash { return sha256.New() }
	default:
		return nil, errors.New("linehasher: unknown algorithm: " + string(algo))
	}
	if pos != Prepend && pos != Append {
		return nil, errors.New("linehasher: position must be \"prepend\" or \"append\"")
	}
	if digestLen < 1 {
		digestLen = 1
	}
	if digestLen > 64 {
		digestLen = 64
	}
	return &Hasher{algo: algo, pos: pos, len: digestLen, newFn: newFn}, nil
}

// Apply returns the line annotated with its hash digest.
func (h *Hasher) Apply(line string) string {
	hw := h.newFn()
	_, _ = fmt.Fprint(hw, line)
	full := hex.EncodeToString(hw.Sum(nil))
	digest := full
	if h.len < len(full) {
		digest = full[:h.len]
	}
	if h.pos == Prepend {
		return "[" + digest + "] " + line
	}
	return line + " [" + digest + "]"
}
