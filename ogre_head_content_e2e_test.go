package ogre

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

// inkBounds returns the bounding box of every non-white pixel, compositing
// over white first because Ogre writes premultiplied RGBA (a transparent pixel
// reads as (0,0,0,0) and would otherwise count as ink).
func inkBounds(img image.Image) (x0, y0, x1, y1 int, ok bool) {
	b := img.Bounds()
	x0, y0, x1, y1 = b.Max.X, b.Max.Y, -1, -1
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, a := img.At(x, y).RGBA()
			over := func(c uint32) uint32 { return (c + 0xffff - a) / 0x101 }
			if over(r) < 250 || over(g) < 250 || over(bl) < 250 {
				x0 = min(x0, x)
				y0 = min(y0, y)
				x1 = max(x1, x)
				y1 = max(y1, y)
			}
		}
	}
	return x0, y0, x1, y1, x1 >= 0
}

// Head content (<style>, <script>, <title>) is not rendered by a browser, so it
// must not reach the layout tree either.
//
// It used to: the fragment parse path converted the element's text into a text
// node, so a document with a <style> block in <head> and a plain <body> painted
// its own CSS across the top of the image — the shape of every real HTML file,
// and the cause of the "heavily corrupted render" reported in #20. Only the
// document path (taken when <body> carries an attribute) dropped <head>.
func TestHeadContentNotRendered(t *testing.T) {
	// The box is the only thing that should leave ink: 100x50 at the origin.
	box := `<div style="width:100px;height:50px;background-color:#ff0000"></div>`
	const wantW, wantH = 100, 50

	cases := []struct {
		name string
		html string
	}{
		{
			name: "style in head",
			html: `<html><head><style>body{margin:0;font-size:40px}</style></head><body>` + box + `</body></html>`,
		},
		{
			name: "script in head",
			html: `<html><head><script>var answer = 42;</script></head><body>` + box + `</body></html>`,
		},
		{
			name: "title in head",
			html: `<html><head><title>Quarterly report</title></head><body>` + box + `</body></html>`,
		},
		{
			name: "style in body",
			html: `<html><body><style>body{margin:0;font-size:40px}</style>` + box + `</body></html>`,
		},
		{
			name: "doctype with meta and style",
			html: `<!DOCTYPE html><html><head><meta charset="utf-8"><style>body{margin:0;font-size:40px}</style></head><body>` + box + `</body></html>`,
		},
		{
			name: "control: no head content",
			html: `<html><body>` + box + `</body></html>`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := Render(c.html, Options{Width: 400, Height: 200, Format: FormatPNG})
			if err != nil {
				t.Fatalf("Render: %v", err)
			}
			img, err := png.Decode(bytes.NewReader(result.Data))
			if err != nil {
				t.Fatalf("decode png: %v", err)
			}

			x0, y0, x1, y1, ok := inkBounds(img)
			if !ok {
				t.Fatal("no ink at all: the box itself is missing")
			}
			gotW, gotH := x1-x0+1, y1-y0+1
			if x0 != 0 || y0 != 0 || gotW != wantW || gotH != wantH {
				t.Fatalf("ink spans x=%d y=%d w=%d h=%d, want x=0 y=0 w=%d h=%d: "+
					"non-rendered head content was laid out as visible text",
					x0, y0, gotW, gotH, wantW, wantH)
			}
		})
	}
}
