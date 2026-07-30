package style

import (
	"strconv"
	"strings"
)

type VariationAxis struct {
	Tag   string
	Value float32
}

// ParseFontVariationSettings returns nil for "normal" or entirely malformed
// input; malformed axes inside an otherwise valid list are dropped silently.
func ParseFontVariationSettings(s string) []VariationAxis {
	s = strings.TrimSpace(s)
	if s == "" || strings.EqualFold(s, "normal") {
		return nil
	}

	parts := splitCSSList(s)
	if len(parts) == 0 {
		return nil
	}

	axes := make([]VariationAxis, 0, len(parts))
	for _, part := range parts {
		axis, ok := parseVariationAxis(part)
		if ok {
			axes = append(axes, axis)
		}
	}
	if len(axes) == 0 {
		return nil
	}
	return axes
}

func FormatFontVariationSettings(axes []VariationAxis) string {
	if len(axes) == 0 {
		return ""
	}
	var b strings.Builder
	for i, a := range axes {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(`'`)
		b.WriteString(a.Tag)
		b.WriteString(`' `)
		b.WriteString(strconv.FormatFloat(float64(a.Value), 'f', -1, 32))
	}
	return b.String()
}

func parseVariationAxis(s string) (VariationAxis, bool) {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return VariationAxis{}, false
	}

	quote := s[0]
	if quote != '\'' && quote != '"' {
		return VariationAxis{}, false
	}
	end := strings.IndexByte(s[1:], quote)
	if end < 0 {
		return VariationAxis{}, false
	}
	tag := s[1 : 1+end]
	if len(tag) != 4 {
		return VariationAxis{}, false
	}

	rest := strings.TrimSpace(s[1+end+1:])
	if rest == "" {
		return VariationAxis{}, false
	}
	v, err := strconv.ParseFloat(rest, 32)
	if err != nil {
		return VariationAxis{}, false
	}
	return VariationAxis{Tag: tag, Value: float32(v)}, true
}

func splitCSSList(s string) []string {
	var parts []string
	var buf strings.Builder
	inSingle := false
	inDouble := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '\'' && !inDouble:
			inSingle = !inSingle
			buf.WriteByte(c)
		case c == '"' && !inSingle:
			inDouble = !inDouble
			buf.WriteByte(c)
		case c == ',' && !inSingle && !inDouble:
			part := strings.TrimSpace(buf.String())
			if part != "" {
				parts = append(parts, part)
			}
			buf.Reset()
		default:
			buf.WriteByte(c)
		}
	}
	if buf.Len() > 0 {
		part := strings.TrimSpace(buf.String())
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}
