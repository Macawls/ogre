package font

import (
	"fmt"
	"strings"

	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// GlyphPath holds the SVG path data and advance width for a single glyph.
type GlyphPath struct {
	D       string
	Advance float64
}

// GlyphToPath converts a single rune to an SVG path at the given size.
// GlyphToPath extracts SVG path data for a single glyph.
func GlyphToPath(f *opentype.Font, r rune, size float64) (GlyphPath, error) {
	var buf sfnt.Buffer

	idx, err := f.GlyphIndex(&buf, r)
	if err != nil {
		return GlyphPath{}, fmt.Errorf("glyph index for %q: %w", string(r), err)
	}
	if idx == 0 {
		return GlyphPath{}, fmt.Errorf("no glyph for %q", string(r))
	}

	ppem := fixed.Int26_6(size * 64)

	segments, err := f.LoadGlyph(&buf, idx, ppem, nil)
	if err != nil {
		return GlyphPath{}, fmt.Errorf("load glyph for %q: %w", string(r), err)
	}

	advance, err := f.GlyphAdvance(&buf, idx, ppem, 0)
	if err != nil {
		return GlyphPath{}, fmt.Errorf("glyph advance for %q: %w", string(r), err)
	}

	var b strings.Builder
	for _, seg := range segments {
		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			fmt.Fprintf(&b, "M%.4g %.4g", fix(seg.Args[0].X), fix(seg.Args[0].Y))
		case sfnt.SegmentOpLineTo:
			fmt.Fprintf(&b, "L%.4g %.4g", fix(seg.Args[0].X), fix(seg.Args[0].Y))
		case sfnt.SegmentOpQuadTo:
			fmt.Fprintf(&b, "Q%.4g %.4g %.4g %.4g",
				fix(seg.Args[0].X), fix(seg.Args[0].Y),
				fix(seg.Args[1].X), fix(seg.Args[1].Y))
		case sfnt.SegmentOpCubeTo:
			fmt.Fprintf(&b, "C%.4g %.4g %.4g %.4g %.4g %.4g",
				fix(seg.Args[0].X), fix(seg.Args[0].Y),
				fix(seg.Args[1].X), fix(seg.Args[1].Y),
				fix(seg.Args[2].X), fix(seg.Args[2].Y))
		}
	}

	if b.Len() > 0 {
		b.WriteString("Z")
	}

	return GlyphPath{
		D:       b.String(),
		Advance: fix(advance),
	}, nil
}

// TextToPath converts a text string to an SVG path using the resolved font.
// TextToPath converts a text string to SVG path data using the font manager.
func ShapedTextToPath(mgr *Manager, text string, family string, weight int, style string, size float64, rtl bool) (string, float64) {
	face := mgr.Resolve(family, weight, style)
	if face == nil || len(face.RawData) == 0 {
		return TextToPath(mgr, text, family, weight, style, size)
	}

	shaper := NewShaper()
	run, err := shaper.ShapeBytes(text, face.RawData, size, rtl)
	if err != nil || len(run.Glyphs) == 0 {
		return TextToPath(mgr, text, family, weight, style, size)
	}

	path := ShapedRunToPath(face.Font, run, size)
	return path, run.Advance
}

func TextToPath(mgr *Manager, text string, family string, weight int, style string, size float64) (string, float64) {
	face := mgr.Resolve(family, weight, style)
	if face == nil {
		return "", 0
	}
	return textToPathCached(mgr, face.Name, face.Font, text, size, 0)
}

// TextToPathWithFont converts text to an SVG path using a specific font and letter spacing.
func TextToPathWithFont(f *opentype.Font, text string, size float64, letterSpacing float64) (string, float64) {
	return textToPathCached(nil, "", f, text, size, letterSpacing)
}

func textToPathCached(mgr *Manager, fontName string, f *opentype.Font, text string, size float64, letterSpacing float64) (string, float64) {
	var combined strings.Builder
	var cursor float64

	for i, r := range text {
		var gp GlyphPath
		var err error
		if mgr != nil {
			gp, err = mgr.CachedGlyphPath(fontName, r, size, f)
		} else {
			gp, err = GlyphToPath(f, r, size)
		}
		if err != nil {
			continue
		}

		if gp.D != "" {
			if cursor != 0 {
				fmt.Fprintf(&combined, "M%.4g 0", cursor)
			}
			combined.WriteString(translatePath(gp.D, cursor, 0))
		}

		cursor += gp.Advance
		if i < len([]rune(text))-1 {
			cursor += letterSpacing
		}
	}

	return combined.String(), cursor
}

func translatePath(d string, dx, dy float64) string {
	if dx == 0 && dy == 0 {
		return d
	}

	var result strings.Builder
	i := 0
	runes := []rune(d)

	for i < len(runes) {
		ch := runes[i]
		switch ch {
		case 'M', 'L':
			result.WriteRune(ch)
			i++
			x, y, next := parseTwoFloats(runes, i)
			fmt.Fprintf(&result, "%.4g %.4g", x+dx, y+dy)
			i = next
		case 'Q':
			result.WriteRune(ch)
			i++
			cx, cy, next := parseTwoFloats(runes, i)
			x, y, next2 := parseTwoFloats(runes, next)
			fmt.Fprintf(&result, "%.4g %.4g %.4g %.4g", cx+dx, cy+dy, x+dx, y+dy)
			i = next2
		case 'C':
			result.WriteRune(ch)
			i++
			cx1, cy1, next := parseTwoFloats(runes, i)
			cx2, cy2, next2 := parseTwoFloats(runes, next)
			x, y, next3 := parseTwoFloats(runes, next2)
			fmt.Fprintf(&result, "%.4g %.4g %.4g %.4g %.4g %.4g",
				cx1+dx, cy1+dy, cx2+dx, cy2+dy, x+dx, y+dy)
			i = next3
		case 'Z':
			result.WriteRune(ch)
			i++
		default:
			result.WriteRune(ch)
			i++
		}
	}

	return result.String()
}

func parseTwoFloats(runes []rune, start int) (float64, float64, int) {
	x, next := parseFloat(runes, start)
	y, next2 := parseFloat(runes, next)
	return x, y, next2
}

func parseFloat(runes []rune, start int) (float64, int) {
	for start < len(runes) && runes[start] == ' ' {
		start++
	}

	end := start
	if end < len(runes) && (runes[end] == '-' || runes[end] == '+') {
		end++
	}
	for end < len(runes) && ((runes[end] >= '0' && runes[end] <= '9') || runes[end] == '.') {
		end++
	}
	if end < len(runes) && (runes[end] == 'e' || runes[end] == 'E') {
		end++
		if end < len(runes) && (runes[end] == '-' || runes[end] == '+') {
			end++
		}
		for end < len(runes) && runes[end] >= '0' && runes[end] <= '9' {
			end++
		}
	}

	var val float64
	fmt.Sscanf(string(runes[start:end]), "%g", &val)
	return val, end
}

func GlyphIDToPath(f *opentype.Font, gid sfnt.GlyphIndex, size float64) (GlyphPath, error) {
	var buf sfnt.Buffer
	ppem := fixed.Int26_6(size * 64)

	segments, err := f.LoadGlyph(&buf, gid, ppem, nil)
	if err != nil {
		return GlyphPath{}, err
	}

	advance, err := f.GlyphAdvance(&buf, gid, ppem, 0)
	if err != nil {
		return GlyphPath{}, err
	}

	var b strings.Builder
	for _, seg := range segments {
		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			fmt.Fprintf(&b, "M%.4g %.4g", fix(seg.Args[0].X), fix(seg.Args[0].Y))
		case sfnt.SegmentOpLineTo:
			fmt.Fprintf(&b, "L%.4g %.4g", fix(seg.Args[0].X), fix(seg.Args[0].Y))
		case sfnt.SegmentOpQuadTo:
			fmt.Fprintf(&b, "Q%.4g %.4g %.4g %.4g",
				fix(seg.Args[0].X), fix(seg.Args[0].Y),
				fix(seg.Args[1].X), fix(seg.Args[1].Y))
		case sfnt.SegmentOpCubeTo:
			fmt.Fprintf(&b, "C%.4g %.4g %.4g %.4g %.4g %.4g",
				fix(seg.Args[0].X), fix(seg.Args[0].Y),
				fix(seg.Args[1].X), fix(seg.Args[1].Y),
				fix(seg.Args[2].X), fix(seg.Args[2].Y))
		}
	}

	if b.Len() > 0 {
		b.WriteString("Z")
	}

	return GlyphPath{D: b.String(), Advance: fix(advance)}, nil
}

func ShapedRunToPath(f *opentype.Font, run *ShapedRun, size float64) string {
	var combined strings.Builder
	var cursor float64

	for _, g := range run.Glyphs {
		gid := sfnt.GlyphIndex(g.GlyphID)
		gp, err := GlyphIDToPath(f, gid, size)
		if err != nil {
			cursor += g.Advance
			continue
		}

		x := cursor + g.XOffset
		y := g.YOffset

		if gp.D != "" {
			combined.WriteString(translatePath(gp.D, x, y))
		}

		cursor += g.Advance
	}

	return combined.String()
}

func fix(v fixed.Int26_6) float64 {
	return float64(v) / 64.0
}
