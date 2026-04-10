package render

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/image/vector"
)

type svgElement struct {
	XMLName  xml.Name
	Attrs    []xml.Attr   `xml:",any,attr"`
	Children []svgElement `xml:",any"`
	CharData string       `xml:",chardata"`
}

func rasterizeSVG(data []byte, targetW, targetH int) (image.Image, error) {
	var root svgElement
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	if root.XMLName.Local != "svg" {
		return nil, fmt.Errorf("not an SVG element")
	}

	attrs := attrMap(root.Attrs)

	var vbX, vbY, vbW, vbH float64
	if vb, ok := attrs["viewBox"]; ok {
		parts := strings.Fields(strings.ReplaceAll(vb, ",", " "))
		if len(parts) == 4 {
			vbX, _ = strconv.ParseFloat(parts[0], 64)
			vbY, _ = strconv.ParseFloat(parts[1], 64)
			vbW, _ = strconv.ParseFloat(parts[2], 64)
			vbH, _ = strconv.ParseFloat(parts[3], 64)
		}
	}
	if vbW == 0 || vbH == 0 {
		vbW = parseDimension(attrs["width"])
		vbH = parseDimension(attrs["height"])
	}
	if vbW == 0 || vbH == 0 {
		vbW = float64(targetW)
		vbH = float64(targetH)
	}

	if targetW == 0 || targetH == 0 {
		targetW = int(vbW)
		targetH = int(vbH)
	}

	sx := float64(targetW) / vbW
	sy := float64(targetH) / vbH

	img := image.NewRGBA(image.Rect(0, 0, targetW, targetH))

	if fill, ok := attrs["fill"]; ok && fill != "none" {
		c := parseSVGColor(fill, color.RGBA{0, 0, 0, 0})
		if c.A > 0 {
			draw.Draw(img, img.Bounds(), image.NewUniform(c), image.Point{}, draw.Src)
		}
	}

	defaultFill := color.RGBA{0, 0, 0, 255}
	if f, ok := attrs["fill"]; ok {
		if f == "none" {
			defaultFill = color.RGBA{}
		} else {
			defaultFill = parseSVGColor(f, color.RGBA{0, 0, 0, 255})
		}
	}

	ctx := &svgRenderCtx{
		img:     img,
		sx:      sx,
		sy:      sy,
		ox:      -vbX * sx,
		oy:      -vbY * sy,
		fill:    defaultFill,
		targetW: targetW,
		targetH: targetH,
	}

	ctx.renderChildren(root.Children)
	return img, nil
}

type svgRenderCtx struct {
	img     *image.RGBA
	sx, sy  float64
	ox, oy  float64
	fill    color.RGBA
	targetW int
	targetH int
}

func (ctx *svgRenderCtx) renderChildren(children []svgElement) {
	for _, child := range children {
		ctx.renderElement(child)
	}
}

func (ctx *svgRenderCtx) renderElement(el svgElement) {
	attrs := attrMap(el.Attrs)
	fill := ctx.resolveFill(attrs)
	stroke, strokeWidth := ctx.resolveStroke(attrs)

	switch el.XMLName.Local {
	case "g":
		if t, ok := attrs["transform"]; ok {
			dx, dy := parseTranslate(t)
			saved := *ctx
			ctx.ox += dx * ctx.sx
			ctx.oy += dy * ctx.sy
			ctx.renderChildren(el.Children)
			*ctx = saved
		} else {
			ctx.renderChildren(el.Children)
		}
	case "path":
		d, ok := attrs["d"]
		if !ok {
			return
		}
		if fill.A > 0 {
			if attrs["fill-rule"] == "evenodd" || attrs["clip-rule"] == "evenodd" {
				ctx.renderPathEvenOdd(d, fill)
			} else {
				ctx.renderPath(d, fill)
			}
		}
		if stroke.A > 0 && strokeWidth > 0 {
			ctx.renderPathStroke(d, stroke, strokeWidth)
		}
	case "rect":
		if fill.A > 0 {
			ctx.renderRect(attrs, fill)
		}
	case "circle":
		if fill.A > 0 {
			ctx.renderCircle(attrs, fill)
		}
		if stroke.A > 0 && strokeWidth > 0 {
			cx := parseAttrFloat(attrs, "cx")
			cy := parseAttrFloat(attrs, "cy")
			radius := parseAttrFloat(attrs, "r")
			ctx.renderCircleStroke(cx, cy, radius, strokeWidth, stroke)
		}
	case "ellipse":
		if fill.A > 0 {
			ctx.renderEllipse(attrs, fill)
		}
	case "polygon":
		pts, ok := attrs["points"]
		if !ok {
			return
		}
		if fill.A > 0 {
			ctx.renderPolygon(pts, fill)
		}
	case "polyline":
		pts, ok := attrs["points"]
		if !ok {
			return
		}
		if fill.A > 0 {
			ctx.renderPolygon(pts, fill)
		}
	case "line":
		if stroke.A > 0 && strokeWidth > 0 {
			x1 := parseAttrFloat(attrs, "x1")
			y1 := parseAttrFloat(attrs, "y1")
			x2 := parseAttrFloat(attrs, "x2")
			y2 := parseAttrFloat(attrs, "y2")
			ctx.renderLineStroke(x1, y1, x2, y2, strokeWidth, stroke)
		}
	}
}

func (ctx *svgRenderCtx) resolveFill(attrs map[string]string) color.RGBA {
	if f, ok := attrs["fill"]; ok {
		if f == "none" {
			return color.RGBA{}
		}
		return parseSVGColor(f, ctx.fill)
	}
	return ctx.fill
}

func (ctx *svgRenderCtx) resolveStroke(attrs map[string]string) (color.RGBA, float64) {
	s, ok := attrs["stroke"]
	if !ok || s == "none" {
		return color.RGBA{}, 0
	}
	c := parseSVGColor(s, ctx.fill)
	w := 1.0
	if sw, ok := attrs["stroke-width"]; ok {
		w = parseDimension(sw)
	}
	return c, w
}

func (ctx *svgRenderCtx) renderCircleStroke(cx, cy, radius, strokeWidth float64, stroke color.RGBA) {
	outerR := radius + strokeWidth/2
	innerR := radius - strokeWidth/2
	if innerR < 0 {
		innerR = 0
	}

	r := vector.NewRasterizer(ctx.targetW, ctx.targetH)

	appendEllipse(r, ctx, cx, cy, outerR, outerR, false)
	if innerR > 0 {
		appendEllipse(r, ctx, cx, cy, innerR, innerR, true)
	}

	r.Draw(ctx.img, ctx.img.Bounds(), image.NewUniform(stroke), image.Point{})
}

func appendEllipse(r *vector.Rasterizer, ctx *svgRenderCtx, cx, cy, rx, ry float64, reverse bool) {
	const n = 16
	type pt struct{ x, y float64 }
	var pts [n]pt
	for i := range n {
		angle := 2 * math.Pi * float64(i) / n
		pts[i] = pt{cx + rx*math.Cos(angle), cy + ry*math.Sin(angle)}
	}

	if reverse {
		r.MoveTo(ctx.tx(pts[0].x), ctx.ty(pts[0].y))
		for i := n - 1; i >= 1; i-- {
			a1 := 2 * math.Pi * float64(i) / n
			a0 := 2 * math.Pi * float64(i-1) / n
			k := 4.0 / 3.0 * math.Tan((a1-a0)/4)
			cp2x := cx + rx*(math.Cos(a1)+k*math.Sin(a1))
			cp2y := cy + ry*(math.Sin(a1)-k*math.Cos(a1))
			cp1x := cx + rx*(math.Cos(a0)-k*math.Sin(a0))
			cp1y := cy + ry*(math.Sin(a0)+k*math.Cos(a0))
			r.CubeTo(ctx.tx(cp2x), ctx.ty(cp2y), ctx.tx(cp1x), ctx.ty(cp1y), ctx.tx(pts[i-1].x), ctx.ty(pts[i-1].y))
		}
	} else {
		r.MoveTo(ctx.tx(pts[0].x), ctx.ty(pts[0].y))
		for i := 1; i < n; i++ {
			a0 := 2 * math.Pi * float64(i-1) / n
			a1 := 2 * math.Pi * float64(i) / n
			k := 4.0 / 3.0 * math.Tan((a1-a0)/4)
			cp1x := cx + rx*(math.Cos(a0)-k*math.Sin(a0))
			cp1y := cy + ry*(math.Sin(a0)+k*math.Cos(a0))
			cp2x := cx + rx*(math.Cos(a1)+k*math.Sin(a1))
			cp2y := cy + ry*(math.Sin(a1)-k*math.Cos(a1))
			r.CubeTo(ctx.tx(cp1x), ctx.ty(cp1y), ctx.tx(cp2x), ctx.ty(cp2y), ctx.tx(pts[i].x), ctx.ty(pts[i].y))
		}
	}
	r.ClosePath()
}

func (ctx *svgRenderCtx) renderPathEvenOdd(d string, fill color.RGBA) {
	subPaths := splitSubPaths(d)
	combined := image.NewAlpha(image.Rect(0, 0, ctx.targetW, ctx.targetH))

	for _, sp := range subPaths {
		r := vector.NewRasterizer(ctx.targetW, ctx.targetH)
		cmds := parseSVGPath(sp)
		var cx, cy, startX, startY, lastCPX, lastCPY float64
		var lastCmd byte
		for _, cmd := range cmds {
			switch cmd.cmd {
			case 'M':
				for i := 0; i < len(cmd.args)-1; i += 2 {
					cx, cy = cmd.args[i], cmd.args[i+1]
					if i == 0 {
						startX, startY = cx, cy
						r.MoveTo(ctx.tx(cx), ctx.ty(cy))
					} else {
						r.LineTo(ctx.tx(cx), ctx.ty(cy))
					}
				}
			case 'm':
				for i := 0; i < len(cmd.args)-1; i += 2 {
					cx += cmd.args[i]
					cy += cmd.args[i+1]
					if i == 0 {
						startX, startY = cx, cy
						r.MoveTo(ctx.tx(cx), ctx.ty(cy))
					} else {
						r.LineTo(ctx.tx(cx), ctx.ty(cy))
					}
				}
			case 'L':
				for i := 0; i < len(cmd.args)-1; i += 2 {
					cx, cy = cmd.args[i], cmd.args[i+1]
					r.LineTo(ctx.tx(cx), ctx.ty(cy))
				}
			case 'l':
				for i := 0; i < len(cmd.args)-1; i += 2 {
					cx += cmd.args[i]
					cy += cmd.args[i+1]
					r.LineTo(ctx.tx(cx), ctx.ty(cy))
				}
			case 'H':
				for _, a := range cmd.args {
					cx = a
					r.LineTo(ctx.tx(cx), ctx.ty(cy))
				}
			case 'h':
				for _, a := range cmd.args {
					cx += a
					r.LineTo(ctx.tx(cx), ctx.ty(cy))
				}
			case 'V':
				for _, a := range cmd.args {
					cy = a
					r.LineTo(ctx.tx(cx), ctx.ty(cy))
				}
			case 'v':
				for _, a := range cmd.args {
					cy += a
					r.LineTo(ctx.tx(cx), ctx.ty(cy))
				}
			case 'C':
				for i := 0; i < len(cmd.args)-5; i += 6 {
					x1, y1 := cmd.args[i], cmd.args[i+1]
					x2, y2 := cmd.args[i+2], cmd.args[i+3]
					cx, cy = cmd.args[i+4], cmd.args[i+5]
					lastCPX, lastCPY = x2, y2
					r.CubeTo(ctx.tx(x1), ctx.ty(y1), ctx.tx(x2), ctx.ty(y2), ctx.tx(cx), ctx.ty(cy))
				}
			case 'c':
				for i := 0; i < len(cmd.args)-5; i += 6 {
					x1, y1 := cx+cmd.args[i], cy+cmd.args[i+1]
					x2, y2 := cx+cmd.args[i+2], cy+cmd.args[i+3]
					cx += cmd.args[i+4]
					cy += cmd.args[i+5]
					lastCPX, lastCPY = x2, y2
					r.CubeTo(ctx.tx(x1), ctx.ty(y1), ctx.tx(x2), ctx.ty(y2), ctx.tx(cx), ctx.ty(cy))
				}
			case 'Z', 'z':
				r.ClosePath()
				cx, cy = startX, startY
			}
			lastCmd = cmd.cmd
		}
		_, _ = lastCPX, lastCPY
		if lastCmd != 'Z' && lastCmd != 'z' && lastCmd != 0 {
			r.ClosePath()
		}

		tmp := image.NewAlpha(image.Rect(0, 0, ctx.targetW, ctx.targetH))
		r.Draw(tmp, tmp.Bounds(), image.NewUniform(color.Alpha{255}), image.Point{})

		for y := range ctx.targetH {
			for x := range ctx.targetW {
				ta := tmp.AlphaAt(x, y).A
				if ta > 128 {
					ca := combined.AlphaAt(x, y).A
					if ca > 128 {
						combined.SetAlpha(x, y, color.Alpha{0})
					} else {
						combined.SetAlpha(x, y, color.Alpha{255})
					}
				}
			}
		}
	}

	draw.DrawMask(ctx.img, ctx.img.Bounds(), image.NewUniform(fill), image.Point{}, combined, image.Point{}, draw.Over)
}

func splitSubPaths(d string) []string {
	var paths []string
	var current strings.Builder
	i := 0
	for i < len(d) {
		if i > 0 && (d[i] == 'M' || d[i] == 'm') && current.Len() > 0 {
			paths = append(paths, current.String())
			current.Reset()
		}
		current.WriteByte(d[i])
		i++
	}
	if current.Len() > 0 {
		paths = append(paths, current.String())
	}
	return paths
}

func (ctx *svgRenderCtx) renderPathStroke(d string, stroke color.RGBA, strokeWidth float64) {
	cmds := parseSVGPath(d)
	var cx, cy float64
	var startX, startY float64
	var lastCPX, lastCPY float64
	var lastCmd byte
	hw := strokeWidth / 2

	segments := collectPathSegments(cmds, &cx, &cy, &startX, &startY, &lastCPX, &lastCPY, &lastCmd)

	for _, seg := range segments {
		for i := 0; i < len(seg)-1; i++ {
			ctx.renderLineStroke(seg[i][0], seg[i][1], seg[i+1][0], seg[i+1][1], hw*2, stroke)
		}
	}
}

func collectPathSegments(cmds []pathCmd, cx, cy, startX, startY, lastCPX, lastCPY *float64, lastCmd *byte) [][][2]float64 {
	var segments [][][2]float64
	var current [][2]float64

	for _, cmd := range cmds {
		switch cmd.cmd {
		case 'M':
			if len(current) > 0 {
				segments = append(segments, current)
			}
			current = nil
			for i := 0; i < len(cmd.args)-1; i += 2 {
				*cx, *cy = cmd.args[i], cmd.args[i+1]
				if i == 0 {
					*startX, *startY = *cx, *cy
				}
				current = append(current, [2]float64{*cx, *cy})
			}
		case 'm':
			if len(current) > 0 {
				segments = append(segments, current)
			}
			current = nil
			for i := 0; i < len(cmd.args)-1; i += 2 {
				*cx += cmd.args[i]
				*cy += cmd.args[i+1]
				if i == 0 {
					*startX, *startY = *cx, *cy
				}
				current = append(current, [2]float64{*cx, *cy})
			}
		case 'L':
			for i := 0; i < len(cmd.args)-1; i += 2 {
				*cx, *cy = cmd.args[i], cmd.args[i+1]
				current = append(current, [2]float64{*cx, *cy})
			}
		case 'l':
			for i := 0; i < len(cmd.args)-1; i += 2 {
				*cx += cmd.args[i]
				*cy += cmd.args[i+1]
				current = append(current, [2]float64{*cx, *cy})
			}
		case 'H':
			for _, a := range cmd.args {
				*cx = a
				current = append(current, [2]float64{*cx, *cy})
			}
		case 'h':
			for _, a := range cmd.args {
				*cx += a
				current = append(current, [2]float64{*cx, *cy})
			}
		case 'V':
			for _, a := range cmd.args {
				*cy = a
				current = append(current, [2]float64{*cx, *cy})
			}
		case 'v':
			for _, a := range cmd.args {
				*cy += a
				current = append(current, [2]float64{*cx, *cy})
			}
		case 'C':
			for i := 0; i < len(cmd.args)-5; i += 6 {
				*lastCPX, *lastCPY = cmd.args[i+2], cmd.args[i+3]
				*cx, *cy = cmd.args[i+4], cmd.args[i+5]
				current = append(current, [2]float64{*cx, *cy})
			}
		case 'c':
			for i := 0; i < len(cmd.args)-5; i += 6 {
				*lastCPX = *cx + cmd.args[i+2]
				*lastCPY = *cy + cmd.args[i+3]
				*cx += cmd.args[i+4]
				*cy += cmd.args[i+5]
				current = append(current, [2]float64{*cx, *cy})
			}
		case 'Z', 'z':
			current = append(current, [2]float64{*startX, *startY})
			*cx, *cy = *startX, *startY
		}
		*lastCmd = cmd.cmd
	}
	if len(current) > 0 {
		segments = append(segments, current)
	}
	return segments
}

func (ctx *svgRenderCtx) renderLineStroke(x1, y1, x2, y2, strokeWidth float64, stroke color.RGBA) {
	dx := x2 - x1
	dy := y2 - y1
	length := math.Sqrt(dx*dx + dy*dy)
	if length == 0 {
		return
	}
	hw := strokeWidth / 2
	px := -dy / length * hw
	py := dx / length * hw

	r := vector.NewRasterizer(ctx.targetW, ctx.targetH)
	r.MoveTo(ctx.tx(x1+px), ctx.ty(y1+py))
	r.LineTo(ctx.tx(x2+px), ctx.ty(y2+py))
	r.LineTo(ctx.tx(x2-px), ctx.ty(y2-py))
	r.LineTo(ctx.tx(x1-px), ctx.ty(y1-py))
	r.ClosePath()
	r.Draw(ctx.img, ctx.img.Bounds(), image.NewUniform(stroke), image.Point{})
}

func (ctx *svgRenderCtx) tx(x float64) float32 { return float32(x*ctx.sx + ctx.ox) }
func (ctx *svgRenderCtx) ty(y float64) float32 { return float32(y*ctx.sy + ctx.oy) }

func parseTranslate(s string) (float64, float64) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "translate(") {
		return 0, 0
	}
	s = strings.TrimPrefix(s, "translate(")
	s = strings.TrimSuffix(s, ")")
	parts := strings.Fields(strings.ReplaceAll(s, ",", " "))
	var dx, dy float64
	if len(parts) >= 1 {
		dx, _ = strconv.ParseFloat(parts[0], 64)
	}
	if len(parts) >= 2 {
		dy, _ = strconv.ParseFloat(parts[1], 64)
	}
	return dx, dy
}

func (ctx *svgRenderCtx) renderPath(d string, fill color.RGBA) {
	r := vector.NewRasterizer(ctx.targetW, ctx.targetH)
	cmds := parseSVGPath(d)
	var cx, cy float64
	var startX, startY float64
	var lastCPX, lastCPY float64
	var lastCmd byte

	for _, cmd := range cmds {
		switch cmd.cmd {
		case 'M':
			for i := 0; i < len(cmd.args)-1; i += 2 {
				cx, cy = cmd.args[i], cmd.args[i+1]
				if i == 0 {
					startX, startY = cx, cy
					r.MoveTo(ctx.tx(cx), ctx.ty(cy))
				} else {
					r.LineTo(ctx.tx(cx), ctx.ty(cy))
				}
			}
		case 'm':
			for i := 0; i < len(cmd.args)-1; i += 2 {
				cx += cmd.args[i]
				cy += cmd.args[i+1]
				if i == 0 {
					startX, startY = cx, cy
					r.MoveTo(ctx.tx(cx), ctx.ty(cy))
				} else {
					r.LineTo(ctx.tx(cx), ctx.ty(cy))
				}
			}
		case 'L':
			for i := 0; i < len(cmd.args)-1; i += 2 {
				cx, cy = cmd.args[i], cmd.args[i+1]
				r.LineTo(ctx.tx(cx), ctx.ty(cy))
			}
		case 'l':
			for i := 0; i < len(cmd.args)-1; i += 2 {
				cx += cmd.args[i]
				cy += cmd.args[i+1]
				r.LineTo(ctx.tx(cx), ctx.ty(cy))
			}
		case 'H':
			for _, a := range cmd.args {
				cx = a
				r.LineTo(ctx.tx(cx), ctx.ty(cy))
			}
		case 'h':
			for _, a := range cmd.args {
				cx += a
				r.LineTo(ctx.tx(cx), ctx.ty(cy))
			}
		case 'V':
			for _, a := range cmd.args {
				cy = a
				r.LineTo(ctx.tx(cx), ctx.ty(cy))
			}
		case 'v':
			for _, a := range cmd.args {
				cy += a
				r.LineTo(ctx.tx(cx), ctx.ty(cy))
			}
		case 'C':
			for i := 0; i < len(cmd.args)-5; i += 6 {
				x1, y1 := cmd.args[i], cmd.args[i+1]
				x2, y2 := cmd.args[i+2], cmd.args[i+3]
				cx, cy = cmd.args[i+4], cmd.args[i+5]
				lastCPX, lastCPY = x2, y2
				r.CubeTo(ctx.tx(x1), ctx.ty(y1), ctx.tx(x2), ctx.ty(y2), ctx.tx(cx), ctx.ty(cy))
			}
		case 'c':
			for i := 0; i < len(cmd.args)-5; i += 6 {
				x1, y1 := cx+cmd.args[i], cy+cmd.args[i+1]
				x2, y2 := cx+cmd.args[i+2], cy+cmd.args[i+3]
				cx += cmd.args[i+4]
				cy += cmd.args[i+5]
				lastCPX, lastCPY = x2, y2
				r.CubeTo(ctx.tx(x1), ctx.ty(y1), ctx.tx(x2), ctx.ty(y2), ctx.tx(cx), ctx.ty(cy))
			}
		case 'S':
			for i := 0; i < len(cmd.args)-3; i += 4 {
				x1, y1 := reflectCP(cx, cy, lastCPX, lastCPY, lastCmd)
				x2, y2 := cmd.args[i], cmd.args[i+1]
				cx, cy = cmd.args[i+2], cmd.args[i+3]
				lastCPX, lastCPY = x2, y2
				r.CubeTo(ctx.tx(x1), ctx.ty(y1), ctx.tx(x2), ctx.ty(y2), ctx.tx(cx), ctx.ty(cy))
			}
		case 's':
			for i := 0; i < len(cmd.args)-3; i += 4 {
				x1, y1 := reflectCP(cx, cy, lastCPX, lastCPY, lastCmd)
				x2, y2 := cx+cmd.args[i], cy+cmd.args[i+1]
				cx += cmd.args[i+2]
				cy += cmd.args[i+3]
				lastCPX, lastCPY = x2, y2
				r.CubeTo(ctx.tx(x1), ctx.ty(y1), ctx.tx(x2), ctx.ty(y2), ctx.tx(cx), ctx.ty(cy))
			}
		case 'Q':
			for i := 0; i < len(cmd.args)-3; i += 4 {
				x1, y1 := cmd.args[i], cmd.args[i+1]
				cx, cy = cmd.args[i+2], cmd.args[i+3]
				lastCPX, lastCPY = x1, y1
				r.QuadTo(ctx.tx(x1), ctx.ty(y1), ctx.tx(cx), ctx.ty(cy))
			}
		case 'q':
			for i := 0; i < len(cmd.args)-3; i += 4 {
				x1, y1 := cx+cmd.args[i], cy+cmd.args[i+1]
				cx += cmd.args[i+2]
				cy += cmd.args[i+3]
				lastCPX, lastCPY = x1, y1
				r.QuadTo(ctx.tx(x1), ctx.ty(y1), ctx.tx(cx), ctx.ty(cy))
			}
		case 'T':
			for i := 0; i < len(cmd.args)-1; i += 2 {
				x1, y1 := reflectCP(cx, cy, lastCPX, lastCPY, lastCmd)
				cx, cy = cmd.args[i], cmd.args[i+1]
				lastCPX, lastCPY = x1, y1
				r.QuadTo(ctx.tx(x1), ctx.ty(y1), ctx.tx(cx), ctx.ty(cy))
			}
		case 't':
			for i := 0; i < len(cmd.args)-1; i += 2 {
				x1, y1 := reflectCP(cx, cy, lastCPX, lastCPY, lastCmd)
				cx += cmd.args[i]
				cy += cmd.args[i+1]
				lastCPX, lastCPY = x1, y1
				r.QuadTo(ctx.tx(x1), ctx.ty(y1), ctx.tx(cx), ctx.ty(cy))
			}
		case 'A':
			for i := 0; i < len(cmd.args)-6; i += 7 {
				rx, ry := cmd.args[i], cmd.args[i+1]
				rot := cmd.args[i+2]
				largeArc := cmd.args[i+3] != 0
				sweep := cmd.args[i+4] != 0
				ex, ey := cmd.args[i+5], cmd.args[i+6]
				arcToCubic(r, ctx, cx, cy, rx, ry, rot, largeArc, sweep, ex, ey)
				cx, cy = ex, ey
			}
		case 'a':
			for i := 0; i < len(cmd.args)-6; i += 7 {
				rx, ry := cmd.args[i], cmd.args[i+1]
				rot := cmd.args[i+2]
				largeArc := cmd.args[i+3] != 0
				sweep := cmd.args[i+4] != 0
				ex, ey := cx+cmd.args[i+5], cy+cmd.args[i+6]
				arcToCubic(r, ctx, cx, cy, rx, ry, rot, largeArc, sweep, ex, ey)
				cx, cy = ex, ey
			}
		case 'Z', 'z':
			r.ClosePath()
			cx, cy = startX, startY
		}
		lastCmd = cmd.cmd
	}

	if lastCmd != 'Z' && lastCmd != 'z' && lastCmd != 0 {
		r.ClosePath()
	}

	r.Draw(ctx.img, ctx.img.Bounds(), image.NewUniform(fill), image.Point{})
}

func reflectCP(cx, cy, cpx, cpy float64, lastCmd byte) (float64, float64) {
	lc := lastCmd | 0x20
	if lc == 'c' || lc == 's' || lc == 'q' || lc == 't' {
		return 2*cx - cpx, 2*cy - cpy
	}
	return cx, cy
}

func (ctx *svgRenderCtx) renderRect(attrs map[string]string, fill color.RGBA) {
	x := parseAttrFloat(attrs, "x")
	y := parseAttrFloat(attrs, "y")
	w := parseAttrFloat(attrs, "width")
	h := parseAttrFloat(attrs, "height")
	rx := parseAttrFloat(attrs, "rx")
	ry := parseAttrFloat(attrs, "ry")
	if ry == 0 {
		ry = rx
	}
	if rx == 0 {
		rx = ry
	}

	r := vector.NewRasterizer(ctx.targetW, ctx.targetH)

	if rx > 0 || ry > 0 {
		rx = math.Min(rx, w/2)
		ry = math.Min(ry, h/2)
		r.MoveTo(ctx.tx(x+rx), ctx.ty(y))
		r.LineTo(ctx.tx(x+w-rx), ctx.ty(y))
		arcQuarter(r, ctx, x+w-rx, y+ry, rx, ry, -math.Pi/2, 0)
		r.LineTo(ctx.tx(x+w), ctx.ty(y+h-ry))
		arcQuarter(r, ctx, x+w-rx, y+h-ry, rx, ry, 0, math.Pi/2)
		r.LineTo(ctx.tx(x+rx), ctx.ty(y+h))
		arcQuarter(r, ctx, x+rx, y+h-ry, rx, ry, math.Pi/2, math.Pi)
		r.LineTo(ctx.tx(x), ctx.ty(y+ry))
		arcQuarter(r, ctx, x+rx, y+ry, rx, ry, math.Pi, 3*math.Pi/2)
	} else {
		r.MoveTo(ctx.tx(x), ctx.ty(y))
		r.LineTo(ctx.tx(x+w), ctx.ty(y))
		r.LineTo(ctx.tx(x+w), ctx.ty(y+h))
		r.LineTo(ctx.tx(x), ctx.ty(y+h))
	}
	r.ClosePath()
	r.Draw(ctx.img, ctx.img.Bounds(), image.NewUniform(fill), image.Point{})
}

func (ctx *svgRenderCtx) renderCircle(attrs map[string]string, fill color.RGBA) {
	cx := parseAttrFloat(attrs, "cx")
	cy := parseAttrFloat(attrs, "cy")
	radius := parseAttrFloat(attrs, "r")
	ctx.renderEllipseAt(cx, cy, radius, radius, fill)
}

func (ctx *svgRenderCtx) renderEllipse(attrs map[string]string, fill color.RGBA) {
	cx := parseAttrFloat(attrs, "cx")
	cy := parseAttrFloat(attrs, "cy")
	rx := parseAttrFloat(attrs, "rx")
	ry := parseAttrFloat(attrs, "ry")
	ctx.renderEllipseAt(cx, cy, rx, ry, fill)
}

func (ctx *svgRenderCtx) renderEllipseAt(cx, cy, rx, ry float64, fill color.RGBA) {
	r := vector.NewRasterizer(ctx.targetW, ctx.targetH)
	const n = 16
	for i := range n {
		angle := 2 * math.Pi * float64(i) / n
		px := cx + rx*math.Cos(angle)
		py := cy + ry*math.Sin(angle)
		if i == 0 {
			r.MoveTo(ctx.tx(px), ctx.ty(py))
		} else {
			a0 := 2 * math.Pi * float64(i-1) / n
			a1 := angle
			am := (a0 + a1) / 2
			k := 4.0 / 3.0 * math.Tan((a1-a0)/4)
			cp1x := cx + rx*(math.Cos(a0)-k*math.Sin(a0))
			cp1y := cy + ry*(math.Sin(a0)+k*math.Cos(a0))
			cp2x := cx + rx*(math.Cos(a1)+k*math.Sin(a1))
			cp2y := cy + ry*(math.Sin(a1)-k*math.Cos(a1))
			_ = am
			r.CubeTo(ctx.tx(cp1x), ctx.ty(cp1y), ctx.tx(cp2x), ctx.ty(cp2y), ctx.tx(px), ctx.ty(py))
		}
	}
	r.ClosePath()
	r.Draw(ctx.img, ctx.img.Bounds(), image.NewUniform(fill), image.Point{})
}

func (ctx *svgRenderCtx) renderPolygon(pts string, fill color.RGBA) {
	coords := parsePointList(pts)
	if len(coords) < 4 {
		return
	}
	r := vector.NewRasterizer(ctx.targetW, ctx.targetH)
	r.MoveTo(ctx.tx(coords[0]), ctx.ty(coords[1]))
	for i := 2; i < len(coords)-1; i += 2 {
		r.LineTo(ctx.tx(coords[i]), ctx.ty(coords[i+1]))
	}
	r.ClosePath()
	r.Draw(ctx.img, ctx.img.Bounds(), image.NewUniform(fill), image.Point{})
}

type pathCmd struct {
	cmd  byte
	args []float64
}

func parseSVGPath(d string) []pathCmd {
	var cmds []pathCmd
	i := 0
	for i < len(d) {
		for i < len(d) && (d[i] == ' ' || d[i] == '\t' || d[i] == '\n' || d[i] == '\r' || d[i] == ',') {
			i++
		}
		if i >= len(d) {
			break
		}
		c := d[i]
		if isPathCmd(c) {
			cmd := pathCmd{cmd: c}
			i++
			cmd.args, i = parsePathArgs(d, i)
			cmds = append(cmds, cmd)
		} else {
			if len(cmds) > 0 {
				args, ni := parsePathArgs(d, i)
				cmds[len(cmds)-1].args = append(cmds[len(cmds)-1].args, args...)
				i = ni
			} else {
				i++
			}
		}
	}
	return cmds
}

func isPathCmd(c byte) bool {
	switch c | 0x20 {
	case 'm', 'l', 'h', 'v', 'c', 's', 'q', 't', 'a', 'z':
		return true
	}
	return false
}

func parsePathArgs(d string, i int) ([]float64, int) {
	var args []float64
	for i < len(d) {
		for i < len(d) && (d[i] == ' ' || d[i] == '\t' || d[i] == '\n' || d[i] == '\r' || d[i] == ',') {
			i++
		}
		if i >= len(d) || isPathCmd(d[i]) {
			break
		}
		v, ni := parsePathNumber(d, i)
		if ni == i {
			break
		}
		args = append(args, v)
		i = ni
	}
	return args, i
}

func parsePathNumber(d string, i int) (float64, int) {
	start := i
	if i < len(d) && (d[i] == '-' || d[i] == '+') {
		i++
	}
	hasDot := false
	for i < len(d) && (d[i] >= '0' && d[i] <= '9' || d[i] == '.') {
		if d[i] == '.' {
			if hasDot {
				break
			}
			hasDot = true
		}
		i++
	}
	if i < len(d) && (d[i] == 'e' || d[i] == 'E') {
		i++
		if i < len(d) && (d[i] == '-' || d[i] == '+') {
			i++
		}
		for i < len(d) && d[i] >= '0' && d[i] <= '9' {
			i++
		}
	}
	if i == start {
		return 0, start
	}
	v, err := strconv.ParseFloat(d[start:i], 64)
	if err != nil {
		return 0, start
	}
	return v, i
}

func arcToCubic(r *vector.Rasterizer, ctx *svgRenderCtx, x1, y1, rxIn, ryIn, phi float64, largeArc, sweep bool, x2, y2 float64) {
	if (x1 == x2 && y1 == y2) || rxIn == 0 || ryIn == 0 {
		r.LineTo(ctx.tx(x2), ctx.ty(y2))
		return
	}

	rx := math.Abs(rxIn)
	ry := math.Abs(ryIn)
	sinPhi := math.Sin(phi * math.Pi / 180)
	cosPhi := math.Cos(phi * math.Pi / 180)

	dx := (x1 - x2) / 2
	dy := (y1 - y2) / 2
	x1p := cosPhi*dx + sinPhi*dy
	y1p := -sinPhi*dx + cosPhi*dy

	lambda := (x1p*x1p)/(rx*rx) + (y1p*y1p)/(ry*ry)
	if lambda > 1 {
		s := math.Sqrt(lambda)
		rx *= s
		ry *= s
	}

	num := rx*rx*ry*ry - rx*rx*y1p*y1p - ry*ry*x1p*x1p
	den := rx*rx*y1p*y1p + ry*ry*x1p*x1p
	if den == 0 {
		r.LineTo(ctx.tx(x2), ctx.ty(y2))
		return
	}
	sq := math.Max(0, num/den)
	root := math.Sqrt(sq)
	if largeArc == sweep {
		root = -root
	}
	cxp := root * rx * y1p / ry
	cyp := -root * ry * x1p / rx

	centerX := cosPhi*cxp - sinPhi*cyp + (x1+x2)/2
	centerY := sinPhi*cxp + cosPhi*cyp + (y1+y2)/2

	theta1 := vecAngle(1, 0, (x1p-cxp)/rx, (y1p-cyp)/ry)
	dTheta := vecAngle((x1p-cxp)/rx, (y1p-cyp)/ry, (-x1p-cxp)/rx, (-y1p-cyp)/ry)

	if !sweep && dTheta > 0 {
		dTheta -= 2 * math.Pi
	} else if sweep && dTheta < 0 {
		dTheta += 2 * math.Pi
	}

	segments := int(math.Ceil(math.Abs(dTheta) / (math.Pi / 2)))
	if segments == 0 {
		segments = 1
	}
	step := dTheta / float64(segments)

	for s := range segments {
		a1 := theta1 + float64(s)*step
		a2 := theta1 + float64(s+1)*step
		alpha := 4.0 / 3.0 * math.Tan((a2-a1)/4)

		cos1 := math.Cos(a1)
		sin1 := math.Sin(a1)
		cos2 := math.Cos(a2)
		sin2 := math.Sin(a2)

		ep1x := rx * cos1
		ep1y := ry * sin1
		ep2x := rx * cos2
		ep2y := ry * sin2

		cp1x := ep1x - alpha*rx*sin1
		cp1y := ep1y + alpha*ry*cos1
		cp2x := ep2x + alpha*rx*sin2
		cp2y := ep2y - alpha*ry*cos2

		p1x := cosPhi*cp1x - sinPhi*cp1y + centerX
		p1y := sinPhi*cp1x + cosPhi*cp1y + centerY
		p2x := cosPhi*cp2x - sinPhi*cp2y + centerX
		p2y := sinPhi*cp2x + cosPhi*cp2y + centerY
		px := cosPhi*ep2x - sinPhi*ep2y + centerX
		py := sinPhi*ep2x + cosPhi*ep2y + centerY

		r.CubeTo(ctx.tx(p1x), ctx.ty(p1y), ctx.tx(p2x), ctx.ty(p2y), ctx.tx(px), ctx.ty(py))
	}
}

func vecAngle(ux, uy, vx, vy float64) float64 {
	dot := ux*vx + uy*vy
	lenU := math.Sqrt(ux*ux + uy*uy)
	lenV := math.Sqrt(vx*vx + vy*vy)
	cos := dot / (lenU * lenV)
	cos = math.Max(-1, math.Min(1, cos))
	angle := math.Acos(cos)
	if ux*vy-uy*vx < 0 {
		angle = -angle
	}
	return angle
}

func arcQuarter(r *vector.Rasterizer, ctx *svgRenderCtx, cx, cy, rx, ry, a1, a2 float64) {
	alpha := 4.0 / 3.0 * math.Tan((a2-a1)/4)
	cos1, sin1 := math.Cos(a1), math.Sin(a1)
	cos2, sin2 := math.Cos(a2), math.Sin(a2)

	cp1x := cx + rx*(cos1-alpha*sin1)
	cp1y := cy + ry*(sin1+alpha*cos1)
	cp2x := cx + rx*(cos2+alpha*sin2)
	cp2y := cy + ry*(sin2-alpha*cos2)
	ex := cx + rx*cos2
	ey := cy + ry*sin2

	r.CubeTo(ctx.tx(cp1x), ctx.ty(cp1y), ctx.tx(cp2x), ctx.ty(cp2y), ctx.tx(ex), ctx.ty(ey))
}

func parseSVGColor(s string, fallback color.RGBA) color.RGBA {
	s = strings.TrimSpace(s)
	if s == "" || s == "none" || s == "transparent" {
		return color.RGBA{}
	}
	if s == "currentColor" {
		return fallback
	}

	if s[0] == '#' {
		return parseHexColor(s)
	}
	if strings.HasPrefix(s, "rgb") {
		return parseRGBFunc(s)
	}
	if c, ok := svgNamedColors[strings.ToLower(s)]; ok {
		return c
	}
	return fallback
}

func parseHexColor(s string) color.RGBA {
	s = strings.TrimPrefix(s, "#")
	switch len(s) {
	case 3:
		r, _ := strconv.ParseUint(string(s[0])+string(s[0]), 16, 8)
		g, _ := strconv.ParseUint(string(s[1])+string(s[1]), 16, 8)
		b, _ := strconv.ParseUint(string(s[2])+string(s[2]), 16, 8)
		return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	case 4:
		r, _ := strconv.ParseUint(string(s[0])+string(s[0]), 16, 8)
		g, _ := strconv.ParseUint(string(s[1])+string(s[1]), 16, 8)
		b, _ := strconv.ParseUint(string(s[2])+string(s[2]), 16, 8)
		a, _ := strconv.ParseUint(string(s[3])+string(s[3]), 16, 8)
		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	case 6:
		r, _ := strconv.ParseUint(s[0:2], 16, 8)
		g, _ := strconv.ParseUint(s[2:4], 16, 8)
		b, _ := strconv.ParseUint(s[4:6], 16, 8)
		return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	case 8:
		r, _ := strconv.ParseUint(s[0:2], 16, 8)
		g, _ := strconv.ParseUint(s[2:4], 16, 8)
		b, _ := strconv.ParseUint(s[4:6], 16, 8)
		a, _ := strconv.ParseUint(s[6:8], 16, 8)
		return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	}
	return color.RGBA{0, 0, 0, 255}
}

func parseRGBFunc(s string) color.RGBA {
	s = strings.TrimSpace(s)
	hasAlpha := strings.HasPrefix(s, "rgba(")
	if hasAlpha {
		s = strings.TrimPrefix(s, "rgba(")
	} else {
		s = strings.TrimPrefix(s, "rgb(")
	}
	s = strings.TrimSuffix(s, ")")
	parts := strings.FieldsFunc(s, func(r rune) bool { return r == ',' || r == '/' || unicode.IsSpace(r) })
	if len(parts) < 3 {
		return color.RGBA{0, 0, 0, 255}
	}
	r, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	g, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	b, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
	a := 255
	if len(parts) >= 4 {
		af, _ := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
		if af <= 1 {
			a = int(af * 255)
		} else {
			a = int(af)
		}
	}
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func parsePointList(s string) []float64 {
	s = strings.ReplaceAll(s, ",", " ")
	parts := strings.Fields(s)
	var coords []float64
	for _, p := range parts {
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			continue
		}
		coords = append(coords, v)
	}
	return coords
}

func parseDimension(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "px")
	s = strings.TrimSuffix(s, "pt")
	s = strings.TrimSuffix(s, "em")
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func attrMap(attrs []xml.Attr) map[string]string {
	m := make(map[string]string, len(attrs))
	for _, a := range attrs {
		m[a.Name.Local] = a.Value
	}
	return m
}

func parseAttrFloat(attrs map[string]string, key string) float64 {
	s, ok := attrs[key]
	if !ok {
		return 0
	}
	return parseDimension(s)
}

var svgNamedColors = map[string]color.RGBA{
	"black":   {0, 0, 0, 255},
	"white":   {255, 255, 255, 255},
	"red":     {255, 0, 0, 255},
	"green":   {0, 128, 0, 255},
	"blue":    {0, 0, 255, 255},
	"yellow":  {255, 255, 0, 255},
	"cyan":    {0, 255, 255, 255},
	"magenta": {255, 0, 255, 255},
	"gray":    {128, 128, 128, 255},
	"grey":    {128, 128, 128, 255},
	"orange":  {255, 165, 0, 255},
	"purple":  {128, 0, 128, 255},
	"pink":    {255, 192, 203, 255},
	"brown":   {165, 42, 42, 255},
	"silver":  {192, 192, 192, 255},
	"gold":    {255, 215, 0, 255},
	"navy":    {0, 0, 128, 255},
	"teal":    {0, 128, 128, 255},
	"maroon":  {128, 0, 0, 255},
	"olive":   {128, 128, 0, 255},
	"lime":    {0, 255, 0, 255},
	"aqua":    {0, 255, 255, 255},
	"fuchsia": {255, 0, 255, 255},
}
