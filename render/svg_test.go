package render

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/macawls/ogre/layout"
	"github.com/macawls/ogre/parse"
	"github.com/macawls/ogre/style"
)

func buildSimpleTree(pn *parse.Node, styles map[*parse.Node]*style.ComputedStyle, w, h float64) *layout.LayoutTree {
	measure := func(_ *parse.Node, text string, cs *style.ComputedStyle, maxWidth float64) (float64, float64) {
		size := cs.FontSize
		if size == 0 {
			size = 16
		}
		return float64(len(text)) * size * 0.6, size * 1.2
	}
	return layout.ComputeLayout(pn, styles, w, h, measure)
}

func TestSimpleDivBackground(t *testing.T) {
	pn := &parse.Node{Type: parse.ElementNode, Tag: "div"}
	cs := style.NewComputedStyle()
	cs.BackgroundColor = style.Color{R: 255, G: 0, B: 0, A: 1}
	cs.Width = style.Value{Raw: 100, Unit: style.UnitPx}
	cs.Height = style.Value{Raw: 50, Unit: style.UnitPx}
	styles := map[*parse.Node]*style.ComputedStyle{pn: cs}
	tree := buildSimpleTree(pn, styles, 800, 600)

	svg := RenderSVG(tree, styles, nil, 800, 600)

	if !strings.Contains(svg, "<svg") {
		t.Fatal("missing <svg> root")
	}
	if !strings.Contains(svg, "</svg>") {
		t.Fatal("missing </svg> close")
	}
	if !strings.Contains(svg, `fill="#ff0000"`) {
		t.Fatalf("expected red fill, got: %s", svg)
	}
	if !strings.Contains(svg, "<rect") {
		t.Fatal("expected <rect> for background")
	}
}

func TestNestedDivs(t *testing.T) {
	child := &parse.Node{Type: parse.ElementNode, Tag: "div"}
	parent := &parse.Node{Type: parse.ElementNode, Tag: "div", Children: []*parse.Node{child}}

	csParent := style.NewComputedStyle()
	csParent.BackgroundColor = style.Color{R: 0, G: 0, B: 255, A: 1}
	csParent.Width = style.Value{Raw: 200, Unit: style.UnitPx}
	csParent.Height = style.Value{Raw: 200, Unit: style.UnitPx}
	csChild := style.NewComputedStyle()
	csChild.BackgroundColor = style.Color{R: 0, G: 255, B: 0, A: 1}
	csChild.Width = style.Value{Raw: 100, Unit: style.UnitPx}
	csChild.Height = style.Value{Raw: 100, Unit: style.UnitPx}
	styles := map[*parse.Node]*style.ComputedStyle{parent: csParent, child: csChild}
	tree := buildSimpleTree(parent, styles, 800, 600)

	svg := RenderSVG(tree, styles, nil, 800, 600)

	rects := strings.Count(svg, "<rect")
	if rects != 2 {
		t.Fatalf("expected 2 rects for nested divs, got %d: %s", rects, svg)
	}
	if !strings.Contains(svg, `fill="#0000ff"`) {
		t.Fatal("missing blue parent rect")
	}
	if !strings.Contains(svg, `fill="#00ff00"`) {
		t.Fatal("missing green child rect")
	}
}

func TestTextNode(t *testing.T) {
	textNode := &parse.Node{Type: parse.TextNode, Text: "Hello"}
	parent := &parse.Node{Type: parse.ElementNode, Tag: "div", Children: []*parse.Node{textNode}}

	csParent := &style.ComputedStyle{
		Width:  style.Value{Raw: 200, Unit: style.UnitPx},
		Height: style.Value{Raw: 50, Unit: style.UnitPx},
	}
	csText := &style.ComputedStyle{
		FontFamily: "Arial",
		FontSize:   14,
		FontWeight: 700,
		Color:      style.Color{R: 0, G: 0, B: 0, A: 1},
	}
	styles := map[*parse.Node]*style.ComputedStyle{parent: csParent, textNode: csText}
	tree := buildSimpleTree(parent, styles, 800, 600)

	svg := RenderSVG(tree, styles, nil, 800, 600)

	if !strings.Contains(svg, "<text") {
		t.Fatalf("expected <text> element, got: %s", svg)
	}
	if !strings.Contains(svg, "Hello") {
		t.Fatal("expected text content")
	}
	if !strings.Contains(svg, `font-family="Arial"`) {
		t.Fatalf("expected font-family, got: %s", svg)
	}
	if !strings.Contains(svg, `font-weight="700"`) {
		t.Fatalf("expected font-weight 700, got: %s", svg)
	}
}

func TestWellFormedXML(t *testing.T) {
	textNode := &parse.Node{Type: parse.TextNode, Text: "A & B <C>"}
	parent := &parse.Node{Type: parse.ElementNode, Tag: "div", Children: []*parse.Node{textNode}}

	csParent := &style.ComputedStyle{
		BackgroundColor: style.Color{R: 200, G: 200, B: 200, A: 1},
		Width:           style.Value{Raw: 300, Unit: style.UnitPx},
		Height:          style.Value{Raw: 100, Unit: style.UnitPx},
	}
	csText := &style.ComputedStyle{
		FontSize: 16,
		Color:    style.Color{R: 0, G: 0, B: 0, A: 1},
	}
	styles := map[*parse.Node]*style.ComputedStyle{parent: csParent, textNode: csText}
	tree := buildSimpleTree(parent, styles, 800, 600)

	svg := RenderSVG(tree, styles, nil, 800, 600)

	decoder := xml.NewDecoder(strings.NewReader(svg))
	for {
		_, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			t.Fatalf("SVG is not well-formed XML: %v\nSVG: %s", err, svg)
		}
	}
}

func TestBorderRadius(t *testing.T) {
	pn := &parse.Node{Type: parse.ElementNode, Tag: "div"}
	cs := &style.ComputedStyle{
		BackgroundColor:         style.Color{R: 100, G: 100, B: 100, A: 1},
		Width:                   style.Value{Raw: 100, Unit: style.UnitPx},
		Height:                  style.Value{Raw: 100, Unit: style.UnitPx},
		BorderTopLeftRadius:     8,
		BorderTopRightRadius:    8,
		BorderBottomLeftRadius:  8,
		BorderBottomRightRadius: 8,
	}
	styles := map[*parse.Node]*style.ComputedStyle{pn: cs}
	tree := buildSimpleTree(pn, styles, 800, 600)

	svg := RenderSVG(tree, styles, nil, 800, 600)

	if !strings.Contains(svg, `rx="8"`) {
		t.Fatalf("expected rx attribute, got: %s", svg)
	}
	if !strings.Contains(svg, `ry="8"`) {
		t.Fatalf("expected ry attribute, got: %s", svg)
	}
}

func TestOpacityGroup(t *testing.T) {
	pn := &parse.Node{Type: parse.ElementNode, Tag: "div"}
	cs := &style.ComputedStyle{
		BackgroundColor: style.Color{R: 255, G: 0, B: 0, A: 1},
		Width:           style.Value{Raw: 100, Unit: style.UnitPx},
		Height:          style.Value{Raw: 50, Unit: style.UnitPx},
		Opacity:         0.5,
	}
	styles := map[*parse.Node]*style.ComputedStyle{pn: cs}
	tree := buildSimpleTree(pn, styles, 800, 600)

	svg := RenderSVG(tree, styles, nil, 800, 600)

	if !strings.Contains(svg, `<g opacity="0.5"`) {
		t.Fatalf("expected opacity group, got: %s", svg)
	}
	if !strings.Contains(svg, "</g>") {
		t.Fatal("expected closing </g>")
	}
}

func TestTransparentBackgroundSkipped(t *testing.T) {
	pn := &parse.Node{Type: parse.ElementNode, Tag: "div"}
	cs := &style.ComputedStyle{
		BackgroundColor: style.Color{R: 0, G: 0, B: 0, A: 0},
		Width:           style.Value{Raw: 100, Unit: style.UnitPx},
		Height:          style.Value{Raw: 50, Unit: style.UnitPx},
	}
	styles := map[*parse.Node]*style.ComputedStyle{pn: cs}
	tree := buildSimpleTree(pn, styles, 800, 600)

	svg := RenderSVG(tree, styles, nil, 800, 600)

	if strings.Contains(svg, "<rect") {
		t.Fatalf("transparent background should not produce a rect, got: %s", svg)
	}
}
