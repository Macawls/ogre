package style

import (
	"fmt"
	"strconv"
	"strings"
)

// Shadow represents a parsed CSS box-shadow or text-shadow value.
type Shadow struct {
	OffsetX float64
	OffsetY float64
	Blur    float64
	Spread  float64
	Color   Color
	Inset   bool
}

// ParseBoxShadow parses a CSS box-shadow property value into a slice of Shadows.
// ParseBoxShadow parses CSS box-shadow values.
func ParseBoxShadow(s string) ([]Shadow, error) {
	s = strings.TrimSpace(s)
	if s == "" || strings.ToLower(s) == "none" {
		return nil, nil
	}

	parts := splitShadowList(s)
	var shadows []Shadow
	for _, part := range parts {
		shadow, err := parseSingleBoxShadow(strings.TrimSpace(part))
		if err != nil {
			return nil, err
		}
		shadows = append(shadows, shadow)
	}
	return shadows, nil
}

// ParseTextShadow parses a CSS text-shadow property value into a slice of Shadows.
// ParseTextShadow parses CSS text-shadow values.
func ParseTextShadow(s string) ([]Shadow, error) {
	s = strings.TrimSpace(s)
	if s == "" || strings.ToLower(s) == "none" {
		return nil, nil
	}

	parts := splitShadowList(s)
	var shadows []Shadow
	for _, part := range parts {
		shadow, err := parseSingleTextShadow(strings.TrimSpace(part))
		if err != nil {
			return nil, err
		}
		shadows = append(shadows, shadow)
	}
	return shadows, nil
}

func splitShadowList(s string) []string {
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
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, s[start:])
	return parts
}

func parseSingleBoxShadow(s string) (Shadow, error) {
	tokens := tokenizeShadow(s)
	if len(tokens) < 2 {
		return Shadow{}, fmt.Errorf("box-shadow requires at least 2 values: %q", s)
	}

	var shadow Shadow
	var lengths []float64
	var colorParts []string

	for _, tok := range tokens {
		lower := strings.ToLower(tok)
		if lower == "inset" {
			shadow.Inset = true
			continue
		}
		if v, err := parseShadowLength(tok); err == nil {
			lengths = append(lengths, v)
			continue
		}
		colorParts = append(colorParts, tok)
	}

	switch len(lengths) {
	case 2:
		shadow.OffsetX = lengths[0]
		shadow.OffsetY = lengths[1]
	case 3:
		shadow.OffsetX = lengths[0]
		shadow.OffsetY = lengths[1]
		shadow.Blur = lengths[2]
	case 4:
		shadow.OffsetX = lengths[0]
		shadow.OffsetY = lengths[1]
		shadow.Blur = lengths[2]
		shadow.Spread = lengths[3]
	default:
		return Shadow{}, fmt.Errorf("box-shadow expects 2-4 length values, got %d: %q", len(lengths), s)
	}

	if len(colorParts) > 0 {
		colorStr := strings.Join(colorParts, " ")
		c, err := ParseColor(colorStr)
		if err != nil {
			return Shadow{}, fmt.Errorf("invalid shadow color %q: %w", colorStr, err)
		}
		shadow.Color = c
	} else {
		shadow.Color = Color{0, 0, 0, 1}
	}

	return shadow, nil
}

func parseSingleTextShadow(s string) (Shadow, error) {
	tokens := tokenizeShadow(s)
	if len(tokens) < 2 {
		return Shadow{}, fmt.Errorf("text-shadow requires at least 2 values: %q", s)
	}

	var shadow Shadow
	var lengths []float64
	var colorParts []string

	for _, tok := range tokens {
		if v, err := parseShadowLength(tok); err == nil {
			lengths = append(lengths, v)
			continue
		}
		colorParts = append(colorParts, tok)
	}

	switch len(lengths) {
	case 2:
		shadow.OffsetX = lengths[0]
		shadow.OffsetY = lengths[1]
	case 3:
		shadow.OffsetX = lengths[0]
		shadow.OffsetY = lengths[1]
		shadow.Blur = lengths[2]
	default:
		return Shadow{}, fmt.Errorf("text-shadow expects 2-3 length values, got %d: %q", len(lengths), s)
	}

	if len(colorParts) > 0 {
		colorStr := strings.Join(colorParts, " ")
		c, err := ParseColor(colorStr)
		if err != nil {
			return Shadow{}, fmt.Errorf("invalid shadow color %q: %w", colorStr, err)
		}
		shadow.Color = c
	} else {
		shadow.Color = Color{0, 0, 0, 1}
	}

	return shadow, nil
}

func tokenizeShadow(s string) []string {
	var tokens []string
	s = strings.TrimSpace(s)
	i := 0
	for i < len(s) {
		for i < len(s) && s[i] == ' ' {
			i++
		}
		if i >= len(s) {
			break
		}
		start := i
		if s[i] == '-' || s[i] == '+' || s[i] == '.' || (s[i] >= '0' && s[i] <= '9') {
			for i < len(s) && s[i] != ' ' && s[i] != '(' {
				i++
			}
			tokens = append(tokens, s[start:i])
		} else {
			depth := 0
			for i < len(s) {
				if s[i] == '(' {
					depth++
				} else if s[i] == ')' {
					depth--
					if depth == 0 {
						i++
						break
					}
				} else if s[i] == ' ' && depth == 0 {
					break
				}
				i++
			}
			tokens = append(tokens, s[start:i])
		}
	}
	return tokens
}

func parseShadowLength(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "0" {
		return 0, nil
	}
	for _, suffix := range []string{"px", "em", "rem", "pt", "vh", "vw"} {
		if strings.HasSuffix(s, suffix) {
			v, err := strconv.ParseFloat(s[:len(s)-len(suffix)], 64)
			if err != nil {
				return 0, err
			}
			return v, nil
		}
	}
	if len(s) > 0 && (s[0] == '-' || s[0] == '+' || s[0] == '.' || (s[0] >= '0' && s[0] <= '9')) {
		if _, err := strconv.ParseFloat(s, 64); err == nil {
			return 0, fmt.Errorf("unitless non-zero number: %q", s)
		}
	}
	return 0, fmt.Errorf("not a length: %q", s)
}
