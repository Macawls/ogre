package render

import (
	"fmt"
	"strings"

	"github.com/macawls/ogre/style"
)

// RenderBoxShadow generates the corresponding output format.
// RenderBoxShadow generates SVG filter elements for box shadows.
func RenderBoxShadow(shadows []style.Shadow, x, y, w, h float64, borderRadius float64, idGen func(string) string) string {
	if len(shadows) == 0 {
		return ""
	}

	var defs strings.Builder
	var shapes strings.Builder

	defs.WriteString("<defs>")

	for _, s := range shadows {
		if s.Inset {
			continue
		}

		filterID := idGen("shadow")
		stdDev := s.Blur / 2

		defs.WriteString(fmt.Sprintf(
			`<filter id="%s" x="-50%%" y="-50%%" width="200%%" height="200%%">`+
				`<feGaussianBlur in="SourceAlpha" stdDeviation="%.4g"/>`+
				`<feOffset dx="%.4g" dy="%.4g" result="offsetBlur"/>`+
				`<feFlood flood-color="%s" flood-opacity="%.4g"/>`+
				`<feComposite in2="offsetBlur" operator="in"/>`+
				`<feMerge><feMergeNode/><feMergeNode in="SourceGraphic"/></feMerge>`+
				`</filter>`,
			filterID, stdDev, s.OffsetX, s.OffsetY,
			shadowColorHex(s.Color), s.Color.A))

		sx := x - s.Spread
		sy := y - s.Spread
		sw := w + s.Spread*2
		sh := h + s.Spread*2

		fmt.Fprintf(&shapes, `<rect x="%.4g" y="%.4g" width="%.4g" height="%.4g"`, sx, sy, sw, sh)
		fmt.Fprintf(&shapes, ` fill="%s"`, shadowColorCSS(s.Color))
		if borderRadius > 0 {
			fmt.Fprintf(&shapes, ` rx="%.4g"`, borderRadius)
		}
		fmt.Fprintf(&shapes, ` filter="url(#%s)"`, filterID)
		shapes.WriteString("/>")
	}

	defs.WriteString("</defs>")

	return defs.String() + shapes.String()
}

// RenderInsetBoxShadow generates the corresponding output format.
// RenderInsetBoxShadow generates SVG elements for inset shadows.
func RenderInsetBoxShadow(shadows []style.Shadow, x, y, w, h float64, borderRadius float64, idGen func(string) string) string {
	if len(shadows) == 0 {
		return ""
	}

	var defs strings.Builder
	var shapes strings.Builder

	defs.WriteString("<defs>")

	for _, s := range shadows {
		if !s.Inset {
			continue
		}

		filterID := idGen("inset-shadow")
		clipID := idGen("inset-clip")
		stdDev := s.Blur / 2

		defs.WriteString(fmt.Sprintf(
			`<clipPath id="%s"><rect x="%.4g" y="%.4g" width="%.4g" height="%.4g"`,
			clipID, x, y, w, h))
		if borderRadius > 0 {
			fmt.Fprintf(&defs, ` rx="%.4g"`, borderRadius)
		}
		defs.WriteString("/></clipPath>")

		defs.WriteString(fmt.Sprintf(
			`<filter id="%s" x="-50%%" y="-50%%" width="200%%" height="200%%">`+
				`<feGaussianBlur in="SourceGraphic" stdDeviation="%.4g"/>`+
				`</filter>`,
			filterID, stdDev))

		spread := s.Spread
		border := s.Blur + spread
		ox := s.OffsetX
		oy := s.OffsetY

		fmt.Fprintf(&shapes, `<g clip-path="url(#%s)">`, clipID)

		sides := []struct {
			rx, ry, rw, rh float64
		}{
			{x + ox, y + oy, w, border},
			{x + ox, y + h - border + oy, w, border},
			{x + ox, y + oy, border, h},
			{x + w - border + ox, y + oy, border, h},
		}

		for _, side := range sides {
			fmt.Fprintf(&shapes, `<rect x="%.4g" y="%.4g" width="%.4g" height="%.4g"`,
				side.rx, side.ry, side.rw, side.rh)
			fmt.Fprintf(&shapes, ` fill="%s"`, shadowColorCSS(s.Color))
			fmt.Fprintf(&shapes, ` filter="url(#%s)"`, filterID)
			shapes.WriteString("/>")
		}

		shapes.WriteString("</g>")
	}

	defs.WriteString("</defs>")

	return defs.String() + shapes.String()
}

func shadowColorHex(c style.Color) string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

func shadowColorCSS(c style.Color) string {
	if c.A == 1.0 {
		return shadowColorHex(c)
	}
	return fmt.Sprintf("rgba(%d,%d,%d,%.4g)", c.R, c.G, c.B, c.A)
}
