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

// opacityCompositeSlow is the pre-optimization inner loop, kept only in the
// bench file so BenchmarkOpacityCompositeFast can be compared apples-to-apples.
func opacityCompositeSlow(dst, src *image.RGBA, bounds image.Rectangle, opacity float64) {
	for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
		for px := bounds.Min.X; px < bounds.Max.X; px++ {
			off := src.PixOffset(px, py)
			sa := src.Pix[off+3]
			if sa == 0 {
				continue
			}
			na := uint8(float64(sa) * opacity)
			s := color.RGBA{R: src.Pix[off], G: src.Pix[off+1], B: src.Pix[off+2], A: na}
			doff := dst.PixOffset(px, py)
			dr, dg, db, da := blendOver(s.R, s.G, s.B, s.A, dst.Pix[doff], dst.Pix[doff+1], dst.Pix[doff+2], dst.Pix[doff+3])
			dst.Pix[doff] = dr
			dst.Pix[doff+1] = dg
			dst.Pix[doff+2] = db
			dst.Pix[doff+3] = da
		}
	}
}

func opacityCompositeFast(dst, src *image.RGBA, bounds image.Rectangle, opacity float64) {
	stride := dst.Stride
	srcPix := src.Pix
	dstPix := dst.Pix
	for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
		rowOff := py * stride
		for px := bounds.Min.X; px < bounds.Max.X; px++ {
			off := rowOff + px*4
			sa := srcPix[off+3]
			if sa == 0 {
				continue
			}
			na := uint8(float64(sa) * opacity)
			dr, dg, db, da := blendOver(
				srcPix[off], srcPix[off+1], srcPix[off+2], na,
				dstPix[off], dstPix[off+1], dstPix[off+2], dstPix[off+3],
			)
			dstPix[off] = dr
			dstPix[off+1] = dg
			dstPix[off+2] = db
			dstPix[off+3] = da
		}
	}
}

func opacityBenchSrc() (dst, src *image.RGBA, bounds image.Rectangle) {
	dst = image.NewRGBA(image.Rect(0, 0, 1200, 630))
	src = image.NewRGBA(image.Rect(0, 0, 1200, 630))
	// Fill src with a mixed alpha pattern so the fast-skip (sa==0) branch
	// exercises both paths.
	for y := 0; y < 630; y++ {
		for x := 0; x < 1200; x++ {
			off := y*src.Stride + x*4
			src.Pix[off] = uint8(x & 0xFF)
			src.Pix[off+1] = uint8(y & 0xFF)
			src.Pix[off+2] = 128
			if (x+y)&7 == 0 {
				src.Pix[off+3] = 0
			} else {
				src.Pix[off+3] = 200
			}
		}
	}
	bounds = image.Rect(100, 50, 900, 580)
	return
}

func BenchmarkOpacityCompositeSlow(b *testing.B) {
	dst, src, bounds := opacityBenchSrc()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		opacityCompositeSlow(dst, src, bounds, 0.5)
	}
}

func BenchmarkOpacityCompositeFast(b *testing.B) {
	dst, src, bounds := opacityBenchSrc()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		opacityCompositeFast(dst, src, bounds, 0.5)
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
