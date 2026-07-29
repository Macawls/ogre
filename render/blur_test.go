package render

import (
	"image"
	"math"
	"testing"
)

// blurHRowsFloat is the pre-fast-path reference implementation, kept here
// only to verify the r=1 integer fast path produces identical output.
func blurHRowsFloat(src, dst *image.Alpha, b image.Rectangle, r, y0, y1 int) {
	w := b.Dx()
	div := float64(2*r + 1)
	stride := dst.Stride
	for y := y0; y < y1; y++ {
		sum := 0.0
		for x := -r; x <= r; x++ {
			cx := clampInt(x, 0, w-1)
			sum += float64(src.Pix[y*stride+cx])
		}
		for x := b.Min.X; x < b.Max.X; x++ {
			dst.Pix[y*stride+x] = uint8(math.Round(sum / div))
			nx := clampInt(x+r+1, 0, w-1)
			ox := clampInt(x-r, 0, w-1)
			sum += float64(src.Pix[y*stride+nx]) - float64(src.Pix[y*stride+ox])
		}
	}
}

func TestBlurR1MatchesFloat(t *testing.T) {
	// Fixed-seed deterministic pattern with varied values.
	src := image.NewAlpha(image.Rect(0, 0, 47, 31))
	for i := range src.Pix {
		src.Pix[i] = uint8((i*37 + 13) & 0xFF)
	}
	got := image.NewAlpha(src.Bounds())
	want := image.NewAlpha(src.Bounds())
	blurHRows(src, got, src.Bounds(), 1, 0, 31)
	blurHRowsFloat(src, want, src.Bounds(), 1, 0, 31)
	for i := range got.Pix {
		if got.Pix[i] != want.Pix[i] {
			t.Fatalf("mismatch at %d: fast=%d float=%d", i, got.Pix[i], want.Pix[i])
		}
	}
}
