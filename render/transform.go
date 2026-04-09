package render

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// RenderTransform generates the corresponding output format.
// RenderTransform converts CSS transform to SVG transform attribute.
func RenderTransform(transform, transformOrigin string, x, y, w, h float64) string {
	transform = strings.TrimSpace(transform)
	if transform == "" || strings.ToLower(transform) == "none" {
		return ""
	}

	originX, originY := parseTransformOrigin(transformOrigin, x, y, w, h)
	inner := parseTransformFunctions(transform)
	if inner == "" {
		return ""
	}

	return fmt.Sprintf("translate(%.4g, %.4g) %s translate(%.4g, %.4g)",
		originX, originY, inner, -originX, -originY)
}

func parseTransformOrigin(origin string, x, y, w, h float64) (float64, float64) {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return x + w/2, y + h/2
	}

	parts := strings.Fields(origin)
	ox := resolveOriginValue(parts[0], x, w)
	oy := y + h/2
	if len(parts) > 1 {
		oy = resolveOriginValue(parts[1], y, h)
	}
	return ox, oy
}

func resolveOriginValue(s string, offset, size float64) float64 {
	s = strings.TrimSpace(s)
	switch s {
	case "center":
		return offset + size/2
	case "left", "top":
		return offset
	case "right", "bottom":
		return offset + size
	}
	if strings.HasSuffix(s, "%") {
		var pct float64
		fmt.Sscanf(s, "%f%%", &pct)
		return offset + size*pct/100
	}
	if strings.HasSuffix(s, "px") {
		var v float64
		fmt.Sscanf(s, "%fpx", &v)
		return offset + v
	}
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		return offset + v
	}
	return offset + size/2
}

func parseTransformFunctions(s string) string {
	var parts []string
	i := 0
	for i < len(s) {
		for i < len(s) && s[i] == ' ' {
			i++
		}
		if i >= len(s) {
			break
		}

		start := i
		for i < len(s) && s[i] != '(' {
			i++
		}
		if i >= len(s) {
			break
		}
		fname := strings.TrimSpace(s[start:i])
		i++

		argStart := i
		for i < len(s) && s[i] != ')' {
			i++
		}
		if i >= len(s) {
			break
		}
		args := strings.TrimSpace(s[argStart:i])
		i++

		if svgPart := convertTransformFunc(fname, args); svgPart != "" {
			parts = append(parts, svgPart)
		}
	}
	return strings.Join(parts, " ")
}

func convertTransformFunc(name, args string) string {
	switch name {
	case "translate":
		vals := splitArgs(args)
		tx := parseLengthArg(vals[0])
		ty := 0.0
		if len(vals) > 1 {
			ty = parseLengthArg(vals[1])
		}
		return fmt.Sprintf("translate(%.4g, %.4g)", tx, ty)
	case "translateX":
		return fmt.Sprintf("translate(%.4g, 0)", parseLengthArg(args))
	case "translateY":
		return fmt.Sprintf("translate(0, %.4g)", parseLengthArg(args))
	case "scale":
		vals := splitArgs(args)
		sx := parseNumericArg(vals[0])
		sy := sx
		if len(vals) > 1 {
			sy = parseNumericArg(vals[1])
		}
		if sx == sy {
			return fmt.Sprintf("scale(%.4g)", sx)
		}
		return fmt.Sprintf("scale(%.4g, %.4g)", sx, sy)
	case "scaleX":
		return fmt.Sprintf("scale(%.4g, 1)", parseNumericArg(args))
	case "scaleY":
		return fmt.Sprintf("scale(1, %.4g)", parseNumericArg(args))
	case "rotate":
		return fmt.Sprintf("rotate(%.4g)", parseAngleArg(args))
	case "skewX":
		return fmt.Sprintf("skewX(%.4g)", parseAngleArg(args))
	case "skewY":
		return fmt.Sprintf("skewY(%.4g)", parseAngleArg(args))
	case "matrix":
		vals := splitArgs(args)
		if len(vals) == 6 {
			a := parseNumericArg(vals[0])
			b := parseNumericArg(vals[1])
			c := parseNumericArg(vals[2])
			d := parseNumericArg(vals[3])
			e := parseNumericArg(vals[4])
			f := parseNumericArg(vals[5])
			return fmt.Sprintf("matrix(%.4g, %.4g, %.4g, %.4g, %.4g, %.4g)", a, b, c, d, e, f)
		}
	}
	return ""
}

func splitArgs(s string) []string {
	s = strings.TrimSpace(s)
	var parts []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			parts = append(parts, p)
		}
	}
	if len(parts) <= 1 {
		parts = nil
		for _, p := range strings.Fields(s) {
			parts = append(parts, p)
		}
	}
	return parts
}

func parseLengthArg(s string) float64 {
	s = strings.TrimSpace(s)
	for _, suffix := range []string{"px", "em", "rem"} {
		if strings.HasSuffix(s, suffix) {
			s = s[:len(s)-len(suffix)]
			break
		}
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parseNumericArg(s string) float64 {
	s = strings.TrimSpace(s)
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parseAngleArg(s string) float64 {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "deg") {
		v, _ := strconv.ParseFloat(s[:len(s)-3], 64)
		return v
	}
	if strings.HasSuffix(s, "rad") {
		v, _ := strconv.ParseFloat(s[:len(s)-3], 64)
		return v * 180 / math.Pi
	}
	if strings.HasSuffix(s, "turn") {
		v, _ := strconv.ParseFloat(s[:len(s)-4], 64)
		return v * 360
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
