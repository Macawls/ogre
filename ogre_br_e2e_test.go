package ogre

import (
	"bytes"
	"image"
	"image/png"
	"strings"
	"testing"
)

// inkBandTops returns the first row of every horizontal band of ink, top to
// bottom. Used to measure line pitch without depending on glyph shapes.
func inkBandTops(img image.Image) []int {
	bounds := img.Bounds()
	isInk := make([]bool, bounds.Max.Y)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Composite over white: Ogre writes premultiplied RGBA, so a fully
			// transparent pixel reads as (0,0,0,0) and would count as ink.
			r, g, b, a := img.At(x, y).RGBA()
			over := func(c uint32) uint32 { return (c*255 + 255*(0xffff-a)) / 0xffff }
			if over(r) < 250 || over(g) < 250 || over(b) < 250 {
				isInk[y] = true
				break
			}
		}
	}

	var tops []int
	inBand := false
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		if isInk[y] && !inBand {
			tops = append(tops, y)
			inBand = true
		} else if !isInk[y] {
			inBand = false
		}
	}
	return tops
}

// A <br> must advance the baseline by exactly one line box. It used to advance
// by two: the break became its own text leaf, which took a full line box on top
// of the one the following text already occupied.
func TestRenderBrAdvancesOneLineBox(t *testing.T) {
	measure := func(html string) []int {
		result, err := Render(html, Options{Width: 600, Height: 300, Format: FormatPNG})
		if err != nil {
			t.Fatalf("Render: %v", err)
		}
		img, err := png.Decode(bytes.NewReader(result.Data))
		if err != nil {
			t.Fatalf("decode png: %v", err)
		}
		return inkBandTops(img)
	}

	style := `padding:20px;font-size:20px;line-height:1.15;color:#000;`

	// Control: the same three lines separated by block boxes, which already
	// advanced one line box each.
	blockTops := measure(`<body style="margin:0;padding:0;width:600px;height:300px;background:#fff;">` +
		`<div style="` + style + `"><span>Первая</span><span style="display:block">Вторая</span>` +
		`<span style="display:block">Третья</span></div></body>`)

	brTops := measure(`<body style="margin:0;padding:0;width:600px;height:300px;background:#fff;">` +
		`<div style="` + style + `">Первая<br>Вторая<br>Третья</div></body>`)

	if len(blockTops) != 3 {
		t.Fatalf("control produced %d ink bands, want 3: %v", len(blockTops), blockTops)
	}
	if len(brTops) != 3 {
		t.Fatalf("three lines separated by <br> produced %d ink bands, want 3: %v", len(brTops), brTops)
	}

	for i := 1; i < 3; i++ {
		gotStep := brTops[i] - brTops[i-1]
		wantStep := blockTops[i] - blockTops[i-1]
		if gotStep != wantStep {
			t.Fatalf("<br> line %d starts %d px after the previous line, want %d (one line box): %v",
				i+1, gotStep, wantStep, brTops)
		}
	}
}

// Consecutive <br> must each advance a line box, so A<br>B<br><br>C leaves a
// blank line before C. Reported by the maintainer on PR #37: the single <br>
// now advances one box, but a run of two still advances only one.
func TestRenderConsecutiveBrLeavesBlankLine(t *testing.T) {
	measure := func(html string) []int {
		result, err := Render(html, Options{Width: 600, Height: 400, Format: FormatPNG})
		if err != nil {
			t.Fatalf("Render: %v", err)
		}
		img, err := png.Decode(bytes.NewReader(result.Data))
		if err != nil {
			t.Fatalf("decode png: %v", err)
		}
		return inkBandTops(img)
	}

	head := `<body style="margin:0;padding:0;width:600px;height:400px;background:#fff;">`
	style := `padding:20px;font-size:20px;line-height:1.15;color:#000;`

	single := measure(head + `<div style="` + style + `">Первая<br>Вторая<br>Третья</div></body>`)
	double := measure(head + `<div style="` + style + `">Первая<br>Вторая<br><br>Третья</div></body>`)

	if len(single) != 3 || len(double) != 3 {
		t.Fatalf("want 3 ink bands each, got single=%v double=%v", single, double)
	}

	oneBox := single[1] - single[0]
	wantGap := 2 * oneBox
	gotGap := double[2] - double[1]

	if gotGap < wantGap-1 || gotGap > wantGap+1 {
		t.Fatalf("two consecutive <br> advance %d px, want %d (two line boxes of %d): single=%v double=%v",
			gotGap, wantGap, oneBox, single, double)
	}
}

// A trailing run of <br> must open one line box per break except the last:
// A<br> is still one line, A<br><br> is two. Chromium puts a following block
// at 23/46/69px for none/one/two trailing breaks, i.e. one box per break minus
// the last. Reported alongside the consecutive-break case on PR #37.
func TestRenderTrailingBrOpensNoExtraLine(t *testing.T) {
	render := func(inner string) []int {
		html := `<body style="margin:0;padding:0;width:600px;height:400px;background:#fff;">` +
			`<div style="font-size:20px;line-height:1.15;color:#000;">` + inner + `</div>` +
			`<div style="font-size:20px;line-height:1.15;color:#000;">Sibling</div></body>`
		result, err := Render(html, Options{Width: 600, Height: 400, Format: FormatPNG})
		if err != nil {
			t.Fatalf("Render: %v", err)
		}
		img, err := png.Decode(bytes.NewReader(result.Data))
		if err != nil {
			t.Fatalf("decode png: %v", err)
		}
		return inkBandTops(img)
	}

	// One line box, measured from two lines the engine already separates
	// correctly: the <br> case must not be the thing that defines it.
	// The sibling is a band of its own, so a two-line case yields three.
	two := render("Первая<br>Вторая")
	if len(two) != 3 {
		t.Fatalf("want 3 ink bands (two lines plus the sibling), got %v", two)
	}
	oneBox := two[1] - two[0]

	base := render("Первая")
	baseTop := base[len(base)-1]

	for _, tc := range []struct {
		breaks int
		want   int // extra line boxes beyond the single line
	}{
		{1, 0}, // a single trailing <br> adds nothing
		{2, 1},
		{3, 2},
	} {
		inner := "Первая" + strings.Repeat("<br>", tc.breaks)
		tops := render(inner)
		got := tops[len(tops)-1] - baseTop
		if got != tc.want*oneBox {
			t.Fatalf("%d trailing <br> push the sibling %d px down, want %d (%d line boxes of %d): %q",
				tc.breaks, got, tc.want*oneBox, tc.want, oneBox, inner)
		}
	}
}
