package style

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// GradientType identifies whether a gradient is linear or radial, and whether it repeats.
type GradientType int

const (
	LinearGradient GradientType = iota
	RadialGradient
	RepeatingLinearGradient
	RepeatingRadialGradient
)

// ColorStop is a color at an optional position within a gradient.
type ColorStop struct {
	Color    Color
	Position float64
	HasPos   bool
}

// Gradient represents a parsed CSS gradient with type, angle, and color stops.
type Gradient struct {
	Type      GradientType
	Angle     float64
	Stops     []ColorStop
	Shape     string
	Size      string
	PositionX float64
	PositionY float64
	Repeating bool
}

var directionKeywords = map[string]float64{
	"to top":          0,
	"to top right":    45,
	"to right top":    45,
	"to right":        90,
	"to bottom right": 135,
	"to right bottom": 135,
	"to bottom":       180,
	"to bottom left":  225,
	"to left bottom":  225,
	"to left":         270,
	"to top left":     315,
	"to left top":     315,
}

// ParseGradient parses a CSS gradient function string into a Gradient.
// ParseGradient parses CSS gradient syntax into a structured Gradient.
func ParseGradient(s string) (Gradient, error) {
	s = strings.TrimSpace(s)
	lower := strings.ToLower(s)

	var g Gradient
	var inner string

	switch {
	case strings.HasPrefix(lower, "repeating-linear-gradient(") && strings.HasSuffix(s, ")"):
		g.Type = RepeatingLinearGradient
		g.Repeating = true
		inner = s[len("repeating-linear-gradient(") : len(s)-1]
	case strings.HasPrefix(lower, "repeating-radial-gradient(") && strings.HasSuffix(s, ")"):
		g.Type = RepeatingRadialGradient
		g.Repeating = true
		inner = s[len("repeating-radial-gradient(") : len(s)-1]
	case strings.HasPrefix(lower, "linear-gradient(") && strings.HasSuffix(s, ")"):
		g.Type = LinearGradient
		inner = s[len("linear-gradient(") : len(s)-1]
	case strings.HasPrefix(lower, "radial-gradient(") && strings.HasSuffix(s, ")"):
		g.Type = RadialGradient
		inner = s[len("radial-gradient(") : len(s)-1]
	default:
		return Gradient{}, fmt.Errorf("unsupported gradient: %q", s)
	}

	inner = strings.TrimSpace(inner)
	args := splitGradientArgs(inner)

	if g.Type == LinearGradient || g.Type == RepeatingLinearGradient {
		return parseLinearGradient(g, args)
	}
	return parseRadialGradient(g, args)
}

func splitGradientArgs(s string) []string {
	var parts []string
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
				parts = append(parts, strings.TrimSpace(s[start:i]))
				start = i + 1
			}
		}
	}
	parts = append(parts, strings.TrimSpace(s[start:]))
	return parts
}

func parseLinearGradient(g Gradient, args []string) (Gradient, error) {
	if len(args) < 2 {
		return g, fmt.Errorf("linear-gradient requires at least 2 arguments")
	}

	g.Angle = 180
	stopStart := 0

	first := strings.TrimSpace(args[0])
	firstLower := strings.ToLower(first)

	if angle, ok := directionKeywords[firstLower]; ok {
		g.Angle = angle
		stopStart = 1
	} else if strings.HasSuffix(firstLower, "deg") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(firstLower, "deg"), 64)
		if err != nil {
			return g, fmt.Errorf("invalid angle: %q", first)
		}
		g.Angle = v
		stopStart = 1
	} else if strings.HasSuffix(firstLower, "rad") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(firstLower, "rad"), 64)
		if err != nil {
			return g, fmt.Errorf("invalid angle: %q", first)
		}
		g.Angle = v * 180 / math.Pi
		stopStart = 1
	} else if strings.HasSuffix(firstLower, "turn") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(firstLower, "turn"), 64)
		if err != nil {
			return g, fmt.Errorf("invalid angle: %q", first)
		}
		g.Angle = v * 360
		stopStart = 1
	}

	stops, err := parseColorStops(args[stopStart:])
	if err != nil {
		return g, err
	}
	g.Stops = stops
	return g, nil
}

func parseRadialGradient(g Gradient, args []string) (Gradient, error) {
	if len(args) < 2 {
		return g, fmt.Errorf("radial-gradient requires at least 2 arguments")
	}

	g.Shape = "ellipse"
	g.Size = "farthest-corner"
	g.PositionX = 50
	g.PositionY = 50

	stopStart := 0
	first := strings.TrimSpace(args[0])
	firstLower := strings.ToLower(first)

	if _, err := ParseColor(first); err != nil || isRadialConfig(firstLower) {
		if isRadialConfig(firstLower) {
			parseRadialConfig(&g, firstLower)
			stopStart = 1
		}
	}

	stops, err := parseColorStops(args[stopStart:])
	if err != nil {
		return g, err
	}
	g.Stops = stops
	return g, nil
}

func isRadialConfig(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	for _, kw := range []string{"circle", "ellipse", "closest-side", "closest-corner", "farthest-side", "farthest-corner", "at "} {
		if strings.Contains(s, kw) {
			return true
		}
	}
	return false
}

func parseRadialConfig(g *Gradient, s string) {
	atIdx := strings.Index(s, " at ")
	shapePart := s
	if atIdx >= 0 {
		shapePart = strings.TrimSpace(s[:atIdx])
		posPart := strings.TrimSpace(s[atIdx+4:])
		parseRadialPosition(g, posPart)
	}

	if shapePart == "" {
		return
	}

	words := strings.Fields(shapePart)
	for _, w := range words {
		switch w {
		case "circle":
			g.Shape = "circle"
		case "ellipse":
			g.Shape = "ellipse"
		case "closest-side", "closest-corner", "farthest-side", "farthest-corner":
			g.Size = w
		}
	}
}

func parseRadialPosition(g *Gradient, s string) {
	posKeywords := map[string]float64{
		"left":   0,
		"center": 50,
		"right":  100,
		"top":    0,
		"bottom": 100,
	}

	parts := strings.Fields(s)
	if len(parts) == 1 {
		if v, ok := posKeywords[parts[0]]; ok {
			g.PositionX = v
			g.PositionY = v
			if parts[0] == "center" {
				g.PositionX = 50
				g.PositionY = 50
			}
		} else if strings.HasSuffix(parts[0], "%") {
			if v, err := strconv.ParseFloat(strings.TrimSuffix(parts[0], "%"), 64); err == nil {
				g.PositionX = v
				g.PositionY = 50
			}
		}
	} else if len(parts) == 2 {
		if v, ok := posKeywords[parts[0]]; ok {
			g.PositionX = v
		} else if strings.HasSuffix(parts[0], "%") {
			if v, err := strconv.ParseFloat(strings.TrimSuffix(parts[0], "%"), 64); err == nil {
				g.PositionX = v
			}
		}
		if v, ok := posKeywords[parts[1]]; ok {
			g.PositionY = v
		} else if strings.HasSuffix(parts[1], "%") {
			if v, err := strconv.ParseFloat(strings.TrimSuffix(parts[1], "%"), 64); err == nil {
				g.PositionY = v
			}
		}
	}
}

func parseColorStops(args []string) ([]ColorStop, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("gradient requires at least one color stop")
	}

	var stops []ColorStop
	for _, arg := range args {
		stop, err := parseColorStop(strings.TrimSpace(arg))
		if err != nil {
			return nil, err
		}
		stops = append(stops, stop)
	}
	return stops, nil
}

func parseColorStop(s string) (ColorStop, error) {
	s = strings.TrimSpace(s)

	if c, err := ParseColor(s); err == nil {
		return ColorStop{Color: c, Position: math.NaN(), HasPos: false}, nil
	}

	funcEnd := findFuncEnd(s)
	if funcEnd > 0 {
		colorPart := s[:funcEnd]
		rest := strings.TrimSpace(s[funcEnd:])
		c, err := ParseColor(colorPart)
		if err != nil {
			return ColorStop{}, fmt.Errorf("invalid color in stop: %q", s)
		}
		if rest == "" {
			return ColorStop{Color: c, Position: math.NaN(), HasPos: false}, nil
		}
		pos, err := parseStopPosition(rest)
		if err != nil {
			return ColorStop{}, err
		}
		return ColorStop{Color: c, Position: pos, HasPos: true}, nil
	}

	lastSpace := strings.LastIndex(s, " ")
	if lastSpace < 0 {
		c, err := ParseColor(s)
		if err != nil {
			return ColorStop{}, fmt.Errorf("invalid color stop: %q", s)
		}
		return ColorStop{Color: c, Position: math.NaN(), HasPos: false}, nil
	}

	colorPart := s[:lastSpace]
	posPart := s[lastSpace+1:]

	pos, posErr := parseStopPosition(posPart)
	if posErr == nil {
		c, err := ParseColor(strings.TrimSpace(colorPart))
		if err != nil {
			return ColorStop{}, fmt.Errorf("invalid color in stop: %q", s)
		}
		return ColorStop{Color: c, Position: pos, HasPos: true}, nil
	}

	c, err := ParseColor(s)
	if err != nil {
		return ColorStop{}, fmt.Errorf("invalid color stop: %q", s)
	}
	return ColorStop{Color: c, Position: math.NaN(), HasPos: false}, nil
}

func findFuncEnd(s string) int {
	depth := 0
	inFunc := false
	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			depth++
			inFunc = true
		} else if s[i] == ')' {
			depth--
			if depth == 0 && inFunc {
				return i + 1
			}
		}
	}
	return 0
}

func parseStopPosition(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err != nil {
			return 0, err
		}
		return v / 100, nil
	}
	return 0, fmt.Errorf("unsupported stop position: %q", s)
}
