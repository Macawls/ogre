package ogre

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// This file exercises letter-spacing through the public Render entry point.
// Text is emitted as one path per run, and every glyph outline inside that path
// starts with an absolute M, so the M coordinates are the glyph origins.

var letterSpacingOriginRe = regexp.MustCompile(`M([-+0-9.eE]+) ([-+0-9.eE]+)`)

func renderLetterSpacingSVG(t *testing.T, html string, w, h int) string {
	t.Helper()
	res, err := Render(html, Options{Width: w, Height: h})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	return string(res.Data)
}

// letterSpacingXRange returns the horizontal extent of every glyph origin in
// the SVG, in the run's local coordinate space.
func letterSpacingXRange(t *testing.T, svg string) (minX, maxX float64) {
	t.Helper()

	minX, maxX = math.Inf(1), math.Inf(-1)
	seen := false

	for _, chunk := range strings.Split(svg, `d="`)[1:] {
		d := chunk
		if end := strings.Index(d, `"`); end >= 0 {
			d = d[:end]
		}
		for _, m := range letterSpacingOriginRe.FindAllStringSubmatch(d, -1) {
			x, err := strconv.ParseFloat(m[1], 64)
			if err != nil {
				continue
			}
			seen = true
			minX = math.Min(minX, x)
			maxX = math.Max(maxX, x)
		}
	}

	if !seen {
		t.Fatal("no glyph origins found in SVG")
	}
	return minX, maxX
}

func TestRenderLetterSpacingWidensText(t *testing.T) {
	// letter-spacing is measured correctly but was never passed to the glyph
	// path builder, so the rendered text kept its unspaced width while layout
	// reserved the wider box.
	text := "ABCDEFGHIJ"
	render := func(ls string) (float64, float64) {
		html := `<div style="font-family:sans-serif;font-size:20px;letter-spacing:` + ls + `">` + text + `</div>`
		return letterSpacingXRange(t, renderLetterSpacingSVG(t, html, 800, 200))
	}

	plainMin, plainMax := render("0px")
	spacedMin, spacedMax := render("10px")

	if math.Abs(plainMin-spacedMin) > 0.01 {
		t.Errorf("letter-spacing moved the start of the text: %.2f vs %.2f", plainMin, spacedMin)
	}
	grew := (spacedMax - spacedMin) - (plainMax - plainMin)
	want := 10.0 * float64(len(text)-1)
	if math.Abs(grew-want) > 0.5 {
		t.Errorf("letter-spacing 10px widened the text by %.2f, want %.2f", grew, want)
	}
}

func TestRenderLetterSpacingEmUnits(t *testing.T) {
	// 0.5em at font-size 20px is 10px, so it must widen the text exactly like
	// the px declaration does.
	text := "ABCDEFGHIJ"
	render := func(ls string) float64 {
		html := `<div style="font-family:sans-serif;font-size:20px;letter-spacing:` + ls + `">` + text + `</div>`
		minX, maxX := letterSpacingXRange(t, renderLetterSpacingSVG(t, html, 800, 200))
		return maxX - minX
	}

	em := render("0.5em")
	px := render("10px")
	if math.Abs(em-px) > 0.5 {
		t.Errorf("letter-spacing 0.5em gave %.2f, 10px gave %.2f: em must resolve to px", em, px)
	}
}
