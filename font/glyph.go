package font

import (
	"fmt"
	"strconv"
	"strings"

	gotextfont "github.com/go-text/typesetting/font"
	ot "github.com/go-text/typesetting/font/opentype"
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

func ShapedTextToPath(mgr *Manager, text string, family string, weight int, style string, size float64, rtl bool) (string, float64) {
	return ShapedTextToPathVariations(mgr, text, family, weight, style, size, rtl, nil)
}

func ShapedTextToPathVariations(mgr *Manager, text string, family string, weight int, style string, size float64, rtl bool, vars []Variation) (string, float64) {
	face := mgr.Resolve(family, weight, style)
	if face == nil || len(face.RawData) == 0 {
		return TextToPathVariations(mgr, text, family, weight, style, size, vars)
	}

	run, err := mgr.ShapeTextWithVariations(face, text, size, rtl, vars)
	if err != nil || len(run.Glyphs) == 0 {
		return TextToPathVariations(mgr, text, family, weight, style, size, vars)
	}

	if face.Variable && len(vars) > 0 {
		gtFace, err := face.shaperFace(vars)
		if err == nil {
			return shapedRunToPathVariations(gtFace, run, size), run.Advance
		}
	}
	path := ShapedRunToPath(face.Font, run, size)
	return path, run.Advance
}

func TextToPath(mgr *Manager, text string, family string, weight int, style string, size float64) (string, float64) {
	return TextToPathVariations(mgr, text, family, weight, style, size, nil)
}

func TextToPathVariations(mgr *Manager, text string, family string, weight int, style string, size float64, vars []Variation) (string, float64) {
	face := mgr.Resolve(family, weight, style)
	if face == nil {
		return "", 0
	}
	if face.Variable && len(vars) > 0 {
		if path, adv, ok := textToPathVariable(mgr, face, text, size, vars, 0); ok {
			return path, adv
		}
	}
	return textToPathCached(mgr, face.Name, face.Font, text, size, 0)
}

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
	result.Grow(len(d))
	i := 0
	for i < len(d) {
		ch := d[i]
		switch ch {
		case 'M', 'L':
			result.WriteByte(ch)
			i++
			x, y, next := parseTwoFloats(d, i)
			fmt.Fprintf(&result, "%.4g %.4g", x+dx, y+dy)
			i = next
		case 'Q':
			result.WriteByte(ch)
			i++
			cx, cy, next := parseTwoFloats(d, i)
			x, y, next2 := parseTwoFloats(d, next)
			fmt.Fprintf(&result, "%.4g %.4g %.4g %.4g", cx+dx, cy+dy, x+dx, y+dy)
			i = next2
		case 'C':
			result.WriteByte(ch)
			i++
			cx1, cy1, next := parseTwoFloats(d, i)
			cx2, cy2, next2 := parseTwoFloats(d, next)
			x, y, next3 := parseTwoFloats(d, next2)
			fmt.Fprintf(&result, "%.4g %.4g %.4g %.4g %.4g %.4g",
				cx1+dx, cy1+dy, cx2+dx, cy2+dy, x+dx, y+dy)
			i = next3
		default:
			result.WriteByte(ch)
			i++
		}
	}

	return result.String()
}

func parseTwoFloats(s string, start int) (float64, float64, int) {
	x, next := parseFloat(s, start)
	y, next2 := parseFloat(s, next)
	return x, y, next2
}

// parseFloat scans an ASCII SVG-path number out of s starting at start. SVG
// path syntax is ASCII-only, so byte indexing is safe and faster than []rune.
func parseFloat(s string, start int) (float64, int) {
	for start < len(s) && s[start] == ' ' {
		start++
	}

	end := start
	if end < len(s) && (s[end] == '-' || s[end] == '+') {
		end++
	}
	for end < len(s) && ((s[end] >= '0' && s[end] <= '9') || s[end] == '.') {
		end++
	}
	if end < len(s) && (s[end] == 'e' || s[end] == 'E') {
		end++
		if end < len(s) && (s[end] == '-' || s[end] == '+') {
			end++
		}
		for end < len(s) && s[end] >= '0' && s[end] <= '9' {
			end++
		}
	}

	val, _ := strconv.ParseFloat(s[start:end], 64)
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

func textToPathVariable(mgr *Manager, face *Face, text string, size float64, vars []Variation, letterSpacing float64) (string, float64, bool) {
	gtFace, err := face.shaperFace(vars)
	if err != nil {
		return "", 0, false
	}
	upem := float64(gtFace.Font.Upem())
	if upem == 0 {
		return "", 0, false
	}
	scale := size / upem

	varKey := variationKey(vars)
	var combined strings.Builder
	var cursor float64
	runes := []rune(text)
	for i, r := range runes {
		gp, ok := cachedVariationGlyph(mgr, face.Name, r, size, varKey, gtFace, scale)
		if !ok {
			return "", 0, false
		}
		if gp.D != "" {
			combined.WriteString(translatePath(gp.D, cursor, 0))
		}
		cursor += gp.Advance
		if i < len(runes)-1 {
			cursor += letterSpacing
		}
	}
	return combined.String(), cursor, true
}

func cachedVariationGlyph(mgr *Manager, fontName string, r rune, size float64, varKey string, gtFace *gotextfont.Face, scale float64) (GlyphPath, bool) {
	if mgr != nil {
		if p, ok := mgr.glyphs.Get(fontName, r, size, varKey); ok {
			return p, true
		}
	}
	gid, ok := gtFace.NominalGlyph(r)
	if !ok {
		return GlyphPath{}, false
	}
	gp, ok := glyphOutlineToPath(gtFace, gid, scale)
	if !ok {
		return GlyphPath{}, false
	}
	if mgr != nil {
		mgr.glyphs.Set(fontName, r, size, varKey, gp)
	}
	return gp, true
}

// glyphOutlineToPath emits the path in Y-down (screen) orientation. go-text
// segments are in font units with OpenType Y-up, so Y is negated here.
func glyphOutlineToPath(gtFace *gotextfont.Face, gid ot.GID, scale float64) (GlyphPath, bool) {
	data := gtFace.GlyphData(gid)
	outline, ok := data.(gotextfont.GlyphOutline)
	if !ok {
		return GlyphPath{}, false
	}
	var b strings.Builder
	for _, seg := range outline.Segments {
		switch seg.Op {
		case ot.SegmentOpMoveTo:
			fmt.Fprintf(&b, "M%.4g %.4g",
				float64(seg.Args[0].X)*scale,
				-float64(seg.Args[0].Y)*scale)
		case ot.SegmentOpLineTo:
			fmt.Fprintf(&b, "L%.4g %.4g",
				float64(seg.Args[0].X)*scale,
				-float64(seg.Args[0].Y)*scale)
		case ot.SegmentOpQuadTo:
			fmt.Fprintf(&b, "Q%.4g %.4g %.4g %.4g",
				float64(seg.Args[0].X)*scale, -float64(seg.Args[0].Y)*scale,
				float64(seg.Args[1].X)*scale, -float64(seg.Args[1].Y)*scale)
		case ot.SegmentOpCubeTo:
			fmt.Fprintf(&b, "C%.4g %.4g %.4g %.4g %.4g %.4g",
				float64(seg.Args[0].X)*scale, -float64(seg.Args[0].Y)*scale,
				float64(seg.Args[1].X)*scale, -float64(seg.Args[1].Y)*scale,
				float64(seg.Args[2].X)*scale, -float64(seg.Args[2].Y)*scale)
		}
	}
	if b.Len() > 0 {
		b.WriteString("Z")
	}
	advance := float64(gtFace.HorizontalAdvance(gid)) * scale
	return GlyphPath{D: b.String(), Advance: advance}, true
}

func shapedRunToPathVariations(gtFace *gotextfont.Face, run *ShapedRun, size float64) string {
	upem := float64(gtFace.Font.Upem())
	if upem == 0 {
		return ""
	}
	scale := size / upem
	var combined strings.Builder
	var cursor float64
	for _, g := range run.Glyphs {
		gid := ot.GID(g.GlyphID)
		gp, ok := glyphOutlineToPath(gtFace, gid, scale)
		if !ok {
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
