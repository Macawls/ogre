package render

import (
	"fmt"
	"math"
	"strings"

	"github.com/macawls/ogre/style"
)

type borderSide struct {
	width float64
	style style.BorderStyle
	color style.Color
}

// RenderBorders generates the corresponding output format.
// RenderBorders generates SVG border elements.
func RenderBorders(cs *style.ComputedStyle, x, y, w, h float64) string {
	sides := [4]borderSide{
		{cs.BorderTopWidth, cs.BorderTopStyle, cs.BorderTopColor},
		{cs.BorderRightWidth, cs.BorderRightStyle, cs.BorderRightColor},
		{cs.BorderBottomWidth, cs.BorderBottomStyle, cs.BorderBottomColor},
		{cs.BorderLeftWidth, cs.BorderLeftStyle, cs.BorderLeftColor},
	}

	allNone := true
	for _, s := range sides {
		if s.width > 0 && s.style != style.BorderStyleNone {
			allNone = false
			break
		}
	}
	if allNone {
		return ""
	}

	if isUniform(sides) {
		return renderUniformBorder(sides[0], cs, x, y, w, h)
	}
	return renderMixedBorders(sides, cs, x, y, w, h)
}

func isUniform(sides [4]borderSide) bool {
	s := sides[0]
	for i := 1; i < 4; i++ {
		if sides[i].width != s.width || sides[i].style != s.style || sides[i].color != s.color {
			return false
		}
	}
	return true
}

func renderUniformBorder(s borderSide, cs *style.ComputedStyle, x, y, w, h float64) string {
	if s.style == style.BorderStyleNone || s.width == 0 {
		return ""
	}

	half := s.width / 2
	rx := math.Min(cs.BorderTopLeftRadius, w/2)
	ry := math.Min(cs.BorderTopLeftRadius, h/2)

	if s.style == style.BorderStyleDouble {
		return renderDoubleRect(s, rx, ry, x, y, w, h)
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<rect x="%.4g" y="%.4g" width="%.4g" height="%.4g" fill="none" stroke="%s" stroke-width="%.4g"`,
		x+half, y+half, w-s.width, h-s.width, colorToCSS(s.color), s.width)

	if rx > 0 {
		fmt.Fprintf(&b, ` rx="%.4g"`, rx)
	}
	if ry > 0 {
		fmt.Fprintf(&b, ` ry="%.4g"`, ry)
	}

	da := dashArray(s.style, s.width)
	if da != "" {
		fmt.Fprintf(&b, ` stroke-dasharray="%s"`, da)
	}

	b.WriteString("/>")
	return b.String()
}

func renderDoubleRect(s borderSide, rx, ry, x, y, w, h float64) string {
	third := s.width / 3
	var b strings.Builder

	outerHalf := third / 2
	fmt.Fprintf(&b, `<rect x="%.4g" y="%.4g" width="%.4g" height="%.4g" fill="none" stroke="%s" stroke-width="%.4g"`,
		x+outerHalf, y+outerHalf, w-third, h-third, colorToCSS(s.color), third)
	if rx > 0 {
		fmt.Fprintf(&b, ` rx="%.4g"`, rx)
	}
	if ry > 0 {
		fmt.Fprintf(&b, ` ry="%.4g"`, ry)
	}
	b.WriteString("/>")

	innerHalf := third / 2
	inset := s.width - innerHalf
	fmt.Fprintf(&b, `<rect x="%.4g" y="%.4g" width="%.4g" height="%.4g" fill="none" stroke="%s" stroke-width="%.4g"`,
		x+inset, y+inset, w-2*s.width+third, h-2*s.width+third, colorToCSS(s.color), third)
	innerRx := math.Max(0, rx-s.width+third)
	innerRy := math.Max(0, ry-s.width+third)
	if innerRx > 0 {
		fmt.Fprintf(&b, ` rx="%.4g"`, innerRx)
	}
	if innerRy > 0 {
		fmt.Fprintf(&b, ` ry="%.4g"`, innerRy)
	}
	b.WriteString("/>")

	return b.String()
}

func renderMixedBorders(sides [4]borderSide, cs *style.ComputedStyle, x, y, w, h float64) string {
	radii := [4]float64{
		cs.BorderTopLeftRadius,
		cs.BorderTopRightRadius,
		cs.BorderBottomRightRadius,
		cs.BorderBottomLeftRadius,
	}
	hasRadius := radii[0] > 0 || radii[1] > 0 || radii[2] > 0 || radii[3] > 0

	var b strings.Builder

	if hasRadius {
		renderMixedWithRadius(&b, sides, radii, x, y, w, h)
	} else {
		renderMixedLines(&b, sides, x, y, w, h)
	}

	return b.String()
}

func renderMixedLines(b *strings.Builder, sides [4]borderSide, x, y, w, h float64) {
	type lineSpec struct {
		x1, y1, x2, y2 float64
		side            borderSide
	}

	lines := []lineSpec{
		{x, y, x + w, y, sides[0]},
		{x + w, y, x + w, y + h, sides[1]},
		{x + w, y + h, x, y + h, sides[2]},
		{x, y + h, x, y, sides[3]},
	}

	for _, l := range lines {
		if l.side.width == 0 || l.side.style == style.BorderStyleNone {
			continue
		}

		if l.side.style == style.BorderStyleDouble {
			renderDoubleLine(b, l.x1, l.y1, l.x2, l.y2, l.side)
			continue
		}

		fmt.Fprintf(b, `<line x1="%.4g" y1="%.4g" x2="%.4g" y2="%.4g" stroke="%s" stroke-width="%.4g"`,
			l.x1, l.y1, l.x2, l.y2, colorToCSS(l.side.color), l.side.width)

		da := dashArray(l.side.style, l.side.width)
		if da != "" {
			fmt.Fprintf(b, ` stroke-dasharray="%s"`, da)
		}

		b.WriteString("/>")
	}
}

func renderDoubleLine(b *strings.Builder, x1, y1, x2, y2 float64, s borderSide) {
	third := s.width / 3
	dx := x2 - x1
	dy := y2 - y1
	length := math.Sqrt(dx*dx + dy*dy)
	if length == 0 {
		return
	}

	nx := -dy / length
	ny := dx / length

	offset := s.width/2 - third/2

	fmt.Fprintf(b, `<line x1="%.4g" y1="%.4g" x2="%.4g" y2="%.4g" stroke="%s" stroke-width="%.4g"/>`,
		x1+nx*offset, y1+ny*offset, x2+nx*offset, y2+ny*offset, colorToCSS(s.color), third)

	fmt.Fprintf(b, `<line x1="%.4g" y1="%.4g" x2="%.4g" y2="%.4g" stroke="%s" stroke-width="%.4g"/>`,
		x1-nx*offset, y1-ny*offset, x2-nx*offset, y2-ny*offset, colorToCSS(s.color), third)
}

func renderMixedWithRadius(b *strings.Builder, sides [4]borderSide, radii [4]float64, x, y, w, h float64) {
	type corner struct {
		cx, cy float64
		r      float64
	}

	corners := [4]corner{
		{x + radii[0], y + radii[0], radii[0]},
		{x + w - radii[1], y + radii[1], radii[1]},
		{x + w - radii[2], y + h - radii[2], radii[2]},
		{x + radii[3], y + h - radii[3], radii[3]},
	}

	type segment struct {
		side borderSide
		d    string
	}

	var segments []segment

	topStart := fmt.Sprintf("%.4g %.4g", x+radii[0], y)
	topEnd := fmt.Sprintf("%.4g %.4g", x+w-radii[1], y)
	var topArcStart, topArcEnd string
	if radii[0] > 0 {
		topArcStart = fmt.Sprintf("M %.4g %.4g A %.4g %.4g 0 0 1 %s",
			corners[0].cx-corners[0].r, corners[0].cy, radii[0], radii[0], topStart)
	}
	_ = topArcStart
	if radii[1] > 0 {
		topArcEnd = fmt.Sprintf(" A %.4g %.4g 0 0 1 %.4g %.4g",
			radii[1], radii[1], corners[1].cx+corners[1].r, corners[1].cy)
	}
	_ = topArcEnd

	if sides[0].width > 0 && sides[0].style != style.BorderStyleNone {
		d := fmt.Sprintf("M %s L %s", topStart, topEnd)
		if radii[1] > 0 {
			d += fmt.Sprintf(" A %.4g %.4g 0 0 1 %.4g %.4g", radii[1], radii[1], x+w, y+radii[1])
		}
		segments = append(segments, segment{sides[0], d})
	}

	if sides[1].width > 0 && sides[1].style != style.BorderStyleNone {
		d := fmt.Sprintf("M %.4g %.4g L %.4g %.4g", x+w, y+radii[1], x+w, y+h-radii[2])
		if radii[2] > 0 {
			d += fmt.Sprintf(" A %.4g %.4g 0 0 1 %.4g %.4g", radii[2], radii[2], x+w-radii[2], y+h)
		}
		segments = append(segments, segment{sides[1], d})
	}

	if sides[2].width > 0 && sides[2].style != style.BorderStyleNone {
		d := fmt.Sprintf("M %.4g %.4g L %.4g %.4g", x+w-radii[2], y+h, x+radii[3], y+h)
		if radii[3] > 0 {
			d += fmt.Sprintf(" A %.4g %.4g 0 0 1 %.4g %.4g", radii[3], radii[3], x, y+h-radii[3])
		}
		segments = append(segments, segment{sides[2], d})
	}

	if sides[3].width > 0 && sides[3].style != style.BorderStyleNone {
		d := fmt.Sprintf("M %.4g %.4g L %.4g %.4g", x, y+h-radii[3], x, y+radii[0])
		if radii[0] > 0 {
			d += fmt.Sprintf(" A %.4g %.4g 0 0 1 %.4g %.4g", radii[0], radii[0], x+radii[0], y)
		}
		segments = append(segments, segment{sides[3], d})
	}

	for _, seg := range segments {
		fmt.Fprintf(b, `<path d="%s" fill="none" stroke="%s" stroke-width="%.4g"`,
			seg.d, colorToCSS(seg.side.color), seg.side.width)
		da := dashArray(seg.side.style, seg.side.width)
		if da != "" {
			fmt.Fprintf(b, ` stroke-dasharray="%s"`, da)
		}
		b.WriteString("/>")
	}
}

func dashArray(s style.BorderStyle, width float64) string {
	switch s {
	case style.BorderStyleDashed:
		return fmt.Sprintf("%.4g %.4g", width*2, width)
	case style.BorderStyleDotted:
		return fmt.Sprintf("%.4g %.4g", width, width)
	default:
		return ""
	}
}
