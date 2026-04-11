package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"math"
)

type DiffResult struct {
	MatchPct float64 `json:"matchPct"`
	Heatmap  string  `json:"heatmap"`
}

func pixelDiff(pngA, pngB []byte) (*DiffResult, error) {
	imgA, err := png.Decode(bytes.NewReader(pngA))
	if err != nil {
		return nil, err
	}
	imgB, err := png.Decode(bytes.NewReader(pngB))
	if err != nil {
		return nil, err
	}

	boundsA := imgA.Bounds()
	boundsB := imgB.Bounds()
	w := boundsA.Dx()
	h := boundsA.Dy()
	if boundsB.Dx() < w {
		w = boundsB.Dx()
	}
	if boundsB.Dy() < h {
		h = boundsB.Dy()
	}

	heatmap := image.NewRGBA(image.Rect(0, 0, w, h))
	matching := 0
	total := w * h
	threshold := 50.0

	for y := range h {
		for x := range w {
			rA, gA, bA, _ := imgA.At(x+boundsA.Min.X, y+boundsA.Min.Y).RGBA()
			rB, gB, bB, _ := imgB.At(x+boundsB.Min.X, y+boundsB.Min.Y).RGBA()

			dr := float64(rA>>8) - float64(rB>>8)
			dg := float64(gA>>8) - float64(gB>>8)
			db := float64(bA>>8) - float64(bB>>8)
			dist := math.Sqrt(dr*dr + dg*dg + db*db)

			if dist < threshold {
				matching++
				heatmap.SetRGBA(x, y, color.RGBA{0, 0, 0, 40})
			} else {
				intensity := uint8(math.Min(255, dist))
				heatmap.SetRGBA(x, y, color.RGBA{intensity, 0, 0, 200})
			}
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, heatmap)

	return &DiffResult{
		MatchPct: float64(matching) / float64(total) * 100,
		Heatmap:  base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}
