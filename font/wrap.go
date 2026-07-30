package font

import (
	"strings"
	"unicode"

	"golang.org/x/image/font"
)

// TextLine represents a single line of wrapped text with its measured width and position.
type TextLine struct {
	Text  string
	Width float64
	X     float64
	Y     float64
}

// WrapConfig controls how text is wrapped including max width, font, and white-space behavior.
type WrapConfig struct {
	MaxWidth      float64
	FontFace      font.Face
	FontSize      float64
	LineHeight    float64
	LetterSpacing float64
	WhiteSpace    int
	WordBreak     int
	LineClamp     int
	TextOverflow  string
	Face          *Face
	Variations    []Variation
}

const (
	wsNormal  = 0
	wsNowrap  = 1
	wsPre     = 2
	wsPreWrap = 3
	wsPreLine = 4

	wbNormal    = 0
	wbBreakAll  = 1
	wbBreakWord = 2
	wbKeepAll   = 3
)

// WrapText breaks text into lines according to the wrap configuration.
// WrapText wraps text into lines based on font metrics and constraints.
func WrapText(text string, cfg WrapConfig) []TextLine {
	if text == "" {
		return nil
	}

	var m TextMeasurer
	if cfg.Face != nil && cfg.Face.Variable && len(cfg.Variations) > 0 {
		if gtFace, err := cfg.Face.shaperFace(cfg.Variations); err == nil {
			m = NewVariationMeasurer(gtFace, cfg.FontSize, cfg.LetterSpacing)
		}
	}
	if m == nil {
		m = NewMeasurer(cfg.FontFace, cfg.LetterSpacing)
	}

	collapseWS := cfg.WhiteSpace == wsNormal || cfg.WhiteSpace == wsNowrap || cfg.WhiteSpace == wsPreLine
	preserveNL := cfg.WhiteSpace == wsPre || cfg.WhiteSpace == wsPreWrap || cfg.WhiteSpace == wsPreLine
	allowWrap := cfg.WhiteSpace != wsNowrap && cfg.WhiteSpace != wsPre

	var paragraphs []string
	if preserveNL {
		paragraphs = strings.Split(text, "\n")
		if collapseWS {
			for i, p := range paragraphs {
				paragraphs[i] = collapseWhitespace(p)
			}
		}
	} else {
		merged := strings.ReplaceAll(text, "\n", " ")
		if collapseWS {
			merged = collapseWhitespace(merged)
		}
		paragraphs = []string{merged}
	}

	var lines []TextLine
	lineIdx := 0

	for _, para := range paragraphs {
		pLines := wrapParagraph(m, para, cfg.MaxWidth, allowWrap, cfg.WordBreak, collapseWS)
		for _, l := range pLines {
			lines = append(lines, TextLine{
				Text:  l.Text,
				Width: l.Width,
				Y:     float64(lineIdx) * cfg.LineHeight,
			})
			lineIdx++
		}
	}

	if len(lines) == 0 && text != "" {
		lines = append(lines, TextLine{Text: "", Width: 0, Y: 0})
	}

	ellipsis := cfg.LineClamp > 0 || cfg.TextOverflow == "ellipsis"
	if ellipsis && len(lines) > 0 {
		clamp := cfg.LineClamp
		if clamp <= 0 {
			clamp = 1
		}
		if len(lines) > clamp {
			lines = lines[:clamp]
			lines[clamp-1] = truncateLineWithEllipsis(m, lines[clamp-1], cfg.MaxWidth)
		} else if cfg.TextOverflow == "ellipsis" && len(lines) == 1 && lines[0].Width > cfg.MaxWidth && cfg.MaxWidth > 0 {
			lines[0] = truncateLineWithEllipsis(m, lines[0], cfg.MaxWidth)
		}
	}

	return lines
}

func truncateLineWithEllipsis(m TextMeasurer, line TextLine, maxWidth float64) TextLine {
	const ellipsisStr = "..."
	ellipsisWidth := m.StringWidth(ellipsisStr)
	available := maxWidth - ellipsisWidth
	if available <= 0 {
		return TextLine{Text: ellipsisStr, Width: ellipsisWidth, X: line.X, Y: line.Y}
	}

	runes := []rune(line.Text)
	var width float64
	cutIdx := 0
	for i, r := range runes {
		rw := m.RuneWidth(r)
		spacing := float64(0)
		if i > 0 {
			spacing = m.LetterSpacing()
		}
		if width+spacing+rw > available {
			break
		}
		width += spacing + rw
		cutIdx = i + 1
	}

	truncated := string(runes[:cutIdx]) + ellipsisStr
	return TextLine{
		Text:  truncated,
		Width: width + ellipsisWidth,
		X:     line.X,
		Y:     line.Y,
	}
}

func collapseWhitespace(s string) string {
	var b strings.Builder
	inSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !inSpace {
				b.WriteRune(' ')
				inSpace = true
			}
		} else {
			b.WriteRune(r)
			inSpace = false
		}
	}
	return strings.TrimSpace(b.String())
}

func wrapParagraph(m TextMeasurer, text string, maxWidth float64, allowWrap bool, wordBreak int, collapsed bool) []TextLine {
	if !allowWrap || maxWidth <= 0 {
		w := m.StringWidth(text)
		return []TextLine{{Text: text, Width: w}}
	}

	if wordBreak == wbBreakAll {
		return wrapBreakAll(m, text, maxWidth)
	}

	if collapsed {
		return wrapCollapsed(m, text, maxWidth, wordBreak)
	}
	return wrapPreserved(m, text, maxWidth, wordBreak)
}

func wrapCollapsed(m TextMeasurer, text string, maxWidth float64, wordBreak int) []TextLine {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []TextLine{{Text: "", Width: 0}}
	}

	lines := make([]TextLine, 0, len(words)/8+1)
	var current strings.Builder
	current.Grow(len(text))
	var currentWidth float64
	spaceWidth := m.RuneWidth(' ')

	setCurrent := func(s string, w float64) {
		current.Reset()
		current.WriteString(s)
		currentWidth = w
	}

	for i, word := range words {
		wordWidth := m.StringWidth(word)

		if i == 0 {
			if wordBreak == wbBreakWord && wordWidth > maxWidth {
				broken := breakChars(m, word, maxWidth)
				for j, bl := range broken {
					if j < len(broken)-1 {
						lines = append(lines, bl)
					} else {
						setCurrent(bl.Text, bl.Width)
					}
				}
			} else {
				setCurrent(word, wordWidth)
			}
			continue
		}

		sep := spaceWidth
		newWidth := currentWidth + sep + wordWidth

		if newWidth <= maxWidth {
			current.WriteByte(' ')
			current.WriteString(word)
			currentWidth = newWidth
		} else {
			lines = append(lines, TextLine{Text: current.String(), Width: currentWidth})

			if wordBreak == wbBreakWord && wordWidth > maxWidth {
				broken := breakChars(m, word, maxWidth)
				for j, bl := range broken {
					if j < len(broken)-1 {
						lines = append(lines, bl)
					} else {
						setCurrent(bl.Text, bl.Width)
					}
				}
			} else {
				setCurrent(word, wordWidth)
			}
		}
	}

	lines = append(lines, TextLine{Text: current.String(), Width: currentWidth})
	return lines
}

func wrapPreserved(m TextMeasurer, text string, maxWidth float64, wordBreak int) []TextLine {
	type segment struct {
		text    string
		isSpace bool
	}

	var segments []segment
	runes := []rune(text)
	i := 0
	for i < len(runes) {
		if unicode.IsSpace(runes[i]) && runes[i] != '\n' {
			j := i
			for j < len(runes) && unicode.IsSpace(runes[j]) && runes[j] != '\n' {
				j++
			}
			segments = append(segments, segment{string(runes[i:j]), true})
			i = j
		} else {
			j := i
			for j < len(runes) && !unicode.IsSpace(runes[j]) {
				j++
			}
			segments = append(segments, segment{string(runes[i:j]), false})
			i = j
		}
	}

	if len(segments) == 0 {
		return []TextLine{{Text: "", Width: 0}}
	}

	var lines []TextLine
	var current strings.Builder
	var currentWidth float64

	for _, seg := range segments {
		segWidth := m.StringWidth(seg.text)

		if currentWidth+segWidth > maxWidth && current.Len() > 0 {
			t := current.String()
			lines = append(lines, TextLine{Text: t, Width: currentWidth})
			current.Reset()
			currentWidth = 0
		}

		if !seg.isSpace && wordBreak == wbBreakWord && segWidth > maxWidth && current.Len() == 0 {
			broken := breakChars(m, seg.text, maxWidth)
			for j, bl := range broken {
				if j < len(broken)-1 {
					lines = append(lines, bl)
				} else {
					current.WriteString(bl.Text)
					currentWidth = bl.Width
				}
			}
		} else {
			current.WriteString(seg.text)
			currentWidth += segWidth
		}
	}

	if current.Len() > 0 {
		t := current.String()
		lines = append(lines, TextLine{Text: t, Width: currentWidth})
	}

	if len(lines) == 0 {
		lines = append(lines, TextLine{Text: "", Width: 0})
	}

	return lines
}

func wrapBreakAll(m TextMeasurer, text string, maxWidth float64) []TextLine {
	var lines []TextLine
	var current strings.Builder
	var currentWidth float64

	for _, r := range text {
		rw := m.RuneWidth(r)
		spacing := float64(0)
		if current.Len() > 0 {
			spacing = m.LetterSpacing()
		}

		if currentWidth+spacing+rw > maxWidth && current.Len() > 0 {
			lines = append(lines, TextLine{Text: current.String(), Width: currentWidth})
			current.Reset()
			currentWidth = 0
			spacing = 0
		}

		current.WriteRune(r)
		currentWidth += spacing + rw
	}

	if current.Len() > 0 {
		lines = append(lines, TextLine{Text: current.String(), Width: currentWidth})
	}

	return lines
}

func breakChars(m TextMeasurer, word string, maxWidth float64) []TextLine {
	var lines []TextLine
	var current strings.Builder
	var currentWidth float64

	for _, r := range word {
		rw := m.RuneWidth(r)
		spacing := float64(0)
		if current.Len() > 0 {
			spacing = m.LetterSpacing()
		}

		if currentWidth+spacing+rw > maxWidth && current.Len() > 0 {
			lines = append(lines, TextLine{Text: current.String(), Width: currentWidth})
			current.Reset()
			currentWidth = 0
			spacing = 0
		}

		current.WriteRune(r)
		currentWidth += spacing + rw
	}

	if current.Len() > 0 {
		lines = append(lines, TextLine{Text: current.String(), Width: currentWidth})
	}

	return lines
}
