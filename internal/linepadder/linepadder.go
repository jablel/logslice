package linepadder

import (
	"fmt"
	"strings"
)

// Padder pads or truncates each line to a fixed width using a configurable
// alignment (left or right) and fill character.
type Padder struct {
	width  int
	align  Alignment
	fill   rune
	enabled bool
}

// Alignment controls how lines are padded.
type Alignment int

const (
	AlignLeft  Alignment = iota
	AlignRight
)

// New creates a Padder that pads lines to width characters.
// fill is the character used for padding (e.g. ' ' or '0').
// Returns an error if width < 1 or fill is zero.
func New(width int, align Alignment, fill rune) (*Padder, error) {
	if width < 1 {
		return nil, fmt.Errorf("linepadder: width must be >= 1, got %d", width)
	}
	if fill == 0 {
		return nil, fmt.Errorf("linepadder: fill rune must not be zero")
	}
	if align != AlignLeft && align != AlignRight {
		return nil, fmt.Errorf("linepadder: unknown alignment %d", align)
	}
	return &Padder{width: width, align: align, fill: fill, enabled: true}, nil
}

// Apply pads or truncates line to the configured width.
// If the padder is disabled it returns line unchanged.
func (p *Padder) Apply(line string) string {
	if !p.enabled {
		return line
	}
	runes := []rune(line)
	if len(runes) >= p.width {
		return string(runes[:p.width])
	}
	padLen := p.width - len(runes)
	padding := strings.Repeat(string(p.fill), padLen)
	if p.align == AlignLeft {
		return line + padding
	}
	return padding + line
}

// SetEnabled toggles the padder on or off.
func (p *Padder) SetEnabled(v bool) { p.enabled = v }

// Width returns the configured target width.
func (p *Padder) Width() int { return p.width }
