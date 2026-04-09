package font

import (
	"image"
	"testing"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type mockFace struct {
	advance fixed.Int26_6
	ascent  fixed.Int26_6
	descent fixed.Int26_6
}

func (m *mockFace) Close() error { return nil }

func (m *mockFace) Glyph(dot fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
	return image.Rectangle{}, nil, image.Point{}, m.advance, true
}

func (m *mockFace) GlyphBounds(r rune) (fixed.Rectangle26_6, fixed.Int26_6, bool) {
	return fixed.Rectangle26_6{}, m.advance, true
}

func (m *mockFace) GlyphAdvance(r rune) (fixed.Int26_6, bool) {
	return m.advance, true
}

func (m *mockFace) Kern(r0, r1 rune) fixed.Int26_6 {
	return 0
}

func (m *mockFace) Metrics() font.Metrics {
	return font.Metrics{
		Ascent:  m.ascent,
		Descent: m.descent,
		Height:  m.ascent + m.descent,
	}
}

func newMockFace(advancePx, ascentPx, descentPx float64) *mockFace {
	return &mockFace{
		advance: fixed.I(int(advancePx)),
		ascent:  fixed.I(int(ascentPx)),
		descent: fixed.I(int(descentPx)),
	}
}

func TestMeasureString(t *testing.T) {
	face := newMockFace(10, 12, 4)
	got := MeasureString(face, "hello")
	want := 50.0
	if got != want {
		t.Errorf("MeasureString = %v, want %v", got, want)
	}
}

func TestMeasureStringWithSpacing(t *testing.T) {
	face := newMockFace(10, 12, 4)

	got := MeasureStringWithSpacing(face, "hello", 2.0)
	want := 58.0 // 50 + 2*4
	if got != want {
		t.Errorf("MeasureStringWithSpacing = %v, want %v", got, want)
	}

	got = MeasureStringWithSpacing(face, "a", 5.0)
	want = 10.0
	if got != want {
		t.Errorf("MeasureStringWithSpacing single char = %v, want %v", got, want)
	}

	got = MeasureStringWithSpacing(face, "", 5.0)
	want = 0.0
	if got != want {
		t.Errorf("MeasureStringWithSpacing empty = %v, want %v", got, want)
	}
}

func TestLineHeight(t *testing.T) {
	face := newMockFace(10, 12, 4)
	got := LineHeight(face)
	want := 16.0
	if got != want {
		t.Errorf("LineHeight = %v, want %v", got, want)
	}
}

func TestAscent(t *testing.T) {
	face := newMockFace(10, 12, 4)
	got := Ascent(face)
	want := 12.0
	if got != want {
		t.Errorf("Ascent = %v, want %v", got, want)
	}
}

func TestDescent(t *testing.T) {
	face := newMockFace(10, 12, 4)
	got := Descent(face)
	want := 4.0
	if got != want {
		t.Errorf("Descent = %v, want %v", got, want)
	}
}

func TestMeasurer_StringWidth(t *testing.T) {
	face := newMockFace(10, 12, 4)
	m := NewMeasurer(face, 2.0)

	got := m.StringWidth("hello")
	want := 58.0
	if got != want {
		t.Errorf("StringWidth = %v, want %v", got, want)
	}
}

func TestMeasurer_RuneWidth(t *testing.T) {
	face := newMockFace(10, 12, 4)
	m := NewMeasurer(face, 0)

	got := m.RuneWidth('a')
	want := 10.0
	if got != want {
		t.Errorf("RuneWidth = %v, want %v", got, want)
	}

	got2 := m.RuneWidth('a')
	if got2 != got {
		t.Errorf("cached RuneWidth = %v, want %v", got2, got)
	}
}

func TestMeasurer_EmptyString(t *testing.T) {
	face := newMockFace(10, 12, 4)
	m := NewMeasurer(face, 5.0)

	got := m.StringWidth("")
	want := 0.0
	if got != want {
		t.Errorf("StringWidth empty = %v, want %v", got, want)
	}
}
