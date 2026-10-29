package render

import (
	"bytes"
	"image/png"
	"strings"
	"testing"

	"github.com/macawls/ogre/v3/style"
)

// CSS gradients interpolate in gamma-encoded sRGB, so the SVG defs must ask for
// sRGB as well. The SVG 1.1 default is linearRGB, which is what Ogre emitted
// before #27 and what made every gradient render too bright.
func TestSVGLinearGradientUsesSRGB(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: "linear-gradient(90deg, red, blue)",
	}
	r := RenderBackground(cs, 0, 0, 100, 100, testIDGen)

	if !strings.Contains(r.Defs, `color-interpolation="sRGB"`) {
		t.Errorf("linear gradient should request sRGB interpolation, got %q", r.Defs)
	}
	if strings.Contains(r.Defs, "linearRGB") {
		t.Errorf("linear gradient still requests linearRGB, got %q", r.Defs)
	}
}

func TestSVGRadialGradientUsesSRGB(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: "radial-gradient(circle, red, blue)",
	}
	r := RenderBackground(cs, 0, 0, 100, 100, testIDGen)

	if !strings.Contains(r.Defs, `color-interpolation="sRGB"`) {
		t.Errorf("radial gradient should request sRGB interpolation, got %q", r.Defs)
	}
	if strings.Contains(r.Defs, "linearRGB") {
		t.Errorf("radial gradient still requests linearRGB, got %q", r.Defs)
	}
}

// A red-to-blue gradient at its midpoint must land near (128,0,128). Linear
// interpolation would put it near (188,0,188), which is the classic "washed
// out" look of a gradient blended in the wrong colour space.
func TestPNGGradientMidpointIsSRGB(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: "linear-gradient(90deg, red, blue)",
	}
	tree, styles := buildTestTree(cs, 100, 100)

	data, err := RenderPNG(tree, styles, nil, 100, 100)
	if err != nil {
		t.Fatalf("RenderPNG returned error: %v", err)
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	r, _, b, a := img.At(50, 50).RGBA()
	gotR, gotB, gotA := r>>8, b>>8, a>>8

	if gotA < 250 {
		t.Fatalf("midpoint should be opaque, got alpha %d", gotA)
	}
	// sRGB blend: 127.5 in each channel. Allow for dithering and rounding.
	if gotR < 120 || gotR > 136 {
		t.Errorf("midpoint red = %d, want ~128 (sRGB blend; linear blend would be ~188)", gotR)
	}
	if gotB < 120 || gotB > 136 {
		t.Errorf("midpoint blue = %d, want ~128 (sRGB blend; linear blend would be ~188)", gotB)
	}
}

// image.RGBA holds premultiplied pixels. A gradient that fades to a fully
// transparent stop must therefore write RGB scaled by alpha, otherwise the
// stored pixel is invalid and un-premultiplying yields a dark fringe.
func TestPNGGradientAlphaIsPremultiplied(t *testing.T) {
	cs := &style.ComputedStyle{
		BackgroundImage: "linear-gradient(90deg, rgba(255,0,0,1), rgba(0,0,255,0))",
	}
	tree, styles := buildTestTree(cs, 400, 100)

	data, err := RenderPNG(tree, styles, nil, 400, 100)
	if err != nil {
		t.Fatalf("RenderPNG returned error: %v", err)
	}
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	// Every pixel must be a valid premultiplied pixel: no channel above alpha.
	for x := 0; x < 400; x += 7 {
		r, g, b, a := img.At(x, 50).RGBA()
		rr, gg, bb, aa := r>>8, g>>8, b>>8, a>>8
		if rr > aa || gg > aa || bb > aa {
			t.Fatalf("pixel x=%d = (%d,%d,%d,a=%d) is not premultiplied", x, rr, gg, bb, aa)
		}
	}

	// At the transparent end the visible colour must stay near-white when
	// composited on white, not fade to black.
	r, g, b, a := img.At(395, 50).RGBA()
	rr, gg, bb, aa := r>>8, g>>8, b>>8, a>>8
	if aa > 8 {
		t.Fatalf("x=395 should be almost transparent, got alpha %d", aa)
	}
	compositedG := int(gg) + int(255-aa)*(255-int(gg))/255
	if compositedG < 235 {
		t.Errorf("composited green at x=395 = %d, want >=235 (fade must not darken); raw pixel (%d,%d,%d,a=%d)",
			compositedG, rr, gg, bb, aa)
	}
}

// The interpolation itself must premultiply, so that a fade between two
// different hues keeps the visible side's hue instead of swinging through a
// desaturated middle.
func TestGradientInterpolationIsPremultiplied(t *testing.T) {
	stops := []colorStopF{
		{R: 1, G: 0, B: 0, A: 1, Position: 0},
		{R: 0, G: 0, B: 1, A: 0, Position: 1},
	}
	r, _, b, a := interpolateStops(stops, 0.5)

	if a < 0.45 || a > 0.55 {
		t.Errorf("alpha at midpoint = %v, want ~0.5", a)
	}
	// Premultiplied blend keeps red dominant: r/a ~= 1, b/a ~= 0.
	if a > 0 {
		straightR, straightB := r/a, b/a
		if straightB > 0.1 {
			t.Errorf("blue at midpoint = %v of alpha, want ~0 (hue must not swing)", straightB)
		}
		if straightR < 0.9 {
			t.Errorf("red at midpoint = %v of alpha, want ~1", straightR)
		}
	}
}
