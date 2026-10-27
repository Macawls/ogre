package render

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"testing"

	"github.com/macawls/ogre/v3/style"
)

// renderShadowPNG renders a white box of size boxW x boxH at (boxX, boxY) on a
// white canvas carrying the given box-shadow, and returns the result composited
// over white, which is how a browser presents it.
func renderShadowPNG(t *testing.T, shadow string, w, h, boxW, boxH, boxX, boxY int) *image.RGBA {
	t.Helper()
	parentCS := &style.ComputedStyle{
		BackgroundColor: style.Color{R: 255, G: 255, B: 255, A: 1},
	}
	childCS := &style.ComputedStyle{
		BackgroundColor: style.Color{R: 255, G: 255, B: 255, A: 1},
		BoxShadow:       shadow,
	}
	tree, styles := buildTestTreeWithChild(parentCS, childCS,
		float64(w), float64(h), float64(boxW), float64(boxH), float64(boxX), float64(boxY))
	data, err := RenderPNG(tree, styles, nil, w, h)
	if err != nil {
		t.Fatalf("RenderPNG(%q) returned error: %v", shadow, err)
	}
	src, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}
	out := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(out, out.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	draw.Draw(out, out.Bounds(), src, image.Point{}, draw.Over)
	return out
}

// image.RGBA stores premultiplied samples, so blendOver must return premultiplied
// channels too. It used to divide by the output alpha, which wrote straight-alpha
// values into a premultiplied buffer: a 0.8-alpha red shadow came out at 40 %
// intensity ([101,51,51] instead of [255,51,51] on white), and every inset shadow
// was off by the same class of error.
func TestBlendOverReturnsPremultiplied(t *testing.T) {
	// Over a fully transparent destination the result is the source scaled by its
	// own alpha: 255 * 204/255 = 204, never 255.
	r, g, b, a := blendOver(255, 0, 0, 204, 0, 0, 0, 0)
	if r != 204 || g != 0 || b != 0 || a != 204 {
		t.Fatalf("blend over transparent = (%d,%d,%d,a=%d), want (204,0,0,a=204)", r, g, b, a)
	}
	// Over opaque white the composite is straight-alpha red at 0.8: (255,51,51).
	r, g, b, a = blendOver(255, 0, 0, 204, 255, 255, 255, 255)
	if r != 255 || g != 51 || b != 51 || a != 255 {
		t.Fatalf("blend over white = (%d,%d,%d,a=%d), want (255,51,51,a=255)", r, g, b, a)
	}
}

func TestBlendOverOpaqueSourceIsPassthrough(t *testing.T) {
	r, g, b, a := blendOver(10, 20, 30, 255, 200, 200, 200, 255)
	if r != 10 || g != 20 || b != 30 || a != 255 {
		t.Fatalf("opaque source = (%d,%d,%d,a=%d), want (10,20,30,a=255)", r, g, b, a)
	}
}

// A CSS blur radius b is a Gaussian with sigma = b/2, and four box passes of
// radius r have variance 4(r^2+r)/3. The calibration factor is what makes the
// rendered falloff match Chromium's over blur 0..48 px.
func TestShadowBlurRadiusCalibration(t *testing.T) {
	if got := shadowBlurRadius(0); got != 0 {
		t.Errorf("blur 0 should need no passes, got radius %d", got)
	}
	for _, blur := range []float64{2, 8, 24, 48} {
		want := int(math.Round((math.Sqrt(1+3*blur*blur*shadowBlurScale*shadowBlurScale/4) - 1) / 2))
		if want < 1 {
			want = 1
		}
		if got := shadowBlurRadius(blur); got != want {
			t.Errorf("blur %g: radius = %d, want %d", blur, got, want)
		}
	}
}

// A semi-transparent shadow colour must keep its brightness: Chromium paints
// rgba(255,0,0,0.8) as (255,51,51) on white. Before the premultiplication fix
// the same pixel came out (101,51,51) - the red channel lost 60 % of its value.
func TestOutsetShadowSemiTransparentColorIsNotDarkened(t *testing.T) {
	// A 100x100 box at (40,40) on a 300x300 canvas with the shadow offset by 60 px:
	// the box covers 40..139 and the shadow band lands at 100..199, so (180,180) is
	// shadow only.
	img := renderShadowPNG(t, "60px 60px 0 rgba(255,0,0,0.8)", 300, 300, 100, 100, 40, 40)
	px := img.RGBAAt(180, 180)
	if px.R < 245 || px.G < 40 || px.G > 62 || px.B < 40 || px.B > 62 {
		t.Fatalf("shadow pixel = (%d,%d,%d), want ~(255,51,51)", px.R, px.G, px.B)
	}
}

// An inset shadow is the inside of the box shaded from the edges. At the edge it
// is half the shadow alpha (the blur averages the complement), fading to nothing
// towards the centre; Chromium behaves the same way.
func TestInsetShadowFalloffMatchesBrowser(t *testing.T) {
	const size = 320
	cases := []struct {
		shadow   string
		halfDist int // Chromium's measured distance from the edge to half darkness
	}{
		{"inset 0 0 8px rgba(0,0,0,1)", 3},
		{"inset 0 0 24px rgba(0,0,0,1)", 9},
		{"inset 0 0 48px rgba(0,0,0,1)", 17},
	}
	for _, tc := range cases {
		img := renderShadowPNG(t, tc.shadow, size, size, size, size, 0, 0)
		peak, dist := 0, -1
		for x := 0; x < size; x++ {
			dark := 255 - int(img.RGBAAt(x, size/2).R)
			if x == 0 {
				peak = dark
				continue
			}
			if dist < 0 && dark <= peak/2 {
				dist = x
			}
		}
		if peak < 60 || peak > 200 {
			t.Errorf("%q: edge darkness = %d, want roughly half alpha (~127)", tc.shadow, peak)
		}
		if dist < 0 {
			t.Fatalf("%q: shadow never faded to half darkness", tc.shadow)
		}
		if diff := dist - tc.halfDist; diff < -2 || diff > 2 {
			t.Errorf("%q: half-darkness at %d px, Chromium at %d px", tc.shadow, dist, tc.halfDist)
		}
		if c := img.RGBAAt(size/2, size/2); c.R != 255 || c.G != 255 || c.B != 255 {
			t.Errorf("%q: centre = (%d,%d,%d), want white", tc.shadow, c.R, c.G, c.B)
		}
	}
}

// boxBlurAlpha returns its input unchanged when radius is 0. The caller must not
// release the same pooled buffer twice: that puts one slice in the pool twice and
// the next two acquisitions alias, corrupting unrelated pixels.
func TestShadowBlurZeroRepeatedRendersAreStable(t *testing.T) {
	first := renderShadowPNG(t, "40px 40px 0 rgba(0,0,0,0.5)", 300, 300, 100, 100, 40, 40)
	for i := 0; i < 6; i++ {
		next := renderShadowPNG(t, "40px 40px 0 rgba(0,0,0,0.5)", 300, 300, 100, 100, 40, 40)
		for y := 0; y < 300; y += 7 {
			for x := 0; x < 300; x += 7 {
				if first.RGBAAt(x, y) != next.RGBAAt(x, y) {
					t.Fatalf("render %d differs at (%d,%d): %v vs %v",
						i, x, y, first.RGBAAt(x, y), next.RGBAAt(x, y))
				}
			}
		}
	}
}
