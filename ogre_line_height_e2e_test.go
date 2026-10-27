package ogre

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

// This file exercises line-height through the public Render entry point, so the
// check runs on the real pipeline rather than on the style package alone.

func renderLineHeightSVG(t *testing.T, html string, w, h int) string {
	t.Helper()
	res, err := Render(html, Options{Width: w, Height: h})
	if err != nil {
		t.Fatalf("render: %v", err)
	}
	return string(res.Data)
}

// lineHeightBaselines returns the y of every text run, in document order.
func lineHeightBaselines(t *testing.T, svg string) []float64 {
	t.Helper()

	var ys []float64
	rest := svg
	for {
		i := strings.Index(rest, "translate(")
		if i < 0 {
			break
		}
		rest = rest[i+len("translate("):]
		end := strings.Index(rest, ")")
		if end < 0 {
			break
		}
		parts := strings.Split(rest[:end], ",")
		if len(parts) == 2 {
			if y, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err == nil {
				ys = append(ys, y)
			}
		}
		rest = rest[end:]
	}
	if len(ys) == 0 {
		t.Fatal("no text runs found in SVG")
	}
	return ys
}

func TestRenderLineHeightAffectsLineStep(t *testing.T) {
	// Same text, same width: only line-height differs. The vertical step
	// between wrapped lines must follow the declared value, not a fixed 1.2.
	wrap := "alpha beta gamma delta epsilon zeta eta theta"

	render := func(lh string) []float64 {
		html := `<div style="font-family:sans-serif;font-size:20px;width:160px;line-height:` + lh + `">` + wrap + `</div>`
		return lineHeightBaselines(t, renderLineHeightSVG(t, html, 400, 400))
	}

	tight := render("1.0")
	wide := render("3.0")

	if len(tight) < 2 || len(wide) < 2 {
		t.Fatalf("expected wrapped text, got %d and %d lines", len(tight), len(wide))
	}

	stepTight := tight[1] - tight[0]
	stepWide := wide[1] - wide[0]

	if math.Abs(stepTight-20) > 0.5 {
		t.Errorf("line-height 1.0 step: got %.2f, want 20", stepTight)
	}
	if math.Abs(stepWide-60) > 0.5 {
		t.Errorf("line-height 3.0 step: got %.2f, want 60", stepWide)
	}
	if math.Abs(stepWide-stepTight) < 30 {
		t.Errorf("line-height ignored: steps are %.2f and %.2f", stepTight, stepWide)
	}
}

func TestRenderLineHeightInheritsFromParent(t *testing.T) {
	// The child declares no line-height, so it inherits the parent's 2.5 and
	// must not fall back to the 1.2 default.
	html := `<div style="font-family:sans-serif;font-size:20px;width:160px;line-height:2.5">` +
		`<div>alpha beta gamma delta epsilon zeta eta theta</div></div>`

	ys := lineHeightBaselines(t, renderLineHeightSVG(t, html, 400, 400))
	if len(ys) < 2 {
		t.Fatalf("expected wrapped text, got %d lines", len(ys))
	}
	step := ys[1] - ys[0]
	if math.Abs(step-50) > 0.5 {
		t.Errorf("inherited line-height step: got %.2f, want 50 (2.5 x 20)", step)
	}
}

func TestRenderLineHeightPercentUsesOwnFontSize(t *testing.T) {
	// A percentage line-height resolves against the element's own font-size,
	// not against the container height.
	html := `<div style="font-family:sans-serif;font-size:20px;width:160px;line-height:250%">` +
		`alpha beta gamma delta epsilon zeta eta theta</div>`

	ys := lineHeightBaselines(t, renderLineHeightSVG(t, html, 400, 400))
	if len(ys) < 2 {
		t.Fatalf("expected wrapped text, got %d lines", len(ys))
	}
	step := ys[1] - ys[0]
	if math.Abs(step-50) > 0.5 {
		t.Errorf("line-height 250%% step: got %.2f, want 50 (2.5 x 20)", step)
	}
}
