package font

import (
	"golang.org/x/image/font"
)

// MeasureString returns the advance width of the text using the given font face.
// MeasureString returns the advance width of a string in pixels.
func MeasureString(face font.Face, text string) float64 {
	advance := font.MeasureString(face, text)
	return float64(advance) / 64.0
}

// MeasureStringWithSpacing returns the advance width of the text including letter spacing.
func MeasureStringWithSpacing(face font.Face, text string, letterSpacing float64) float64 {
	width := MeasureString(face, text)
	count := 0
	for range text {
		count++
	}
	if count > 1 {
		width += letterSpacing * float64(count-1)
	}
	return width
}

// LineHeight returns the sum of ascent and descent for the given font face.
func LineHeight(face font.Face) float64 {
	metrics := face.Metrics()
	ascent := float64(metrics.Ascent) / 64.0
	descent := float64(metrics.Descent) / 64.0
	return ascent + descent
}

// Ascent returns the ascent metric of the font face in pixels.
func Ascent(face font.Face) float64 {
	return float64(face.Metrics().Ascent) / 64.0
}

// Descent returns the descent metric of the font face in pixels.
func Descent(face font.Face) float64 {
	return float64(face.Metrics().Descent) / 64.0
}

// Measurer caches per-rune widths for efficient text measurement.
type Measurer struct {
	face    font.Face
	spacing float64
	cache   map[rune]float64
}

// NewMeasurer creates a Measurer for the given font face and letter spacing.
func NewMeasurer(face font.Face, letterSpacing float64) *Measurer {
	return &Measurer{
		face:    face,
		spacing: letterSpacing,
		cache:   make(map[rune]float64),
	}
}

func (m *Measurer) RuneWidth(r rune) float64 {
	if w, ok := m.cache[r]; ok {
		return w
	}
	adv, ok := m.face.GlyphAdvance(r)
	if !ok {
		return 0
	}
	w := float64(adv) / 64.0
	m.cache[r] = w
	return w
}

func (m *Measurer) StringWidth(s string) float64 {
	var total float64
	count := 0
	for _, r := range s {
		total += m.RuneWidth(r)
		count++
	}
	if count > 1 {
		total += m.spacing * float64(count-1)
	}
	return total
}
