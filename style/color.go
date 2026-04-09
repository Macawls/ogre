// Package style defines CSS property types, parsing, inheritance, and computed style resolution.
package style

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Color represents an RGBA color value.
type Color struct {
	R, G, B uint8
	A        float64
}

// CurrentColor is a sentinel value representing the CSS currentColor keyword.
var CurrentColor = Color{R: 0, G: 0, B: 0, A: -1}

func (c Color) String() string {
	return fmt.Sprintf("rgba(%d,%d,%d,%.4g)", c.R, c.G, c.B, c.A)
}

func (c Color) Hex() string {
	if c.A < 1.0 {
		a := uint8(math.Round(c.A * 255))
		return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, a)
	}
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

func (c Color) IsTransparent() bool {
	return c.A == 0
}

// ParseColor parses a CSS color string (hex, rgb, hsl, or named) into a Color.
func ParseColor(s string) (Color, error) {
	s = strings.TrimSpace(s)
	lower := strings.ToLower(s)

	if lower == "currentcolor" {
		return CurrentColor, nil
	}

	if lower == "transparent" {
		return Color{0, 0, 0, 0}, nil
	}

	if c, ok := namedColors[lower]; ok {
		return c, nil
	}

	if strings.HasPrefix(s, "#") {
		return parseHex(s[1:])
	}

	if strings.HasPrefix(lower, "rgb") {
		return parseRGB(s)
	}

	if strings.HasPrefix(lower, "hsl") {
		return parseHSL(s)
	}

	return Color{}, fmt.Errorf("unsupported color: %q", s)
}

func parseHex(hex string) (Color, error) {
	switch len(hex) {
	case 3:
		r, err := strconv.ParseUint(string(hex[0])+string(hex[0]), 16, 8)
		if err != nil {
			return Color{}, err
		}
		g, err := strconv.ParseUint(string(hex[1])+string(hex[1]), 16, 8)
		if err != nil {
			return Color{}, err
		}
		b, err := strconv.ParseUint(string(hex[2])+string(hex[2]), 16, 8)
		if err != nil {
			return Color{}, err
		}
		return Color{uint8(r), uint8(g), uint8(b), 1}, nil
	case 4:
		r, err := strconv.ParseUint(string(hex[0])+string(hex[0]), 16, 8)
		if err != nil {
			return Color{}, err
		}
		g, err := strconv.ParseUint(string(hex[1])+string(hex[1]), 16, 8)
		if err != nil {
			return Color{}, err
		}
		b, err := strconv.ParseUint(string(hex[2])+string(hex[2]), 16, 8)
		if err != nil {
			return Color{}, err
		}
		a, err := strconv.ParseUint(string(hex[3])+string(hex[3]), 16, 8)
		if err != nil {
			return Color{}, err
		}
		return Color{uint8(r), uint8(g), uint8(b), float64(a) / 255}, nil
	case 6:
		r, err := strconv.ParseUint(hex[0:2], 16, 8)
		if err != nil {
			return Color{}, err
		}
		g, err := strconv.ParseUint(hex[2:4], 16, 8)
		if err != nil {
			return Color{}, err
		}
		b, err := strconv.ParseUint(hex[4:6], 16, 8)
		if err != nil {
			return Color{}, err
		}
		return Color{uint8(r), uint8(g), uint8(b), 1}, nil
	case 8:
		r, err := strconv.ParseUint(hex[0:2], 16, 8)
		if err != nil {
			return Color{}, err
		}
		g, err := strconv.ParseUint(hex[2:4], 16, 8)
		if err != nil {
			return Color{}, err
		}
		b, err := strconv.ParseUint(hex[4:6], 16, 8)
		if err != nil {
			return Color{}, err
		}
		a, err := strconv.ParseUint(hex[6:8], 16, 8)
		if err != nil {
			return Color{}, err
		}
		return Color{uint8(r), uint8(g), uint8(b), float64(a) / 255}, nil
	default:
		return Color{}, fmt.Errorf("invalid hex color length: %d", len(hex))
	}
}

func parseRGB(s string) (Color, error) {
	lower := strings.ToLower(s)
	inner, err := extractFuncArgs(lower, "rgba", "rgb")
	if err != nil {
		return Color{}, err
	}

	parts, alpha := splitColorArgs(inner)
	if len(parts) != 3 {
		return Color{}, fmt.Errorf("rgb() requires 3 color components, got %d", len(parts))
	}

	var r, g, b uint8
	if strings.HasSuffix(parts[0], "%") {
		rv, err := parsePercent(parts[0])
		if err != nil {
			return Color{}, err
		}
		gv, err := parsePercent(parts[1])
		if err != nil {
			return Color{}, err
		}
		bv, err := parsePercent(parts[2])
		if err != nil {
			return Color{}, err
		}
		r = clampByte(rv / 100 * 255)
		g = clampByte(gv / 100 * 255)
		b = clampByte(bv / 100 * 255)
	} else {
		rv, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return Color{}, err
		}
		gv, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return Color{}, err
		}
		bv, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return Color{}, err
		}
		r = clampByte(rv)
		g = clampByte(gv)
		b = clampByte(bv)
	}

	a := 1.0
	if alpha != "" {
		a, err = parseAlpha(alpha)
		if err != nil {
			return Color{}, err
		}
	}

	return Color{r, g, b, a}, nil
}

func parseHSL(s string) (Color, error) {
	lower := strings.ToLower(s)
	inner, err := extractFuncArgs(lower, "hsla", "hsl")
	if err != nil {
		return Color{}, err
	}

	parts, alpha := splitColorArgs(inner)
	if len(parts) != 3 {
		return Color{}, fmt.Errorf("hsl() requires 3 components, got %d", len(parts))
	}

	h, err := strconv.ParseFloat(strings.TrimSuffix(parts[0], "deg"), 64)
	if err != nil {
		return Color{}, fmt.Errorf("invalid hue: %w", err)
	}

	sv, err := parsePercent(parts[1])
	if err != nil {
		return Color{}, fmt.Errorf("invalid saturation: %w", err)
	}

	lv, err := parsePercent(parts[2])
	if err != nil {
		return Color{}, fmt.Errorf("invalid lightness: %w", err)
	}

	h = math.Mod(h, 360)
	if h < 0 {
		h += 360
	}
	sat := clampFloat(sv/100, 0, 1)
	light := clampFloat(lv/100, 0, 1)

	r, g, b := hslToRGB(h, sat, light)

	a := 1.0
	if alpha != "" {
		a, err = parseAlpha(alpha)
		if err != nil {
			return Color{}, err
		}
	}

	return Color{
		R: clampByte(r * 255),
		G: clampByte(g * 255),
		B: clampByte(b * 255),
		A: a,
	}, nil
}

func extractFuncArgs(s string, names ...string) (string, error) {
	for _, name := range names {
		if strings.HasPrefix(s, name+"(") && strings.HasSuffix(s, ")") {
			return s[len(name)+1 : len(s)-1], nil
		}
	}
	return "", fmt.Errorf("not a valid color function: %q", s)
}

func splitColorArgs(inner string) (components []string, alpha string) {
	inner = strings.TrimSpace(inner)

	if idx := strings.Index(inner, "/"); idx >= 0 {
		alpha = strings.TrimSpace(inner[idx+1:])
		inner = strings.TrimSpace(inner[:idx])
	}

	if strings.Contains(inner, ",") {
		raw := strings.Split(inner, ",")
		for _, r := range raw {
			r = strings.TrimSpace(r)
			if r != "" {
				components = append(components, r)
			}
		}
		if alpha == "" && len(components) == 4 {
			alpha = components[3]
			components = components[:3]
		}
		return
	}

	fields := strings.Fields(inner)
	components = fields
	return
}

func parsePercent(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if !strings.HasSuffix(s, "%") {
		return 0, fmt.Errorf("expected percentage: %q", s)
	}
	return strconv.ParseFloat(s[:len(s)-1], 64)
}

func parseAlpha(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "%") {
		v, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err != nil {
			return 0, err
		}
		return clampFloat(v/100, 0, 1), nil
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return clampFloat(v, 0, 1), nil
}

func clampByte(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(math.Round(v))
}

func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func hslToRGB(h, s, l float64) (float64, float64, float64) {
	if s == 0 {
		return l, l, l
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	h /= 360

	r := hueToRGB(p, q, h+1.0/3.0)
	g := hueToRGB(p, q, h)
	b := hueToRGB(p, q, h-1.0/3.0)

	return r, g, b
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}

var namedColors = map[string]Color{
	"aliceblue":            {240, 248, 255, 1},
	"antiquewhite":         {250, 235, 215, 1},
	"aqua":                 {0, 255, 255, 1},
	"aquamarine":           {127, 255, 212, 1},
	"azure":                {240, 255, 255, 1},
	"beige":                {245, 245, 220, 1},
	"bisque":               {255, 228, 196, 1},
	"black":                {0, 0, 0, 1},
	"blanchedalmond":       {255, 235, 205, 1},
	"blue":                 {0, 0, 255, 1},
	"blueviolet":           {138, 43, 226, 1},
	"brown":                {165, 42, 42, 1},
	"burlywood":            {222, 184, 135, 1},
	"cadetblue":            {95, 158, 160, 1},
	"chartreuse":           {127, 255, 0, 1},
	"chocolate":            {210, 105, 30, 1},
	"coral":                {255, 127, 80, 1},
	"cornflowerblue":       {100, 149, 237, 1},
	"cornsilk":             {255, 248, 220, 1},
	"crimson":              {220, 20, 60, 1},
	"cyan":                 {0, 255, 255, 1},
	"darkblue":             {0, 0, 139, 1},
	"darkcyan":             {0, 139, 139, 1},
	"darkgoldenrod":        {184, 134, 11, 1},
	"darkgray":             {169, 169, 169, 1},
	"darkgreen":            {0, 100, 0, 1},
	"darkgrey":             {169, 169, 169, 1},
	"darkkhaki":            {189, 183, 107, 1},
	"darkmagenta":          {139, 0, 139, 1},
	"darkolivegreen":       {85, 107, 47, 1},
	"darkorange":           {255, 140, 0, 1},
	"darkorchid":           {153, 50, 204, 1},
	"darkred":              {139, 0, 0, 1},
	"darksalmon":           {233, 150, 122, 1},
	"darkseagreen":         {143, 188, 143, 1},
	"darkslateblue":        {72, 61, 139, 1},
	"darkslategray":        {47, 79, 79, 1},
	"darkslategrey":        {47, 79, 79, 1},
	"darkturquoise":        {0, 206, 209, 1},
	"darkviolet":           {148, 0, 211, 1},
	"deeppink":             {255, 20, 147, 1},
	"deepskyblue":          {0, 191, 255, 1},
	"dimgray":              {105, 105, 105, 1},
	"dimgrey":              {105, 105, 105, 1},
	"dodgerblue":           {30, 144, 255, 1},
	"firebrick":            {178, 34, 34, 1},
	"floralwhite":          {255, 250, 240, 1},
	"forestgreen":          {34, 139, 34, 1},
	"fuchsia":              {255, 0, 255, 1},
	"gainsboro":            {220, 220, 220, 1},
	"ghostwhite":           {248, 248, 255, 1},
	"gold":                 {255, 215, 0, 1},
	"goldenrod":            {218, 165, 32, 1},
	"gray":                 {128, 128, 128, 1},
	"green":                {0, 128, 0, 1},
	"greenyellow":          {173, 255, 47, 1},
	"grey":                 {128, 128, 128, 1},
	"honeydew":             {240, 255, 240, 1},
	"hotpink":              {255, 105, 180, 1},
	"indianred":            {205, 92, 92, 1},
	"indigo":               {75, 0, 130, 1},
	"ivory":                {255, 255, 240, 1},
	"khaki":                {240, 230, 140, 1},
	"lavender":             {230, 230, 250, 1},
	"lavenderblush":        {255, 240, 245, 1},
	"lawngreen":            {124, 252, 0, 1},
	"lemonchiffon":         {255, 250, 205, 1},
	"lightblue":            {173, 216, 230, 1},
	"lightcoral":           {240, 128, 128, 1},
	"lightcyan":            {224, 255, 255, 1},
	"lightgoldenrodyellow": {250, 250, 210, 1},
	"lightgray":            {211, 211, 211, 1},
	"lightgreen":           {144, 238, 144, 1},
	"lightgrey":            {211, 211, 211, 1},
	"lightpink":            {255, 182, 193, 1},
	"lightsalmon":          {255, 160, 122, 1},
	"lightseagreen":        {32, 178, 170, 1},
	"lightskyblue":         {135, 206, 250, 1},
	"lightslategray":       {119, 136, 153, 1},
	"lightslategrey":       {119, 136, 153, 1},
	"lightsteelblue":       {176, 196, 222, 1},
	"lightyellow":          {255, 255, 224, 1},
	"lime":                 {0, 255, 0, 1},
	"limegreen":            {50, 205, 50, 1},
	"linen":                {250, 240, 230, 1},
	"magenta":              {255, 0, 255, 1},
	"maroon":               {128, 0, 0, 1},
	"mediumaquamarine":     {102, 205, 170, 1},
	"mediumblue":           {0, 0, 205, 1},
	"mediumorchid":         {186, 85, 211, 1},
	"mediumpurple":         {147, 112, 219, 1},
	"mediumseagreen":       {60, 179, 113, 1},
	"mediumslateblue":      {123, 104, 238, 1},
	"mediumspringgreen":    {0, 250, 154, 1},
	"mediumturquoise":      {72, 209, 204, 1},
	"mediumvioletred":      {199, 21, 133, 1},
	"midnightblue":         {25, 25, 112, 1},
	"mintcream":            {245, 255, 250, 1},
	"mistyrose":            {255, 228, 225, 1},
	"moccasin":             {255, 228, 181, 1},
	"navajowhite":          {255, 222, 173, 1},
	"navy":                 {0, 0, 128, 1},
	"oldlace":              {253, 245, 230, 1},
	"olive":                {128, 128, 0, 1},
	"olivedrab":            {107, 142, 35, 1},
	"orange":               {255, 165, 0, 1},
	"orangered":            {255, 69, 0, 1},
	"orchid":               {218, 112, 214, 1},
	"palegoldenrod":        {238, 232, 170, 1},
	"palegreen":            {152, 251, 152, 1},
	"paleturquoise":        {175, 238, 238, 1},
	"palevioletred":        {219, 112, 147, 1},
	"papayawhip":           {255, 239, 213, 1},
	"peachpuff":            {255, 218, 185, 1},
	"peru":                 {205, 133, 63, 1},
	"pink":                 {255, 192, 203, 1},
	"plum":                 {221, 160, 221, 1},
	"powderblue":           {176, 224, 230, 1},
	"purple":               {128, 0, 128, 1},
	"rebeccapurple":        {102, 51, 153, 1},
	"red":                  {255, 0, 0, 1},
	"rosybrown":            {188, 143, 143, 1},
	"royalblue":            {65, 105, 225, 1},
	"saddlebrown":          {139, 69, 19, 1},
	"salmon":               {250, 128, 114, 1},
	"sandybrown":           {244, 164, 96, 1},
	"seagreen":             {46, 139, 87, 1},
	"seashell":             {255, 245, 238, 1},
	"sienna":               {160, 82, 45, 1},
	"silver":               {192, 192, 192, 1},
	"skyblue":              {135, 206, 235, 1},
	"slateblue":            {106, 90, 205, 1},
	"slategray":            {112, 128, 144, 1},
	"slategrey":            {112, 128, 144, 1},
	"snow":                 {255, 250, 250, 1},
	"springgreen":          {0, 255, 127, 1},
	"steelblue":            {70, 130, 180, 1},
	"tan":                  {210, 180, 140, 1},
	"teal":                 {0, 128, 128, 1},
	"thistle":              {216, 191, 216, 1},
	"tomato":               {255, 99, 71, 1},
	"turquoise":            {64, 224, 208, 1},
	"violet":               {238, 130, 238, 1},
	"wheat":                {245, 222, 179, 1},
	"white":                {255, 255, 255, 1},
	"whitesmoke":           {245, 245, 245, 1},
	"yellow":               {255, 255, 0, 1},
	"yellowgreen":          {154, 205, 50, 1},
}
