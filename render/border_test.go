package render

import (
	"strings"
	"testing"

	"github.com/macawls/ogre/style"
)

func TestUniformBorder(t *testing.T) {
	cs := &style.ComputedStyle{
		BorderTopWidth:    2,
		BorderRightWidth:  2,
		BorderBottomWidth: 2,
		BorderLeftWidth:   2,
		BorderTopStyle:    style.BorderStyleSolid,
		BorderRightStyle:  style.BorderStyleSolid,
		BorderBottomStyle: style.BorderStyleSolid,
		BorderLeftStyle:   style.BorderStyleSolid,
		BorderTopColor:    style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderRightColor:  style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderBottomColor: style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderLeftColor:   style.Color{R: 0, G: 0, B: 0, A: 1},
	}

	result := RenderBorders(cs, 10, 20, 100, 50)

	if !strings.Contains(result, "<rect") {
		t.Fatalf("expected <rect>, got %q", result)
	}
	if !strings.Contains(result, `stroke="#000000"`) {
		t.Errorf("expected stroke color, got %q", result)
	}
	if !strings.Contains(result, `stroke-width="2"`) {
		t.Errorf("expected stroke-width, got %q", result)
	}
	if !strings.Contains(result, `fill="none"`) {
		t.Errorf("expected fill=none, got %q", result)
	}
}

func TestMixedBorderWidths(t *testing.T) {
	cs := &style.ComputedStyle{
		BorderTopWidth:    2,
		BorderRightWidth:  4,
		BorderBottomWidth: 2,
		BorderLeftWidth:   4,
		BorderTopStyle:    style.BorderStyleSolid,
		BorderRightStyle:  style.BorderStyleSolid,
		BorderBottomStyle: style.BorderStyleSolid,
		BorderLeftStyle:   style.BorderStyleSolid,
		BorderTopColor:    style.Color{R: 255, G: 0, B: 0, A: 1},
		BorderRightColor:  style.Color{R: 0, G: 255, B: 0, A: 1},
		BorderBottomColor: style.Color{R: 0, G: 0, B: 255, A: 1},
		BorderLeftColor:   style.Color{R: 255, G: 255, B: 0, A: 1},
	}

	result := RenderBorders(cs, 0, 0, 100, 100)

	if !strings.Contains(result, "<line") {
		t.Fatalf("expected <line> elements for mixed borders, got %q", result)
	}
	if strings.Count(result, "<line") != 4 {
		t.Errorf("expected 4 <line> elements, got %d in %q", strings.Count(result, "<line"), result)
	}
	if !strings.Contains(result, `stroke="#ff0000"`) {
		t.Errorf("expected red top border, got %q", result)
	}
	if !strings.Contains(result, `stroke="#00ff00"`) {
		t.Errorf("expected green right border, got %q", result)
	}
}

func TestBorderWithRadius(t *testing.T) {
	cs := &style.ComputedStyle{
		BorderTopWidth:          2,
		BorderRightWidth:        2,
		BorderBottomWidth:       2,
		BorderLeftWidth:         2,
		BorderTopStyle:          style.BorderStyleSolid,
		BorderRightStyle:        style.BorderStyleSolid,
		BorderBottomStyle:       style.BorderStyleSolid,
		BorderLeftStyle:         style.BorderStyleSolid,
		BorderTopColor:          style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderRightColor:       style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderBottomColor:      style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderLeftColor:        style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderTopLeftRadius:     8,
		BorderTopRightRadius:    8,
		BorderBottomRightRadius: 8,
		BorderBottomLeftRadius:  8,
	}

	result := RenderBorders(cs, 0, 0, 100, 50)

	if !strings.Contains(result, "<rect") {
		t.Fatalf("expected <rect> for uniform border with radius, got %q", result)
	}
	if !strings.Contains(result, `rx="8"`) {
		t.Errorf("expected rx=8, got %q", result)
	}
}

func TestDashedBorder(t *testing.T) {
	cs := &style.ComputedStyle{
		BorderTopWidth:    3,
		BorderRightWidth:  3,
		BorderBottomWidth: 3,
		BorderLeftWidth:   3,
		BorderTopStyle:    style.BorderStyleDashed,
		BorderRightStyle:  style.BorderStyleDashed,
		BorderBottomStyle: style.BorderStyleDashed,
		BorderLeftStyle:   style.BorderStyleDashed,
		BorderTopColor:    style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderRightColor:  style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderBottomColor: style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderLeftColor:   style.Color{R: 0, G: 0, B: 0, A: 1},
	}

	result := RenderBorders(cs, 0, 0, 100, 50)

	if !strings.Contains(result, "stroke-dasharray") {
		t.Errorf("expected stroke-dasharray for dashed border, got %q", result)
	}
	if !strings.Contains(result, `stroke-dasharray="6 3"`) {
		t.Errorf("expected dasharray 6 3, got %q", result)
	}
}

func TestDottedBorder(t *testing.T) {
	cs := &style.ComputedStyle{
		BorderTopWidth:    2,
		BorderRightWidth:  2,
		BorderBottomWidth: 2,
		BorderLeftWidth:   2,
		BorderTopStyle:    style.BorderStyleDotted,
		BorderRightStyle:  style.BorderStyleDotted,
		BorderBottomStyle: style.BorderStyleDotted,
		BorderLeftStyle:   style.BorderStyleDotted,
		BorderTopColor:    style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderRightColor:  style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderBottomColor: style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderLeftColor:   style.Color{R: 0, G: 0, B: 0, A: 1},
	}

	result := RenderBorders(cs, 0, 0, 100, 50)

	if !strings.Contains(result, `stroke-dasharray="2 2"`) {
		t.Errorf("expected dotted dasharray 2 2, got %q", result)
	}
}

func TestNoBorder(t *testing.T) {
	cs := &style.ComputedStyle{
		BorderTopStyle:    style.BorderStyleNone,
		BorderRightStyle:  style.BorderStyleNone,
		BorderBottomStyle: style.BorderStyleNone,
		BorderLeftStyle:   style.BorderStyleNone,
	}

	result := RenderBorders(cs, 0, 0, 100, 50)

	if result != "" {
		t.Errorf("expected empty string for no border, got %q", result)
	}
}

func TestNoBorderZeroWidth(t *testing.T) {
	cs := &style.ComputedStyle{
		BorderTopWidth:    0,
		BorderRightWidth:  0,
		BorderBottomWidth: 0,
		BorderLeftWidth:   0,
		BorderTopStyle:    style.BorderStyleSolid,
		BorderRightStyle:  style.BorderStyleSolid,
		BorderBottomStyle: style.BorderStyleSolid,
		BorderLeftStyle:   style.BorderStyleSolid,
	}

	result := RenderBorders(cs, 0, 0, 100, 50)

	if result != "" {
		t.Errorf("expected empty string for zero-width border, got %q", result)
	}
}

func TestMixedBorderWithRadius(t *testing.T) {
	cs := &style.ComputedStyle{
		BorderTopWidth:          2,
		BorderRightWidth:        4,
		BorderBottomWidth:       2,
		BorderLeftWidth:         4,
		BorderTopStyle:          style.BorderStyleSolid,
		BorderRightStyle:        style.BorderStyleSolid,
		BorderBottomStyle:       style.BorderStyleSolid,
		BorderLeftStyle:         style.BorderStyleSolid,
		BorderTopColor:          style.Color{R: 255, G: 0, B: 0, A: 1},
		BorderRightColor:       style.Color{R: 0, G: 0, B: 255, A: 1},
		BorderBottomColor:      style.Color{R: 255, G: 0, B: 0, A: 1},
		BorderLeftColor:        style.Color{R: 0, G: 0, B: 255, A: 1},
		BorderTopLeftRadius:     10,
		BorderTopRightRadius:    10,
		BorderBottomRightRadius: 10,
		BorderBottomLeftRadius:  10,
	}

	result := RenderBorders(cs, 0, 0, 200, 100)

	if !strings.Contains(result, "<path") {
		t.Fatalf("expected <path> for mixed borders with radius, got %q", result)
	}
	if !strings.Contains(result, "A") {
		t.Errorf("expected arc commands in path, got %q", result)
	}
}

func TestDoubleBorderUniform(t *testing.T) {
	cs := &style.ComputedStyle{
		BorderTopWidth:    6,
		BorderRightWidth:  6,
		BorderBottomWidth: 6,
		BorderLeftWidth:   6,
		BorderTopStyle:    style.BorderStyleDouble,
		BorderRightStyle:  style.BorderStyleDouble,
		BorderBottomStyle: style.BorderStyleDouble,
		BorderLeftStyle:   style.BorderStyleDouble,
		BorderTopColor:    style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderRightColor:  style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderBottomColor: style.Color{R: 0, G: 0, B: 0, A: 1},
		BorderLeftColor:   style.Color{R: 0, G: 0, B: 0, A: 1},
	}

	result := RenderBorders(cs, 0, 0, 100, 50)

	if strings.Count(result, "<rect") != 2 {
		t.Errorf("expected 2 <rect> elements for double border, got %d in %q", strings.Count(result, "<rect"), result)
	}
}
