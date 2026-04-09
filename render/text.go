package render

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/bidi"

	"github.com/macawls/ogre/font"
	"github.com/macawls/ogre/style"
)

type TextRenderResult struct {
	Shadows     string
	Content     string
	Decorations string
}

// RenderText generates the corresponding output format.
func RenderText(lines []font.TextLine, cs *style.ComputedStyle, boxX, boxY, boxW, boxH float64) TextRenderResult {
	return RenderTextWithIDGen(lines, cs, boxX, boxY, boxW, boxH, nil, nil, nil)
}

// RenderTextWithIDGen generates the corresponding output format.
func RenderTextWithIDGen(lines []font.TextLine, cs *style.ComputedStyle, boxX, boxY, boxW, boxH float64, idGen func(string) string, fontMgr *font.Manager, emojiProvider ...*font.EmojiProvider) TextRenderResult {
	if len(lines) == 0 {
		return TextRenderResult{}
	}

	family := cs.FontFamily
	if family == "" {
		family = "sans-serif"
	}
	size := cs.FontSize
	if size == 0 {
		size = 16
	}
	weight := cs.FontWeight
	if weight == 0 {
		weight = 400
	}
	lineHeight := cs.LineHeight
	if lineHeight == 0 {
		lineHeight = size * 1.2
	}

	fill := colorToCSS(cs.Color)
	if cs.Color.A == -1 {
		fill = "#000000"
	}

	ascent := size * 0.8
	descent := size * 0.2

	if fontMgr != nil {
		fontStyle := cs.FontStyle
		if fontStyle == "" {
			fontStyle = "normal"
		}
		face := fontMgr.Resolve(family, weight, fontStyle)
		if face != nil {
			ff, err := fontMgr.NewFace(face, size)
			if err == nil {
				ascent = font.Ascent(ff)
				descent = font.Descent(ff)
			}
		}
	}

	var shadows strings.Builder
	var content strings.Builder
	var decorations strings.Builder

	var textShadows []style.Shadow
	if cs.TextShadow != "" {
		textShadows, _ = style.ParseTextShadow(cs.TextShadow)
	}

	rtl := isRTL(cs)

	for i, line := range lines {
		text := applyTextTransform(line.Text, cs.TextTransform)
		if rtl {
			text = reorderBidi(text, "rtl")
		}
		align := cs.TextAlign
		if rtl && align == style.TextAlignStart {
			align = style.TextAlignRight
		}
		x := alignX(boxX, boxW, line.Width, align)
		y := boxY + ascent + float64(i)*lineHeight

		for _, s := range textShadows {
			sx := x + s.OffsetX
			sy := y + s.OffsetY
			shadowFill := shadowColorCSS(s.Color)

			shadowPathRendered := false
			if fontMgr != nil {
				fontStyle := cs.FontStyle
				if fontStyle == "" {
					fontStyle = "normal"
				}
				pathD, _ := font.TextToPath(fontMgr, text, family, weight, fontStyle, size)
				if pathD != "" {
					if s.Blur > 0 && idGen != nil {
						filterID := idGen("tshadow")
						fmt.Fprintf(&shadows, `<defs><filter id="%s" x="-50%%" y="-50%%" width="200%%" height="200%%">`+
							`<feGaussianBlur stdDeviation="%.4g"/>`+
							`</filter></defs>`, filterID, s.Blur/2)
						fmt.Fprintf(&shadows, `<path d="%s" fill="%s" filter="url(#%s)" transform="translate(%.4g,%.4g)"/>`,
							pathD, shadowFill, filterID, sx, sy)
					} else {
						fmt.Fprintf(&shadows, `<path d="%s" fill="%s" transform="translate(%.4g,%.4g)"/>`,
							pathD, shadowFill, sx, sy)
					}
					shadowPathRendered = true
				}
			}
			if !shadowPathRendered {
				if s.Blur > 0 && idGen != nil {
					filterID := idGen("tshadow")
					fmt.Fprintf(&shadows, `<defs><filter id="%s" x="-50%%" y="-50%%" width="200%%" height="200%%">`+
						`<feGaussianBlur stdDeviation="%.4g"/>`+
						`</filter></defs>`, filterID, s.Blur/2)
					fmt.Fprintf(&shadows, `<text x="%.4g" y="%.4g" font-family="%s" font-size="%.4g" font-weight="%d"`,
						sx, sy, xmlEscape(family), size, weight)
					if cs.FontStyle != "" && cs.FontStyle != "normal" {
						fmt.Fprintf(&shadows, ` font-style="%s"`, cs.FontStyle)
					}
					fmt.Fprintf(&shadows, ` fill="%s" filter="url(#%s)"`, shadowFill, filterID)
					if cs.LetterSpacing != 0 {
						fmt.Fprintf(&shadows, ` letter-spacing="%.4g"`, cs.LetterSpacing)
					}
					shadows.WriteString(">")
					shadows.WriteString(xmlEscape(text))
					shadows.WriteString("</text>")
				} else {
					fmt.Fprintf(&shadows, `<text x="%.4g" y="%.4g" font-family="%s" font-size="%.4g" font-weight="%d"`,
						sx, sy, xmlEscape(family), size, weight)
					if cs.FontStyle != "" && cs.FontStyle != "normal" {
						fmt.Fprintf(&shadows, ` font-style="%s"`, cs.FontStyle)
					}
					fmt.Fprintf(&shadows, ` fill="%s"`, shadowFill)
					if cs.LetterSpacing != 0 {
						fmt.Fprintf(&shadows, ` letter-spacing="%.4g"`, cs.LetterSpacing)
					}
					shadows.WriteString(">")
					shadows.WriteString(xmlEscape(text))
					shadows.WriteString("</text>")
				}
			}
		}

		var ep *font.EmojiProvider
		if len(emojiProvider) > 0 {
			ep = emojiProvider[0]
		}

		if ep != nil && containsEmoji(text) {
			segments := font.SplitEmoji(text)
			cx := x
			for _, seg := range segments {
				if seg.IsEmoji {
					emojiSize := size
					ey := y - ascent
					href := font.TwemojiURL(seg.Text)
					fmt.Fprintf(&content, `<image href="%s" x="%.4g" y="%.4g" width="%.4g" height="%.4g"/>`,
						href, cx, ey, emojiSize, emojiSize)
					cx += emojiSize
				} else {
					renderTextSegment(&content, seg.Text, cx, y, family, size, weight, cs, fill, fontMgr)
					if fontMgr != nil {
						face := fontMgr.Resolve(family, weight, cs.FontStyle)
						if face != nil {
							ff, err := fontMgr.NewFace(face, size)
							if err == nil {
								m := font.NewMeasurer(ff, cs.LetterSpacing)
								cx += m.StringWidth(seg.Text)
							}
						}
					}
				}
			}
		} else {
			renderTextSegment(&content, text, x, y, family, size, weight, cs, fill, fontMgr)
		}

		if cs.TextDecorationLine != style.TextDecorationNone {
			renderDecoration(&decorations, cs, x, y, line.Width, ascent, descent)
		}
	}

	return TextRenderResult{
		Shadows:     shadows.String(),
		Content:     content.String(),
		Decorations: decorations.String(),
	}
}

func alignX(boxX, boxW, lineWidth float64, align style.TextAlign) float64 {
	switch align {
	case style.TextAlignRight, style.TextAlignEnd:
		return boxX + boxW - lineWidth
	case style.TextAlignCenter:
		return boxX + (boxW-lineWidth)/2
	default:
		return boxX
	}
}

func applyTextTransform(s string, t style.TextTransform) string {
	switch t {
	case style.TextTransformUppercase:
		return strings.ToUpper(s)
	case style.TextTransformLowercase:
		return strings.ToLower(s)
	case style.TextTransformCapitalize:
		return capitalize(s)
	default:
		return s
	}
}

func capitalize(s string) string {
	var b strings.Builder
	wordStart := true
	for _, r := range s {
		if unicode.IsSpace(r) {
			b.WriteRune(r)
			wordStart = true
		} else if wordStart {
			b.WriteRune(unicode.ToUpper(r))
			wordStart = false
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func renderDecoration(b *strings.Builder, cs *style.ComputedStyle, x, baseline, lineWidth, ascent, descent float64) {
	stroke := colorToCSS(cs.Color)
	if !cs.TextDecorationColor.IsTransparent() {
		stroke = colorToCSS(cs.TextDecorationColor)
	}

	var dashArray string
	switch cs.TextDecorationStyle {
	case "dashed":
		dashArray = ` stroke-dasharray="6 3"`
	case "dotted":
		dashArray = ` stroke-dasharray="2 2"`
	}

	var dy float64
	switch cs.TextDecorationLine {
	case style.TextDecorationUnderline:
		dy = baseline + descent/2
	case style.TextDecorationLineThrough:
		dy = baseline - ascent/3
	case style.TextDecorationOverline:
		dy = baseline - ascent
	default:
		return
	}

	fmt.Fprintf(b, `<line x1="%.4g" y1="%.4g" x2="%.4g" y2="%.4g" stroke="%s" stroke-width="1"%s/>`,
		x, dy, x+lineWidth, dy, stroke, dashArray)
}

func renderTextSegment(content *strings.Builder, text string, x, y float64, family string, size float64, weight int, cs *style.ComputedStyle, fill string, fontMgr *font.Manager) {
	pathRendered := false
	rtl := cs.Direction == "rtl"
	if fontMgr != nil {
		fontStyle := cs.FontStyle
		if fontStyle == "" {
			fontStyle = "normal"
		}
		var pathD string
		if rtl || needsShaping(text) {
			pathD, _ = font.ShapedTextToPath(fontMgr, text, family, weight, fontStyle, size, rtl)
		} else {
			pathD, _ = font.TextToPath(fontMgr, text, family, weight, fontStyle, size)
		}
		if pathD != "" {
			fmt.Fprintf(content, `<path d="%s" fill="%s" transform="translate(%.4g,%.4g)"/>`,
				pathD, fill, x, y)
			pathRendered = true
		}
	}
	if !pathRendered {
		fmt.Fprintf(content, `<text x="%.4g" y="%.4g" font-family="%s" font-size="%.4g" font-weight="%d"`,
			x, y, xmlEscape(family), size, weight)
		if cs.FontStyle != "" && cs.FontStyle != "normal" {
			fmt.Fprintf(content, ` font-style="%s"`, cs.FontStyle)
		}
		fmt.Fprintf(content, ` fill="%s"`, fill)
		if cs.LetterSpacing != 0 {
			fmt.Fprintf(content, ` letter-spacing="%.4g"`, cs.LetterSpacing)
		}
		content.WriteString(">")
		content.WriteString(xmlEscape(text))
		content.WriteString("</text>")
	}
}

func reorderBidi(text string, dir string) string {
	var defaultDir bidi.Direction
	if dir == "rtl" {
		defaultDir = bidi.RightToLeft
	} else {
		defaultDir = bidi.LeftToRight
	}

	p := bidi.Paragraph{}
	p.SetString(text, bidi.DefaultDirection(defaultDir))
	ordering, err := p.Order()
	if err != nil {
		return text
	}

	var result strings.Builder
	for i := 0; i < ordering.NumRuns(); i++ {
		run := ordering.Run(i)
		result.WriteString(run.String())
	}
	return result.String()
}

func isRTL(cs *style.ComputedStyle) bool {
	return cs.Direction == "rtl"
}

func needsShaping(s string) bool {
	for _, r := range s {
		if r >= 0x0590 && r <= 0x05FF {
			return true // Hebrew
		}
		if r >= 0x0600 && r <= 0x06FF {
			return true // Arabic
		}
		if r >= 0x0700 && r <= 0x074F {
			return true // Syriac
		}
		if r >= 0x0780 && r <= 0x07BF {
			return true // Thaana
		}
		if r >= 0x0900 && r <= 0x097F {
			return true // Devanagari
		}
		if r >= 0x0E00 && r <= 0x0E7F {
			return true // Thai
		}
		if r >= 0xFB50 && r <= 0xFDFF {
			return true // Arabic Presentation Forms-A
		}
		if r >= 0xFE70 && r <= 0xFEFF {
			return true // Arabic Presentation Forms-B
		}
	}
	return false
}

func containsEmoji(s string) bool {
	for _, r := range s {
		if font.IsEmoji(r) {
			return true
		}
	}
	return false
}
