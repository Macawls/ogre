package font

import (
	"strconv"
	"strings"

	gotextfont "github.com/go-text/typesetting/font"
	ot "github.com/go-text/typesetting/font/opentype"
)

type Variation struct {
	Tag   string
	Value float32
}

func variationKey(vars []Variation) string {
	if len(vars) == 0 {
		return ""
	}
	var b strings.Builder
	for i, v := range vars {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString(v.Tag)
		b.WriteByte('=')
		b.WriteString(strconv.FormatFloat(float64(v.Value), 'f', -1, 32))
	}
	return b.String()
}

func toGoTextVariations(vars []Variation) []gotextfont.Variation {
	if len(vars) == 0 {
		return nil
	}
	out := make([]gotextfont.Variation, 0, len(vars))
	for _, v := range vars {
		if len(v.Tag) != 4 {
			continue
		}
		out = append(out, gotextfont.Variation{
			Tag:   ot.MustNewTag(v.Tag),
			Value: v.Value,
		})
	}
	return out
}
