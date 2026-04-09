package render

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/macawls/ogre/style"
)

type BackgroundLayer struct {
	Defs string
	Fill string
}

type BackgroundResult struct {
	Defs   string
	Fill   string
	Layers []BackgroundLayer
}

func splitBackgroundLayers(s string) []string {
	var layers []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				layers = append(layers, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	layers = append(layers, strings.TrimSpace(s[start:]))
	return layers
}

func renderSingleLayer(image string, cs *style.ComputedStyle, x, y, w, h float64, idGen func(string) string) BackgroundLayer {
	if strings.HasPrefix(image, "url(") {
		tmp := *cs
		tmp.BackgroundImage = image
		r := renderURLBackground(&tmp, x, y, w, h, idGen)
		return BackgroundLayer{Defs: r.Defs, Fill: r.Fill}
	}
	g, err := style.ParseGradient(image)
	if err == nil {
		distributeStops(g.Stops)
		switch g.Type {
		case style.LinearGradient, style.RepeatingLinearGradient:
			r := renderLinearGradient(g, w, h, idGen)
			return BackgroundLayer{Defs: r.Defs, Fill: r.Fill}
		case style.RadialGradient, style.RepeatingRadialGradient:
			r := renderRadialGradient(g, idGen)
			return BackgroundLayer{Defs: r.Defs, Fill: r.Fill}
		}
	}
	c, err := style.ParseColor(image)
	if err == nil && !c.IsTransparent() {
		if c.A == 1.0 {
			return BackgroundLayer{Fill: c.Hex()}
		}
		return BackgroundLayer{Fill: c.String()}
	}
	return BackgroundLayer{}
}

// RenderBackground generates the corresponding output format.
// RenderBackground generates SVG background elements from computed styles.
func RenderBackground(cs *style.ComputedStyle, x, y, w, h float64, idGen func(string) string) BackgroundResult {
	if cs.BackgroundImage != "" {
		rawLayers := splitBackgroundLayers(cs.BackgroundImage)
		if len(rawLayers) > 1 {
			layers := make([]BackgroundLayer, 0, len(rawLayers))
			var allDefs strings.Builder
			for _, raw := range rawLayers {
				layer := renderSingleLayer(raw, cs, x, y, w, h, idGen)
				if layer.Fill != "" {
					layers = append(layers, layer)
					if layer.Defs != "" {
						allDefs.WriteString(layer.Defs)
					}
				}
			}
			if len(layers) == 0 {
				if cs.BackgroundColor.IsTransparent() {
					return BackgroundResult{Fill: "none"}
				}
				fill := cs.BackgroundColor.Hex()
				if cs.BackgroundColor.A != 1.0 {
					fill = cs.BackgroundColor.String()
				}
				return BackgroundResult{Fill: fill}
			}
			return BackgroundResult{
				Defs:   allDefs.String(),
				Fill:   layers[0].Fill,
				Layers: layers,
			}
		}

		if strings.HasPrefix(cs.BackgroundImage, "url(") {
			return renderURLBackground(cs, x, y, w, h, idGen)
		}
		g, err := style.ParseGradient(cs.BackgroundImage)
		if err == nil {
			distributeStops(g.Stops)
			switch g.Type {
			case style.LinearGradient, style.RepeatingLinearGradient:
				return renderLinearGradient(g, w, h, idGen)
			case style.RadialGradient, style.RepeatingRadialGradient:
				return renderRadialGradient(g, idGen)
			}
		}
	}

	if cs.BackgroundColor.IsTransparent() {
		return BackgroundResult{Fill: "none"}
	}

	if cs.BackgroundColor.A == 1.0 {
		return BackgroundResult{Fill: cs.BackgroundColor.Hex()}
	}
	return BackgroundResult{Fill: cs.BackgroundColor.String()}
}

func gradientEndpoints(angleRad, w, h float64) (x1, y1, x2, y2 float64) {
	a := math.Mod(angleRad, 2*math.Pi)
	if a < 0 {
		a += 2 * math.Pi
	}

	var sx, sy, ex, ey float64

	var compute func(d float64)
	compute = func(d float64) {
		if d == 0 {
			sx, sy, ex, ey = 0, h, 0, 0
			return
		}
		if d == math.Pi/2 {
			sx, sy, ex, ey = 0, 0, w, 0
			return
		}

		tanD := math.Tan(d)

		if d > 0 && d < math.Pi/2 {
			invTan := 1.0 / tanD
			I := (w/2*invTan - h/2) / (tanD + invTan)
			E := tanD*I + h
			halfLen := math.Sqrt(math.Pow(w/2-I, 2) + math.Pow(h/2-E, 2))
			sinD := math.Sin(d)
			cosD := math.Cos(d)
			sx = w/2 - halfLen*sinD
			sy = h/2 + halfLen*cosD
			ex = w/2 + halfLen*sinD
			ey = h/2 - halfLen*cosD
			return
		}

		if d > math.Pi/2 && d < math.Pi {
			invTan := 1.0 / tanD
			I := (w/2*invTan + h/2) / (tanD + invTan)
			E := tanD * I
			halfLen := math.Sqrt(math.Pow(w/2-I, 2) + math.Pow(h/2-E, 2))
			sinD := math.Sin(d)
			cosD := math.Cos(d)
			sx = w/2 - halfLen*sinD
			sy = h/2 + halfLen*cosD
			ex = w/2 + halfLen*sinD
			ey = h/2 - halfLen*cosD
			return
		}

		if d >= math.Pi {
			compute(d - math.Pi)
			sx, ex = ex, sx
			sy, ey = ey, sy
		}
	}

	compute(a)
	return sx / w * 100, sy / h * 100, ex / w * 100, ey / h * 100
}

func renderLinearGradient(g style.Gradient, w, h float64, idGen func(string) string) BackgroundResult {
	id := idGen("lg")
	rad := g.Angle * math.Pi / 180
	if w <= 0 {
		w = 1
	}
	if h <= 0 {
		h = 1
	}
	x1, y1, x2, y2 := gradientEndpoints(rad, w, h)

	defs := fmt.Sprintf(`<linearGradient id="%s" x1="%.6g%%" y1="%.6g%%" x2="%.6g%%" y2="%.6g%%">`, id, x1, y1, x2, y2)
	for _, s := range g.Stops {
		defs += fmt.Sprintf(`<stop offset="%.6g%%" stop-color="%s"/>`, s.Position*100, s.Color.String())
	}
	defs += `</linearGradient>`

	return BackgroundResult{
		Defs: defs,
		Fill: fmt.Sprintf("url(#%s)", id),
	}
}

func renderRadialGradient(g style.Gradient, idGen func(string) string) BackgroundResult {
	id := idGen("rg")

	cx := g.PositionX
	cy := g.PositionY

	defs := fmt.Sprintf(`<radialGradient id="%s" cx="%.6g%%" cy="%.6g%%" r="50%%">`, id, cx, cy)
	for _, s := range g.Stops {
		defs += fmt.Sprintf(`<stop offset="%.6g%%" stop-color="%s"/>`, s.Position*100, s.Color.String())
	}
	defs += `</radialGradient>`

	return BackgroundResult{
		Defs: defs,
		Fill: fmt.Sprintf("url(#%s)", id),
	}
}

func extractURL(s string) string {
	s = strings.TrimPrefix(s, "url(")
	s = strings.TrimSuffix(s, ")")
	s = strings.TrimSpace(s)
	if (strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'")) ||
		(strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) {
		s = s[1 : len(s)-1]
	}
	return s
}

func renderURLBackground(cs *style.ComputedStyle, x, y, w, h float64, idGen func(string) string) BackgroundResult {
	href := extractURL(cs.BackgroundImage)
	id := idGen("bg")

	imgW, imgH := resolveBackgroundSize(cs.BackgroundSize, w, h)
	par := resolvePreserveAspectRatio(cs.BackgroundSize)
	patW, patH := resolvePatternSize(cs.BackgroundRepeat, imgW, imgH, w, h)
	offX, offY := resolveBackgroundPosition(cs.BackgroundPosition, w, h, imgW, imgH)

	defs := fmt.Sprintf(
		`<pattern id="%s" patternUnits="userSpaceOnUse" x="%.6g" y="%.6g" width="%.6g" height="%.6g">`+
			`<image href="%s" width="%.6g" height="%.6g" preserveAspectRatio="%s"/>`+
			`</pattern>`,
		id, x+offX, y+offY, patW, patH, xmlEscape(href), imgW, imgH, par,
	)

	return BackgroundResult{
		Defs: defs,
		Fill: fmt.Sprintf("url(#%s)", id),
	}
}

func resolveBackgroundSize(size string, elemW, elemH float64) (float64, float64) {
	size = strings.TrimSpace(size)
	switch size {
	case "", "auto":
		return elemW, elemH
	case "cover", "contain":
		return elemW, elemH
	}

	parts := strings.Fields(size)
	if len(parts) == 1 {
		parts = append(parts, "auto")
	}

	w := parseSizeDimension(parts[0], elemW)
	h := parseSizeDimension(parts[1], elemH)
	return w, h
}

func parseSizeDimension(s string, base float64) float64 {
	if s == "auto" {
		return base
	}
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil {
			return base
		}
		return base * v / 100
	}
	if strings.HasSuffix(s, "px") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "px"), 64)
		if err != nil {
			return base
		}
		return v
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return base
	}
	return v
}

func resolvePreserveAspectRatio(size string) string {
	switch strings.TrimSpace(size) {
	case "cover":
		return "xMidYMid slice"
	case "contain":
		return "xMidYMid meet"
	default:
		return "none"
	}
}

func resolvePatternSize(repeat string, imgW, imgH, elemW, elemH float64) (float64, float64) {
	switch strings.TrimSpace(repeat) {
	case "no-repeat":
		return elemW, elemH
	case "repeat-x":
		return imgW, elemH
	case "repeat-y":
		return elemW, imgH
	default:
		return imgW, imgH
	}
}

func resolveBackgroundPosition(pos string, elemW, elemH, imgW, imgH float64) (float64, float64) {
	pos = strings.TrimSpace(pos)
	if pos == "" {
		return 0, 0
	}

	parts := strings.Fields(pos)
	if len(parts) == 1 {
		switch parts[0] {
		case "center":
			parts = []string{"50%", "50%"}
		case "left":
			parts = []string{"0%", "50%"}
		case "right":
			parts = []string{"100%", "50%"}
		case "top":
			parts = []string{"50%", "0%"}
		case "bottom":
			parts = []string{"50%", "100%"}
		default:
			parts = append(parts, "50%")
		}
	}

	xStr := resolvePositionKeyword(parts[0], true)
	yStr := resolvePositionKeyword(parts[1], false)

	offX := parsePositionValue(xStr, elemW, imgW)
	offY := parsePositionValue(yStr, elemH, imgH)
	return offX, offY
}

func resolvePositionKeyword(s string, isX bool) string {
	switch s {
	case "left":
		return "0%"
	case "right":
		return "100%"
	case "top":
		return "0%"
	case "bottom":
		return "100%"
	case "center":
		return "50%"
	default:
		return s
	}
}

func parsePositionValue(s string, elemDim, imgDim float64) float64 {
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
		if err != nil {
			return 0
		}
		return (elemDim - imgDim) * v / 100
	}
	if strings.HasSuffix(s, "px") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "px"), 64)
		if err != nil {
			return 0
		}
		return v
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func distributeStops(stops []style.ColorStop) {
	n := len(stops)
	if n == 0 {
		return
	}

	if !stops[0].HasPos {
		stops[0].Position = 0
		stops[0].HasPos = true
	}
	if !stops[n-1].HasPos {
		stops[n-1].Position = 1
		stops[n-1].HasPos = true
	}

	i := 1
	for i < n-1 {
		if stops[i].HasPos {
			i++
			continue
		}
		start := i - 1
		end := i + 1
		for end < n && !stops[end].HasPos {
			end++
		}
		count := end - start
		for j := start + 1; j < end; j++ {
			stops[j].Position = stops[start].Position + (stops[end].Position-stops[start].Position)*float64(j-start)/float64(count)
			stops[j].HasPos = true
		}
		i = end
	}
}
