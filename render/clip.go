package render

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/macawls/ogre/style"
)

// RenderOverflowClip generates the corresponding output format.
// RenderOverflowClip generates SVG clip-path for overflow hidden.
func RenderOverflowClip(cs *style.ComputedStyle, x, y, w, h float64, idGen func(string) string) (defsContent string, clipAttr string) {
	if cs.ClipPath != "" {
		return renderCSSClipPath(cs.ClipPath, x, y, w, h, idGen)
	}

	if cs.Overflow != style.OverflowHidden {
		return "", ""
	}

	id := idGen("clip")
	var defs strings.Builder
	fmt.Fprintf(&defs, `<clipPath id="%s">`, id)

	tl := cs.BorderTopLeftRadius
	tr := cs.BorderTopRightRadius
	bl := cs.BorderBottomLeftRadius
	br := cs.BorderBottomRightRadius

	if tl == 0 && tr == 0 && bl == 0 && br == 0 {
		fmt.Fprintf(&defs, `<rect x="%.4g" y="%.4g" width="%.4g" height="%.4g"/>`, x, y, w, h)
	} else if tl == tr && tr == bl && bl == br {
		fmt.Fprintf(&defs, `<rect x="%.4g" y="%.4g" width="%.4g" height="%.4g" rx="%.4g"/>`, x, y, w, h, tl)
	} else {
		defs.WriteString(roundedRectPath(x, y, w, h, tl, tr, br, bl))
	}

	defs.WriteString("</clipPath>")
	return defs.String(), fmt.Sprintf(`clip-path="url(#%s)"`, id)
}

func roundedRectPath(x, y, w, h, tl, tr, br, bl float64) string {
	var b strings.Builder
	fmt.Fprintf(&b, `<path d="M%.4g,%.4g`, x+tl, y)
	fmt.Fprintf(&b, ` H%.4g`, x+w-tr)
	if tr > 0 {
		fmt.Fprintf(&b, ` A%.4g,%.4g 0 0 1 %.4g,%.4g`, tr, tr, x+w, y+tr)
	}
	fmt.Fprintf(&b, ` V%.4g`, y+h-br)
	if br > 0 {
		fmt.Fprintf(&b, ` A%.4g,%.4g 0 0 1 %.4g,%.4g`, br, br, x+w-br, y+h)
	}
	fmt.Fprintf(&b, ` H%.4g`, x+bl)
	if bl > 0 {
		fmt.Fprintf(&b, ` A%.4g,%.4g 0 0 1 %.4g,%.4g`, bl, bl, x, y+h-bl)
	}
	fmt.Fprintf(&b, ` V%.4g`, y+tl)
	if tl > 0 {
		fmt.Fprintf(&b, ` A%.4g,%.4g 0 0 1 %.4g,%.4g`, tl, tl, x+tl, y)
	}
	b.WriteString(` Z"/>`)
	return b.String()
}

func renderCSSClipPath(raw string, x, y, w, h float64, idGen func(string) string) (string, string) {
	id := idGen("clip")
	var defs strings.Builder
	fmt.Fprintf(&defs, `<clipPath id="%s">`, id)

	val := strings.TrimSpace(raw)

	switch {
	case strings.HasPrefix(val, "circle("):
		inner := extractParens(val, "circle(")
		defs.WriteString(parseCircle(inner, x, y, w, h))

	case strings.HasPrefix(val, "ellipse("):
		inner := extractParens(val, "ellipse(")
		defs.WriteString(parseEllipse(inner, x, y, w, h))

	case strings.HasPrefix(val, "polygon("):
		inner := extractParens(val, "polygon(")
		defs.WriteString(parsePolygon(inner, x, y, w, h))

	case strings.HasPrefix(val, "inset("):
		inner := extractParens(val, "inset(")
		defs.WriteString(parseInset(inner, x, y, w, h))

	default:
		return "", ""
	}

	defs.WriteString("</clipPath>")
	return defs.String(), fmt.Sprintf(`clip-path="url(#%s)"`, id)
}

func extractParens(val, prefix string) string {
	s := strings.TrimPrefix(val, prefix)
	s = strings.TrimSuffix(s, ")")
	return strings.TrimSpace(s)
}

func parseCircle(inner string, x, y, w, h float64) string {
	parts := strings.SplitN(inner, " at ", 2)
	r := resolveValue(strings.TrimSpace(parts[0]), w)
	cx := x + w/2
	cy := y + h/2
	if len(parts) == 2 {
		pos := strings.Fields(strings.TrimSpace(parts[1]))
		if len(pos) >= 1 {
			cx = x + resolveValue(pos[0], w)
		}
		if len(pos) >= 2 {
			cy = y + resolveValue(pos[1], h)
		}
	}
	return fmt.Sprintf(`<circle cx="%.4g" cy="%.4g" r="%.4g"/>`, cx, cy, r)
}

func parseEllipse(inner string, x, y, w, h float64) string {
	parts := strings.SplitN(inner, " at ", 2)
	radii := strings.Fields(strings.TrimSpace(parts[0]))
	rx := resolveValue(radii[0], w)
	ry := rx
	if len(radii) >= 2 {
		ry = resolveValue(radii[1], h)
	}
	cx := x + w/2
	cy := y + h/2
	if len(parts) == 2 {
		pos := strings.Fields(strings.TrimSpace(parts[1]))
		if len(pos) >= 1 {
			cx = x + resolveValue(pos[0], w)
		}
		if len(pos) >= 2 {
			cy = y + resolveValue(pos[1], h)
		}
	}
	return fmt.Sprintf(`<ellipse cx="%.4g" cy="%.4g" rx="%.4g" ry="%.4g"/>`, cx, cy, rx, ry)
}

func parsePolygon(inner string, x, y, w, h float64) string {
	points := strings.Split(inner, ",")
	var parts []string
	for _, p := range points {
		coords := strings.Fields(strings.TrimSpace(p))
		if len(coords) < 2 {
			continue
		}
		px := x + resolveValue(coords[0], w)
		py := y + resolveValue(coords[1], h)
		parts = append(parts, fmt.Sprintf("%.4g,%.4g", px, py))
	}
	return fmt.Sprintf(`<polygon points="%s"/>`, strings.Join(parts, " "))
}

func parseInset(inner string, x, y, w, h float64) string {
	fields := strings.Fields(inner)
	var top, right, bottom, left float64
	switch len(fields) {
	case 1:
		top = resolveValuePx(fields[0])
		right, bottom, left = top, top, top
	case 2:
		top = resolveValuePx(fields[0])
		right = resolveValuePx(fields[1])
		bottom, left = top, right
	case 3:
		top = resolveValuePx(fields[0])
		right = resolveValuePx(fields[1])
		bottom = resolveValuePx(fields[2])
		left = right
	case 4:
		top = resolveValuePx(fields[0])
		right = resolveValuePx(fields[1])
		bottom = resolveValuePx(fields[2])
		left = resolveValuePx(fields[3])
	}
	rx := x + left
	ry := y + top
	rw := w - left - right
	rh := h - top - bottom
	return fmt.Sprintf(`<rect x="%.4g" y="%.4g" width="%.4g" height="%.4g"/>`, rx, ry, rw, rh)
}

func resolveValue(s string, ref float64) float64 {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		v, _ := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		return v / 100 * ref
	}
	return resolveValuePx(s)
}

func resolveValuePx(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "px")
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
