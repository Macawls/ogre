package render

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/macawls/ogre/style"
)

func BenchmarkRoundedMask(b *testing.B) {
	for i := 0; i < b.N; i++ {
		roundedMask(400, 300, 24, 24, 24, 24)
	}
}

func BenchmarkLinearGradient(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 1200, 630))
	r := &PNGRenderer{img: img}
	g := style.Gradient{
		Type:  style.LinearGradient,
		Angle: 135,
		Stops: []style.ColorStop{
			{Color: style.Color{R: 102, G: 126, B: 234, A: 1}, Position: 0},
			{Color: style.Color{R: 118, G: 75, B: 162, A: 1}, Position: 1},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.renderLinearGradientPNG(g, 0, 0, 1200, 630)
	}
}

func BenchmarkRadialGradient(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 1200, 630))
	r := &PNGRenderer{img: img}
	g := style.Gradient{
		Type:      style.RadialGradient,
		PositionX: 50,
		PositionY: 50,
		Stops: []style.ColorStop{
			{Color: style.Color{R: 102, G: 126, B: 234, A: 1}, Position: 0},
			{Color: style.Color{R: 118, G: 75, B: 162, A: 1}, Position: 1},
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.renderRadialGradientPNG(g, 0, 0, 1200, 630)
	}
}

func BenchmarkBoxBlurAlpha(b *testing.B) {
	src := image.NewAlpha(image.Rect(0, 0, 500, 350))
	for y := 50; y < 300; y++ {
		for x := 50; x < 450; x++ {
			src.SetAlpha(x, y, color.Alpha{A: 200})
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		boxBlurAlpha(src, 25)
	}
}

func BenchmarkBoxBlurAlphaR1(b *testing.B) {
	src := image.NewAlpha(image.Rect(0, 0, 500, 350))
	for y := 50; y < 300; y++ {
		for x := 50; x < 450; x++ {
			src.SetAlpha(x, y, color.Alpha{A: 200})
		}
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		boxBlurAlpha(src, 1)
	}
}

func BenchmarkFillRect(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 1200, 630))
	c := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fillRect(img, 100, 100, 400, 300, c)
	}
}

func BenchmarkPNGEncode(b *testing.B) {
	img := image.NewRGBA(image.Rect(0, 0, 1200, 630))
	for y := range 630 {
		for x := range 1200 {
			off := y*img.Stride + x*4
			img.Pix[off] = uint8(x & 0xFF)
			img.Pix[off+1] = uint8(y & 0xFF)
			img.Pix[off+2] = 128
			img.Pix[off+3] = 255
		}
	}
	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		png.Encode(&buf, img)
	}
}

func benchImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, 1200, 630))
	for y := range 630 {
		for x := range 1200 {
			off := y*img.Stride + x*4
			img.Pix[off] = uint8(x & 0xFF)
			img.Pix[off+1] = uint8(y & 0xFF)
			img.Pix[off+2] = 128
			img.Pix[off+3] = 255
		}
	}
	return img
}

func BenchmarkPNGEncodeBestSpeed(b *testing.B) {
	img := benchImage()
	enc := png.Encoder{CompressionLevel: png.BestSpeed, BufferPool: pngEncBufPool}
	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := enc.Encode(&buf, img); err != nil {
			b.Fatal(err)
		}
	}
	b.ReportMetric(float64(buf.Len()), "outputBytes")
}

func BenchmarkPNGEncodeBestCompression(b *testing.B) {
	img := benchImage()
	enc := png.Encoder{CompressionLevel: png.BestCompression, BufferPool: pngEncBufPool}
	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		if err := enc.Encode(&buf, img); err != nil {
			b.Fatal(err)
		}
	}
	b.ReportMetric(float64(buf.Len()), "outputBytes")
}
