package ogre

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"testing"
)

// dejaVuSans loads the system DejaVu Sans face used by the browser comparison.
func dejaVuSans(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf")
	if err != nil {
		t.Skipf("DejaVu Sans not available: %v", err)
	}
	return data
}

// decodePNG decodes PNG bytes returned by Render.
func decodePNG(t *testing.T, data []byte) image.Image {
	t.Helper()
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode png: %v", err)
	}
	return img
}

// A document that paints its background on <body> must actually get that
// background. Before the fix the <body> element and its style were dropped by
// fragment parsing, so the canvas stayed blank.
func TestRenderFullDocumentBodyBackground(t *testing.T) {
	html := `<!DOCTYPE html><html><head><meta charset="utf-8"></head>` +
		`<body style="margin:0;padding:0;width:300px;height:200px;background:#ff0000;">` +
		`<div style="padding:10px;font-size:20px;">X</div>` +
		`</body></html>`

	result, err := Render(html, Options{Width: 300, Height: 200, Format: FormatPNG})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	img := decodePNG(t, result.Data)
	r, g, b, _ := img.At(150, 150).RGBA()
	if r>>8 != 0xff || g>>8 != 0x00 || b>>8 != 0x00 {
		t.Fatalf("background pixel = (%d,%d,%d), want (255,0,0): <body> style was dropped",
			r>>8, g>>8, b>>8)
	}
}

// The <body> font must reach the shaper. Declaring the same family on a <div>
// already worked; declaring it on <body> used to fall back to the built-in
// font, which made bold text render ~12% narrower than the browser.
func TestRenderFullDocumentBodyFontAppliesToChildren(t *testing.T) {
	fonts := []FontSource{
		{Name: "DejaVu Sans", Weight: 400, Data: dejaVuSans(t)},
	}
	render := func(html string) int {
		result, err := Render(html, Options{Width: 900, Height: 400, Format: FormatPNG, Fonts: fonts})
		if err != nil {
			t.Fatalf("Render: %v", err)
		}
		img := decodePNG(t, result.Data)
		bounds := img.Bounds()
		minX, maxX := bounds.Max.X, -1
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				if r>>8 < 250 || g>>8 < 250 || b>>8 < 250 {
					if x < minX {
						minX = x
					}
					if x > maxX {
						maxX = x
					}
				}
			}
		}
		if maxX < 0 {
			t.Fatal("no ink rendered")
		}
		return maxX - minX + 1
	}

	bodyDecl := `<body style="margin:0;padding:0;width:900px;height:400px;background:#fff;font-family:'DejaVu Sans';">` +
		`<div style="padding:10px;font-size:20px;">WWWWWWWW</div></body>`
	divDecl := `<body style="margin:0;padding:0;width:900px;height:400px;background:#fff;">` +
		`<div style="font-family:'DejaVu Sans';"><div style="padding:10px;font-size:20px;">WWWWWWWW</div></div></body>`

	viaBody := render(bodyDecl)
	viaDiv := render(divDecl)

	if viaBody != viaDiv {
		t.Fatalf("text width via <body> = %d, via <div> = %d: body font-family did not reach the shaper",
			viaBody, viaDiv)
	}
}
