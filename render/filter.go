package render

import (
	"fmt"
	"strconv"
	"strings"
)

// RenderCSSFilter generates the corresponding output format.
// RenderCSSFilter converts CSS filter functions to SVG filter elements.
func RenderCSSFilter(filter string, idGen func(string) string) (defsContent string, filterAttr string) {
	filter = strings.TrimSpace(filter)
	if filter == "" || strings.ToLower(filter) == "none" {
		return "", ""
	}

	funcs := parseFilterFuncs(filter)
	if len(funcs) == 0 {
		return "", ""
	}

	var primitives []string
	for _, f := range funcs {
		p := filterToPrimitive(f.name, f.value)
		if p != "" {
			primitives = append(primitives, p)
		}
	}

	if len(primitives) == 0 {
		return "", ""
	}

	id := idGen("filter")
	defs := fmt.Sprintf(`<filter id="%s">%s</filter>`, id, strings.Join(primitives, ""))
	return defs, fmt.Sprintf(`filter="url(#%s)"`, id)
}

type filterFunc struct {
	name  string
	value string
}

func parseFilterFuncs(s string) []filterFunc {
	var result []filterFunc
	for s != "" {
		s = strings.TrimSpace(s)
		if s == "" {
			break
		}
		idx := strings.Index(s, "(")
		if idx < 0 {
			break
		}
		name := strings.TrimSpace(s[:idx])

		depth := 0
		end := -1
		for i := idx; i < len(s); i++ {
			if s[i] == '(' {
				depth++
			} else if s[i] == ')' {
				depth--
				if depth == 0 {
					end = i
					break
				}
			}
		}
		if end < 0 {
			break
		}

		value := strings.TrimSpace(s[idx+1 : end])
		result = append(result, filterFunc{name: name, value: value})
		s = s[end+1:]
	}
	return result
}

func filterToPrimitive(name, value string) string {
	switch name {
	case "blur":
		px := parsePxValue(value)
		if px <= 0 {
			return ""
		}
		return fmt.Sprintf(`<feGaussianBlur stdDeviation="%.4g"/>`, px)

	case "grayscale":
		pct := parsePercentOrFraction(value)
		sat := 1 - pct
		if sat < 0 {
			sat = 0
		}
		return fmt.Sprintf(`<feColorMatrix type="saturate" values="%.4g"/>`, sat)

	case "brightness":
		slope := parsePercentOrFraction(value)
		return fmt.Sprintf(
			`<feComponentTransfer>`+
				`<feFuncR type="linear" slope="%.4g"/>`+
				`<feFuncG type="linear" slope="%.4g"/>`+
				`<feFuncB type="linear" slope="%.4g"/>`+
				`</feComponentTransfer>`, slope, slope, slope)

	case "contrast":
		x := parsePercentOrFraction(value)
		intercept := -(0.5*x) + 0.5
		return fmt.Sprintf(
			`<feComponentTransfer>`+
				`<feFuncR type="linear" slope="%.4g" intercept="%.4g"/>`+
				`<feFuncG type="linear" slope="%.4g" intercept="%.4g"/>`+
				`<feFuncB type="linear" slope="%.4g" intercept="%.4g"/>`+
				`</feComponentTransfer>`, x, intercept, x, intercept, x, intercept)

	case "saturate":
		x := parsePercentOrFraction(value)
		return fmt.Sprintf(`<feColorMatrix type="saturate" values="%.4g"/>`, x)

	case "sepia":
		pct := parsePercentOrFraction(value)
		r0 := 1 - 0.607*pct
		r1 := 0.769 * pct
		r2 := 0.189 * pct
		g0 := 0.349 * pct
		g1 := 1 - 0.314*pct
		g2 := 0.168 * pct
		b0 := 0.272 * pct
		b1 := 0.534 * pct
		b2 := 1 - 0.869*pct
		return fmt.Sprintf(
			`<feColorMatrix type="matrix" values="%.4g %.4g %.4g 0 0 %.4g %.4g %.4g 0 0 %.4g %.4g %.4g 0 0 0 0 0 1 0"/>`,
			r0, r1, r2, g0, g1, g2, b0, b1, b2)

	case "hue-rotate":
		deg := parseDegValue(value)
		return fmt.Sprintf(`<feColorMatrix type="hueRotate" values="%.4g"/>`, deg)

	case "invert":
		pct := parsePercentOrFraction(value)
		lo := pct
		hi := 1 - pct
		return fmt.Sprintf(
			`<feComponentTransfer>`+
				`<feFuncR type="table" tableValues="%.4g %.4g"/>`+
				`<feFuncG type="table" tableValues="%.4g %.4g"/>`+
				`<feFuncB type="table" tableValues="%.4g %.4g"/>`+
				`</feComponentTransfer>`, lo, hi, lo, hi, lo, hi)

	case "drop-shadow":
		return parseDropShadow(value)

	case "opacity":
		return ""

	default:
		return ""
	}
}

func parseDropShadow(value string) string {
	parts := strings.Fields(value)
	if len(parts) < 2 {
		return ""
	}

	dx := parsePxValue(parts[0])
	dy := parsePxValue(parts[1])

	var blur float64
	var floodColor string

	if len(parts) >= 3 {
		if isNumericValue(parts[2]) {
			blur = parsePxValue(parts[2])
			if len(parts) >= 4 {
				floodColor = strings.Join(parts[3:], " ")
			}
		} else {
			floodColor = strings.Join(parts[2:], " ")
		}
	}

	stdDev := blur / 2
	if floodColor == "" {
		floodColor = "black"
	}

	return fmt.Sprintf(`<feDropShadow dx="%.4g" dy="%.4g" stdDeviation="%.4g" flood-color="%s"/>`,
		dx, dy, stdDev, floodColor)
}

func isNumericValue(s string) bool {
	s = strings.TrimSuffix(s, "px")
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func parsePxValue(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "px")
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parseDegValue(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "deg")
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parsePercentOrFraction(s string) float64 {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		v, _ := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		return v / 100
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}
