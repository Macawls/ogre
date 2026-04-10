//go:build ignore

package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
)

func main() {
	size := 64
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(img, img.Bounds(), image.Transparent, image.Point{}, draw.Src)

	cx, cy := float64(size)/2, float64(size)/2
	outerR := float64(size)/2 - 2
	innerR := outerR * 0.4
	gold := color.RGBA{227, 179, 65, 255}

	points := make([][2]float64, 10)
	for i := range 10 {
		angle := float64(i)*math.Pi/5 - math.Pi/2
		r := outerR
		if i%2 == 1 {
			r = innerR
		}
		points[i] = [2]float64{cx + r*math.Cos(angle), cy + r*math.Sin(angle)}
	}

	for y := range size {
		for x := range size {
			px := float64(x) + 0.5
			py := float64(y) + 0.5
			if pointInPolygon(px, py, points[:]) {
				minDist := edgeDistance(px, py, points[:])
				if minDist < 1.0 {
					a := uint8(minDist * 255)
					img.SetRGBA(x, y, color.RGBA{gold.R, gold.G, gold.B, a})
				} else {
					img.SetRGBA(x, y, gold)
				}
			} else {
				minDist := edgeDistance(px, py, points[:])
				if minDist < 1.0 {
					a := uint8((1.0 - minDist) * 255)
					img.SetRGBA(x, y, color.RGBA{gold.R, gold.G, gold.B, a})
				}
			}
		}
	}

	var buf bytes.Buffer
	png.Encode(&buf, img)
	fmt.Printf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(buf.Bytes()))
}

func pointInPolygon(px, py float64, poly [][2]float64) bool {
	inside := false
	n := len(poly)
	j := n - 1
	for i := range n {
		xi, yi := poly[i][0], poly[i][1]
		xj, yj := poly[j][0], poly[j][1]
		if ((yi > py) != (yj > py)) && (px < (xj-xi)*(py-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}
	return inside
}

func edgeDistance(px, py float64, poly [][2]float64) float64 {
	minD := math.MaxFloat64
	n := len(poly)
	for i := range n {
		j := (i + 1) % n
		d := pointToSegment(px, py, poly[i][0], poly[i][1], poly[j][0], poly[j][1])
		if d < minD {
			minD = d
		}
	}
	return minD
}

func pointToSegment(px, py, ax, ay, bx, by float64) float64 {
	dx, dy := bx-ax, by-ay
	if dx == 0 && dy == 0 {
		return math.Hypot(px-ax, py-ay)
	}
	t := ((px-ax)*dx + (py-ay)*dy) / (dx*dx + dy*dy)
	t = math.Max(0, math.Min(1, t))
	return math.Hypot(px-(ax+t*dx), py-(ay+t*dy))
}
