package ogre

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

// This file exercises box-sizing through the public Render entry point. Boxes
// are emitted as <rect> elements, so the declared width is readable straight
// from the SVG without parsing glyph geometry.

func renderBoxSizingSVG(t *testing.T, html string, w, h int) string {
	t.Helper()
	res, err := Render(html, Options{Width: w, Height: h})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	return string(res.Data)
}

// boxSizingRects returns the width of every visible <rect> in document order.
func boxSizingRects(t *testing.T, svg string) []float64 {
	t.Helper()

	widths := make([]float64, 0, 4)
	for _, chunk := range strings.Split(svg, "<rect")[1:] {
		start := strings.Index(chunk, `width="`)
		if start < 0 {
			continue
		}
		raw := chunk[start+len(`width="`):]
		if end := strings.Index(raw, `"`); end >= 0 {
			raw = raw[:end]
		}
		v, err := strconv.ParseFloat(raw, 64)
		if err != nil || v <= 1 {
			continue
		}
		widths = append(widths, v)
	}
	if len(widths) < 2 {
		t.Fatalf("expected at least two boxes, got %d", len(widths))
	}
	return widths
}

func TestRenderBoxSizingContentBox(t *testing.T) {
	// border-box (the engine default) keeps the declared width; content-box has
	// to add the horizontal padding on top of it.
	html := `<div style="width:400px">` +
		`<div style="width:200px;padding:0 50px;height:40px;background:#c00"></div>` +
		`<div style="box-sizing:content-box;width:200px;padding:0 50px;height:40px;background:#00c"></div>` +
		`</div>`

	widths := boxSizingRects(t, renderBoxSizingSVG(t, html, 1280, 720))

	if math.Abs(widths[0]-200) > 0.01 {
		t.Errorf("border-box width: got %.2f, want 200", widths[0])
	}
	if math.Abs(widths[1]-300) > 0.01 {
		t.Errorf("content-box width: got %.2f, want 300 (200 + 2 x 50 padding)", widths[1])
	}
}

func TestRenderBoxSizingContentBoxVertical(t *testing.T) {
	// The same conversion must apply on the vertical axis.
	html := `<div style="height:400px">` +
		`<div style="height:100px;padding:25px 0;width:40px;background:#0a0"></div>` +
		`<div style="box-sizing:content-box;height:100px;padding:25px 0;width:40px;background:#a0a"></div>` +
		`</div>`

	svg := renderBoxSizingSVG(t, html, 400, 400)

	var heights []float64
	for _, chunk := range strings.Split(svg, "<rect")[1:] {
		start := strings.Index(chunk, `height="`)
		if start < 0 {
			continue
		}
		raw := chunk[start+len(`height="`):]
		if end := strings.Index(raw, `"`); end >= 0 {
			raw = raw[:end]
		}
		if v, err := strconv.ParseFloat(raw, 64); err == nil && v > 1 {
			heights = append(heights, v)
		}
	}
	if len(heights) < 2 {
		t.Fatalf("expected two boxes, got %d", len(heights))
	}
	if math.Abs(heights[0]-100) > 0.01 {
		t.Errorf("border-box height: got %.2f, want 100", heights[0])
	}
	if math.Abs(heights[1]-150) > 0.01 {
		t.Errorf("content-box height: got %.2f, want 150 (100 + 2 x 25 padding)", heights[1])
	}
}

// boxSizingBoxes returns (x, width) of every visible <rect> in document order.
func boxSizingBoxes(t *testing.T, svg string) [][2]float64 {
	t.Helper()
	num := func(chunk, attr string) (float64, bool) {
		start := strings.Index(chunk, attr+`="`)
		if start < 0 {
			return 0, false
		}
		raw := chunk[start+len(attr)+2:]
		if end := strings.Index(raw, `"`); end >= 0 {
			raw = raw[:end]
		}
		v, err := strconv.ParseFloat(raw, 64)
		return v, err == nil
	}
	boxes := make([][2]float64, 0, 4)
	for _, chunk := range strings.Split(svg, "<rect")[1:] {
		w, ok := num(chunk, "width")
		if !ok || w <= 1 {
			continue
		}
		x, _ := num(chunk, "x")
		boxes = append(boxes, [2]float64{x, w})
	}
	if len(boxes) < 2 {
		t.Fatalf("expected at least two boxes, got %d", len(boxes))
	}
	return boxes
}

// content-box must be honoured by the flex sizing pass too, not only by the
// element's own layout: the flex base size is read straight off the declared
// width, so a padded content-box item overflows and its sibling overlaps it.
// Reported by the maintainer on PR #26.
func TestRenderBoxSizingContentBoxFlexSibling(t *testing.T) {
	html := `<div style="display:flex;width:600px;height:40px">` +
		`<div style="box-sizing:content-box;width:200px;padding:0 50px;height:40px;background:#c00"></div>` +
		`<div style="width:100px;height:40px;background:#00c"></div>` +
		`</div>`

	boxes := boxSizingBoxes(t, renderBoxSizingSVG(t, html, 600, 100))

	if math.Abs(boxes[0][1]-300) > 0.01 {
		t.Errorf("content-box item width: got %.2f, want 300 (200 + 2 x 50 padding)", boxes[0][1])
	}
	if math.Abs(boxes[1][0]-300) > 0.01 {
		t.Errorf("sibling starts at x=%.2f, want 300: it overlaps the padded content-box item", boxes[1][0])
	}
}
