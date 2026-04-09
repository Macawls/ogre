package font

import (
	"bytes"

	"github.com/go-text/typesetting/di"
	gotextfont "github.com/go-text/typesetting/font"
	ot "github.com/go-text/typesetting/font/opentype"
	"github.com/go-text/typesetting/shaping"
	"golang.org/x/image/math/fixed"
)

type ShapedGlyph struct {
	GlyphID  uint32
	Advance  float64
	XOffset  float64
	YOffset  float64
	XBearing float64
	Width    float64
}

type ShapedRun struct {
	Glyphs  []ShapedGlyph
	Advance float64
}

type Shaper struct {
	shaper shaping.HarfbuzzShaper
}

func NewShaper() *Shaper {
	return &Shaper{}
}

func (s *Shaper) ShapeBytes(text string, fontData []byte, size float64, rtl bool) (*ShapedRun, error) {
	reader := bytes.NewReader(fontData)
	ld, err := ot.NewLoader(reader)
	if err != nil {
		return nil, err
	}

	ft, err := gotextfont.NewFont(ld)
	if err != nil {
		return nil, err
	}

	face := gotextfont.NewFace(ft)

	return s.shape(text, face, size, rtl), nil
}

func (s *Shaper) shape(text string, face *gotextfont.Face, size float64, rtl bool) *ShapedRun {
	runes := []rune(text)
	input := shaping.Input{
		Text:     runes,
		RunStart: 0,
		RunEnd:   len(runes),
		Face:     face,
		Size:     fixed.I(int(size)),
	}

	if rtl {
		input.Direction = di.DirectionRTL
	}

	output := s.shaper.Shape(input)

	run := &ShapedRun{
		Advance: float64(output.Advance) / 64.0,
	}

	for _, g := range output.Glyphs {
		run.Glyphs = append(run.Glyphs, ShapedGlyph{
			GlyphID:  uint32(g.GlyphID),
			Advance:  float64(g.Advance) / 64.0,
			XOffset:  float64(g.XOffset) / 64.0,
			YOffset:  float64(g.YOffset) / 64.0,
			XBearing: float64(g.XBearing) / 64.0,
			Width:    float64(g.Width) / 64.0,
		})
	}

	return run
}
