package render

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"strings"
	"sync"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/vector"

	fontpkg "github.com/macawls/ogre/font"
	"github.com/macawls/ogre/layout"
	"github.com/macawls/ogre/parse"
	"github.com/macawls/ogre/style"
)

var rgbaPool sync.Pool

func acquireRGBA(r image.Rectangle) *image.RGBA {
	if v := rgbaPool.Get(); v != nil {
		img := v.(*image.RGBA)
		need := r.Dx() * 4 * r.Dy()
		if cap(img.Pix) >= need {
			img.Pix = img.Pix[:need]
			img.Stride = r.Dx() * 4
			img.Rect = r
			clear(img.Pix)
			return img
		}
	}
	return image.NewRGBA(r)
}

func releaseRGBA(img *image.RGBA) {
	rgbaPool.Put(img)
}

type PNGRenderer struct {
	img           *image.RGBA
	styles        map[*parse.Node]*style.ComputedStyle
	fonts         *fontpkg.Manager
	reverse       map[*layout.Node]*parse.Node
	wrappedText   map[*parse.Node][]fontpkg.TextLine
	emojiProvider *fontpkg.EmojiProvider
	maskCache     map[maskKey]*image.Alpha
}

type maskKey struct {
	w, h             int
	tl, tr, br, bl   int
}

// RenderPNG generates the corresponding output format.
type PNGOptions struct {
	WrappedText   map[*parse.Node][]fontpkg.TextLine
	EmojiProvider *fontpkg.EmojiProvider
}

func RenderPNG(tree *layout.LayoutTree, styles map[*parse.Node]*style.ComputedStyle, fonts *fontpkg.Manager, width, height int, opts ...PNGOptions) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	reverse := make(map[*layout.Node]*parse.Node, len(tree.NodeMap))
	for pn, ln := range tree.NodeMap {
		reverse[ln] = pn
	}

	var o PNGOptions
	if len(opts) > 0 {
		o = opts[0]
	}

	r := &PNGRenderer{
		img:           img,
		styles:        styles,
		fonts:         fonts,
		reverse:       reverse,
		wrappedText:   o.WrappedText,
		emojiProvider: o.EmojiProvider,
	}

	if tree.Root != nil {
		pn := reverse[tree.Root]
		cs := styles[pn]
		r.renderNode(tree.Root, pn, cs, 0, 0)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (r *PNGRenderer) renderNode(node *layout.Node, pn *parse.Node, cs *style.ComputedStyle, parentX, parentY float64) {
	if cs == nil {
		cs = style.NewComputedStyle()
	}

	opacity := cs.Opacity
	if opacity == 0 {
		opacity = 1
	}
	if opacity < 0 {
		return
	}

	l := node.Layout
	absX := parentX + l.X
	absY := parentY + l.Y

	if opacity < 1 {
		tmp := acquireRGBA(r.img.Bounds())
		defer releaseRGBA(tmp)
		sub := &PNGRenderer{
			img:           tmp,
			styles:        r.styles,
			fonts:         r.fonts,
			reverse:       r.reverse,
			wrappedText:   r.wrappedText,
			emojiProvider: r.emojiProvider,
		}
		sub.renderNodeContent(node, pn, cs, absX, absY)
		bounds := image.Rect(int(absX), int(absY), int(absX+l.Width), int(absY+l.Height)).Intersect(r.img.Bounds())
		for py := bounds.Min.Y; py < bounds.Max.Y; py++ {
			for px := bounds.Min.X; px < bounds.Max.X; px++ {
				off := tmp.PixOffset(px, py)
				sa := tmp.Pix[off+3]
				if sa == 0 {
					continue
				}
				na := uint8(float64(sa) * opacity)
				src := color.RGBA{R: tmp.Pix[off], G: tmp.Pix[off+1], B: tmp.Pix[off+2], A: na}
				doff := r.img.PixOffset(px, py)
				dr, dg, db, da := blendOver(src.R, src.G, src.B, src.A, r.img.Pix[doff], r.img.Pix[doff+1], r.img.Pix[doff+2], r.img.Pix[doff+3])
				r.img.Pix[doff] = dr
				r.img.Pix[doff+1] = dg
				r.img.Pix[doff+2] = db
				r.img.Pix[doff+3] = da
			}
		}
		return
	}

	r.renderNodeContent(node, pn, cs, absX, absY)
}

func (r *PNGRenderer) renderNodeContent(node *layout.Node, pn *parse.Node, cs *style.ComputedStyle, absX, absY float64) {
	l := node.Layout

	if pn != nil && pn.Type == parse.TextNode {
		r.renderTextNode(l, pn, cs, absX, absY)
		return
	}

	if pn != nil && pn.Tag == "img" {
		if src := pn.Attrs["src"]; src != "" {
			r.renderImage(src, cs, absX, absY, l.Width, l.Height)
		}
		return
	}

	if pn != nil && pn.Tag == "svg" {
		r.renderInlineSVG(pn, cs, absX, absY, l.Width, l.Height)
		return
	}

	r.renderBoxShadows(l, cs, absX, absY, false)

	hasRadius := cs.BorderTopLeftRadius > 0 || cs.BorderTopRightRadius > 0 ||
		cs.BorderBottomLeftRadius > 0 || cs.BorderBottomRightRadius > 0

	if cs.BackgroundImage != "" {
		if hasRadius {
			tmp := acquireRGBA(r.img.Bounds())
			sub := &PNGRenderer{img: tmp, styles: r.styles, fonts: r.fonts, reverse: r.reverse, wrappedText: r.wrappedText, emojiProvider: r.emojiProvider}
			sub.renderGradient(absX, absY, l.Width, l.Height, cs)
			rect := image.Rect(int(absX), int(absY), int(absX+l.Width), int(absY+l.Height)).Intersect(r.img.Bounds())
			mask := r.cachedMask(int(l.Width), int(l.Height), cs.BorderTopLeftRadius, cs.BorderTopRightRadius, cs.BorderBottomRightRadius, cs.BorderBottomLeftRadius)
			draw.DrawMask(r.img, rect, tmp, rect.Min, mask, image.Point{}, draw.Over)
			releaseRGBA(tmp)
		} else {
			r.renderGradient(absX, absY, l.Width, l.Height, cs)
		}
	} else if !cs.BackgroundColor.IsTransparent() {
		c := styleToColor(cs.BackgroundColor)
		rect := image.Rect(int(absX), int(absY), int(absX+l.Width), int(absY+l.Height)).Intersect(r.img.Bounds())
		if hasRadius {
			mask := r.cachedMask(int(l.Width), int(l.Height), cs.BorderTopLeftRadius, cs.BorderTopRightRadius, cs.BorderBottomRightRadius, cs.BorderBottomLeftRadius)
			draw.DrawMask(r.img, rect, uniformSrc(c), image.Point{}, mask, image.Point{}, draw.Over)
		} else {
			draw.Draw(r.img, rect, uniformSrc(c), image.Point{}, draw.Over)
		}
	}

	r.renderBorders(absX, absY, l.Width, l.Height, cs)

	r.renderBoxShadows(l, cs, absX, absY, true)

	if cs.Overflow == style.OverflowHidden {
		clip := image.Rect(int(absX), int(absY), int(absX+l.Width), int(absY+l.Height))
		tmp := acquireRGBA(r.img.Bounds())
		sub := &PNGRenderer{img: tmp, styles: r.styles, fonts: r.fonts, reverse: r.reverse, wrappedText: r.wrappedText, emojiProvider: r.emojiProvider}
		for _, child := range node.Children {
			cpn := sub.reverse[child]
			ccs := sub.styles[cpn]
			sub.renderNode(child, cpn, ccs, absX, absY)
		}
		hasRadius := cs.BorderTopLeftRadius > 0 || cs.BorderTopRightRadius > 0 ||
			cs.BorderBottomLeftRadius > 0 || cs.BorderBottomRightRadius > 0
		if hasRadius {
			mask := r.cachedMask(int(l.Width), int(l.Height),
				cs.BorderTopLeftRadius, cs.BorderTopRightRadius,
				cs.BorderBottomRightRadius, cs.BorderBottomLeftRadius)
			draw.DrawMask(r.img, clip, tmp, clip.Min, mask, image.Point{}, draw.Over)
		} else {
			draw.Draw(r.img, clip, tmp, clip.Min, draw.Over)
		}
		releaseRGBA(tmp)
		return
	}

	for _, child := range node.Children {
		cpn := r.reverse[child]
		ccs := r.styles[cpn]
		r.renderNode(child, cpn, ccs, absX, absY)
	}
}

func (r *PNGRenderer) renderBorders(absX, absY, w, h float64, cs *style.ComputedStyle) {
	x := int(absX)
	y := int(absY)
	wi := int(w)
	hi := int(h)

	hasBorder := (cs.BorderTopWidth > 0 && cs.BorderTopStyle != style.BorderStyleNone) ||
		(cs.BorderBottomWidth > 0 && cs.BorderBottomStyle != style.BorderStyleNone) ||
		(cs.BorderLeftWidth > 0 && cs.BorderLeftStyle != style.BorderStyleNone) ||
		(cs.BorderRightWidth > 0 && cs.BorderRightStyle != style.BorderStyleNone)
	if !hasBorder {
		return
	}

	hasRadius := cs.BorderTopLeftRadius > 0 || cs.BorderTopRightRadius > 0 ||
		cs.BorderBottomLeftRadius > 0 || cs.BorderBottomRightRadius > 0

	if hasRadius {
		tmp := acquireRGBA(r.img.Bounds())
		if cs.BorderTopWidth > 0 && cs.BorderTopStyle != style.BorderStyleNone {
			c := styleToColor(cs.BorderTopColor)
			fillRect(tmp, x, y, wi, int(math.Max(1, cs.BorderTopWidth)), c)
		}
		if cs.BorderBottomWidth > 0 && cs.BorderBottomStyle != style.BorderStyleNone {
			c := styleToColor(cs.BorderBottomColor)
			by := y + hi - int(math.Max(1, cs.BorderBottomWidth))
			fillRect(tmp, x, by, wi, int(math.Max(1, cs.BorderBottomWidth)), c)
		}
		if cs.BorderLeftWidth > 0 && cs.BorderLeftStyle != style.BorderStyleNone {
			c := styleToColor(cs.BorderLeftColor)
			fillRect(tmp, x, y, int(math.Max(1, cs.BorderLeftWidth)), hi, c)
		}
		if cs.BorderRightWidth > 0 && cs.BorderRightStyle != style.BorderStyleNone {
			c := styleToColor(cs.BorderRightColor)
			bx := x + wi - int(math.Max(1, cs.BorderRightWidth))
			fillRect(tmp, bx, y, int(math.Max(1, cs.BorderRightWidth)), hi, c)
		}
		rect := image.Rect(x, y, x+wi, y+hi).Intersect(r.img.Bounds())
		mask := r.cachedMask(wi, hi, cs.BorderTopLeftRadius, cs.BorderTopRightRadius, cs.BorderBottomRightRadius, cs.BorderBottomLeftRadius)
		draw.DrawMask(r.img, rect, tmp, rect.Min, mask, image.Point{}, draw.Over)
		releaseRGBA(tmp)
		return
	}

	if cs.BorderTopWidth > 0 && cs.BorderTopStyle != style.BorderStyleNone {
		c := styleToColor(cs.BorderTopColor)
		fillRect(r.img, x, y, wi, int(math.Max(1, cs.BorderTopWidth)), c)
	}
	if cs.BorderBottomWidth > 0 && cs.BorderBottomStyle != style.BorderStyleNone {
		c := styleToColor(cs.BorderBottomColor)
		by := y + hi - int(math.Max(1, cs.BorderBottomWidth))
		fillRect(r.img, x, by, wi, int(math.Max(1, cs.BorderBottomWidth)), c)
	}
	if cs.BorderLeftWidth > 0 && cs.BorderLeftStyle != style.BorderStyleNone {
		c := styleToColor(cs.BorderLeftColor)
		fillRect(r.img, x, y, int(math.Max(1, cs.BorderLeftWidth)), hi, c)
	}
	if cs.BorderRightWidth > 0 && cs.BorderRightStyle != style.BorderStyleNone {
		c := styleToColor(cs.BorderRightColor)
		bx := x + wi - int(math.Max(1, cs.BorderRightWidth))
		fillRect(r.img, bx, y, int(math.Max(1, cs.BorderRightWidth)), hi, c)
	}
}

func (r *PNGRenderer) renderTextNode(l layout.Layout, pn *parse.Node, cs *style.ComputedStyle, absX, absY float64) {
	family := cs.FontFamily
	if family == "" {
		family = "default"
	}
	size := cs.FontSize
	if size == 0 {
		size = 16
	}
	weight := cs.FontWeight
	if weight == 0 {
		weight = 400
	}
	fstyle := cs.FontStyle
	if fstyle == "" {
		fstyle = "normal"
	}

	if r.fonts == nil {
		return
	}

	face := r.fonts.Resolve(family, weight, fstyle)
	if face == nil {
		return
	}

	ff, err := r.fonts.NewFace(face, size)
	if err != nil {
		return
	}

	tc := styleToColor(cs.Color)
	if cs.Color.A == -1 {
		tc = color.RGBA{0, 0, 0, 255}
	}

	lineHeight := cs.LineHeight
	if lineHeight == 0 {
		lineHeight = size * 1.2
	}
	ascent := fontpkg.Ascent(ff)

	if lines, ok := r.wrappedText[pn]; ok && len(lines) > 0 {
		for i, line := range lines {
			text := applyTextTransform(line.Text, cs.TextTransform)
			x := alignX(absX, l.Width, line.Width, cs.TextAlign)
			y := absY + ascent + float64(i)*lineHeight
			r.drawTextWithEmoji(text, x, y, ascent, size, tc, ff, cs)
		}
		return
	}

	r.drawTextWithEmoji(pn.Text, absX, absY+ascent, ascent, size, tc, ff, cs)
}

func (r *PNGRenderer) drawShapedText(text string, x, y, size float64, tc color.RGBA, cs *style.ComputedStyle) bool {
	if r.fonts == nil {
		return false
	}
	family := cs.FontFamily
	if family == "" {
		family = "default"
	}
	fstyle := cs.FontStyle
	if fstyle == "" {
		fstyle = "normal"
	}
	rtl := cs.Direction == "rtl"
	if !rtl && !needsShaping(text) {
		return false
	}
	pathD, _ := fontpkg.ShapedTextToPath(r.fonts, text, family, cs.FontWeight, fstyle, size, rtl)
	if pathD == "" {
		return false
	}
	rast := vector.NewRasterizer(r.img.Bounds().Dx(), r.img.Bounds().Dy())
	cmds := parseSVGPath(pathD)
	var cx, cy, startX, startY, lastCPX, lastCPY float64
	var lastCmd byte
	for _, cmd := range cmds {
		switch cmd.cmd {
		case 'M':
			for i := 0; i < len(cmd.args)-1; i += 2 {
				cx, cy = cmd.args[i]+x, cmd.args[i+1]+y
				if i == 0 {
					startX, startY = cx, cy
					rast.MoveTo(float32(cx), float32(cy))
				} else {
					rast.LineTo(float32(cx), float32(cy))
				}
			}
		case 'L':
			for i := 0; i < len(cmd.args)-1; i += 2 {
				cx, cy = cmd.args[i]+x, cmd.args[i+1]+y
				rast.LineTo(float32(cx), float32(cy))
			}
		case 'Q':
			for i := 0; i < len(cmd.args)-3; i += 4 {
				x1, y1 := cmd.args[i]+x, cmd.args[i+1]+y
				cx, cy = cmd.args[i+2]+x, cmd.args[i+3]+y
				rast.QuadTo(float32(x1), float32(y1), float32(cx), float32(cy))
			}
		case 'C':
			for i := 0; i < len(cmd.args)-5; i += 6 {
				x1, y1 := cmd.args[i]+x, cmd.args[i+1]+y
				x2, y2 := cmd.args[i+2]+x, cmd.args[i+3]+y
				cx, cy = cmd.args[i+4]+x, cmd.args[i+5]+y
				lastCPX, lastCPY = x2, y2
				rast.CubeTo(float32(x1), float32(y1), float32(x2), float32(y2), float32(cx), float32(cy))
			}
		case 'Z', 'z':
			rast.ClosePath()
			cx, cy = startX, startY
		}
		lastCmd = cmd.cmd
	}
	_, _ = lastCPX, lastCPY
	_ = lastCmd
	if lastCmd != 'Z' && lastCmd != 'z' && lastCmd != 0 {
		rast.ClosePath()
	}
	rast.Draw(r.img, r.img.Bounds(), image.NewUniform(tc), image.Point{})
	return true
}

func (r *PNGRenderer) drawTextWithEmoji(text string, x, y, ascent, size float64, tc color.RGBA, ff font.Face, cs *style.ComputedStyle) {
	if r.emojiProvider == nil || !containsEmoji(text) {
		if r.drawShapedText(text, x, y, size, tc, cs) {
			return
		}
		drawer := &font.Drawer{
			Dst:  r.img,
			Src:  image.NewUniform(tc),
			Face: ff,
			Dot:  fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)},
		}
		drawer.DrawString(text)
		return
	}

	segments := fontpkg.SplitEmoji(text)
	cx := x
	m := fontpkg.NewMeasurer(ff, cs.LetterSpacing)

	for _, seg := range segments {
		if seg.IsEmoji {
			emojiImg, err := r.emojiProvider.FetchPNG(seg.Text)
			if err == nil && emojiImg != nil {
				ey := int(y - ascent)
				ex := int(cx)
				es := int(size)
				dst := image.Rect(ex, ey, ex+es, ey+es)
				xdraw.BiLinear.Scale(r.img, dst, emojiImg, emojiImg.Bounds(), draw.Over, nil)
				cx += size
			} else {
				cx += size
			}
		} else {
			drawer := &font.Drawer{
				Dst:  r.img,
				Src:  image.NewUniform(tc),
				Face: ff,
				Dot:  fixed.Point26_6{X: fixed.Int26_6(cx * 64), Y: fixed.Int26_6(y * 64)},
			}
			drawer.DrawString(seg.Text)
			cx += m.StringWidth(seg.Text)
		}
	}
}

func (r *PNGRenderer) renderGradient(absX, absY, w, h float64, cs *style.ComputedStyle) {
	g, err := style.ParseGradient(cs.BackgroundImage)
	if err != nil {
		if !cs.BackgroundColor.IsTransparent() {
			c := styleToColor(cs.BackgroundColor)
			fillRect(r.img, int(absX), int(absY), int(w), int(h), c)
		}
		return
	}

	distributeStops(g.Stops)

	x0, y0 := int(absX), int(absY)
	wi, hi := int(w), int(h)

	switch g.Type {
	case style.LinearGradient, style.RepeatingLinearGradient:
		r.renderLinearGradientPNG(g, x0, y0, wi, hi)
	case style.RadialGradient, style.RepeatingRadialGradient:
		r.renderRadialGradientPNG(g, x0, y0, wi, hi)
	}
}

func (r *PNGRenderer) renderLinearGradientPNG(g style.Gradient, rx, ry, rw, rh int) {
	rad := g.Angle * math.Pi / 180
	sinA := math.Sin(rad)
	cosA := math.Cos(rad)

	cx := float64(rw) / 2
	cy := float64(rh) / 2
	length := math.Abs(float64(rw)*sinA) + math.Abs(float64(rh)*cosA)
	if length == 0 {
		length = 1
	}

	linStops := toLinearStops(g.Stops)

	stripLen := int(math.Ceil(length)) + 1
	if stripLen < 2 {
		stripLen = 2
	}
	strip := buildGradientStrip(linStops, stripLen)

	bounds := r.img.Bounds()
	invLength := 1.0 / length
	for py := ry; py < ry+rh; py++ {
		if py < bounds.Min.Y || py >= bounds.Max.Y {
			continue
		}
		dy := float64(py-ry) - cy
		rowBase := -dy*cosA*invLength + 0.5
		rowStep := sinA * invLength
		t := rowBase + (float64(0)-cx)*rowStep
		rowOff := py * r.img.Stride
		for px := rx; px < rx+rw; px++ {
			if px >= bounds.Min.X && px < bounds.Max.X {
				tc := t
				if g.Repeating {
					tc = tc - math.Floor(tc)
				} else if tc < 0 {
					tc = 0
				} else if tc > 1 {
					tc = 1
				}
				idx := int(tc * float64(stripLen-1))
				if idx >= stripLen {
					idx = stripLen - 1
				}
				s := strip[idx]
				dither := bayerMatrix[py&3][px&3] - 0.5
				off := rowOff + px*4
				r.img.Pix[off] = uint8(clampF(math.Floor(s.sr*255+dither+0.5), 0, 255))
				r.img.Pix[off+1] = uint8(clampF(math.Floor(s.sg*255+dither+0.5), 0, 255))
				r.img.Pix[off+2] = uint8(clampF(math.Floor(s.sb*255+dither+0.5), 0, 255))
				r.img.Pix[off+3] = s.a
			}
			t += rowStep
		}
	}
}

type stripEntry struct {
	sr, sg, sb float64
	a          uint8
}

func buildGradientStrip(stops []linearStop, n int) []stripEntry {
	strip := make([]stripEntry, n)
	for i := range n {
		t := float64(i) / float64(n-1)
		lr, lg, lb, la := interpolateLinear(stops, t)
		strip[i] = stripEntry{
			sr: linearToSrgbF(clampF(lr, 0, 1)),
			sg: linearToSrgbF(clampF(lg, 0, 1)),
			sb: linearToSrgbF(clampF(lb, 0, 1)),
			a:  uint8(math.Round(la * 255)),
		}
	}
	return strip
}

func interpolateLinear(stops []linearStop, t float64) (r, g, b, a float64) {
	if len(stops) == 0 {
		return 0, 0, 0, 1
	}
	if t <= stops[0].Position {
		return stops[0].R, stops[0].G, stops[0].B, stops[0].A
	}
	last := stops[len(stops)-1]
	if t >= last.Position {
		return last.R, last.G, last.B, last.A
	}
	for i := 1; i < len(stops); i++ {
		if t <= stops[i].Position {
			prev := stops[i-1]
			curr := stops[i]
			span := curr.Position - prev.Position
			if span <= 0 {
				return curr.R, curr.G, curr.B, curr.A
			}
			f := (t - prev.Position) / span
			return prev.R + f*(curr.R-prev.R),
				prev.G + f*(curr.G-prev.G),
				prev.B + f*(curr.B-prev.B),
				prev.A + f*(curr.A-prev.A)
		}
	}
	return last.R, last.G, last.B, last.A
}

func clampF(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func (r *PNGRenderer) renderRadialGradientPNG(g style.Gradient, rx, ry, rw, rh int) {
	cx := float64(rw) * g.PositionX / 100
	cy := float64(rh) * g.PositionY / 100

	maxDist := 0.0
	for _, corner := range [][2]float64{{0, 0}, {float64(rw), 0}, {0, float64(rh)}, {float64(rw), float64(rh)}} {
		dx := corner[0] - cx
		dy := corner[1] - cy
		d := math.Sqrt(dx*dx + dy*dy)
		if d > maxDist {
			maxDist = d
		}
	}
	if maxDist == 0 {
		maxDist = 1
	}

	linStops := toLinearStops(g.Stops)

	stripLen := int(math.Ceil(maxDist)) + 1
	if stripLen < 2 {
		stripLen = 2
	}
	strip := buildGradientStrip(linStops, stripLen)

	bounds := r.img.Bounds()
	invMax := 1.0 / maxDist
	for py := ry; py < ry+rh; py++ {
		if py < bounds.Min.Y || py >= bounds.Max.Y {
			continue
		}
		dy := float64(py-ry) - cy
		dy2 := dy * dy
		rowOff := py * r.img.Stride
		for px := rx; px < rx+rw; px++ {
			if px < bounds.Min.X || px >= bounds.Max.X {
				continue
			}
			dx := float64(px-rx) - cx
			dist := math.Sqrt(dx*dx + dy2)
			t := dist * invMax
			if g.Repeating {
				t = t - math.Floor(t)
			} else if t > 1 {
				t = 1
			}
			idx := int(t * float64(stripLen-1))
			if idx >= stripLen {
				idx = stripLen - 1
			}
			s := strip[idx]
			dither := bayerMatrix[py&3][px&3] - 0.5
			off := rowOff + px*4
			r.img.Pix[off] = uint8(clampF(math.Floor(s.sr*255+dither+0.5), 0, 255))
			r.img.Pix[off+1] = uint8(clampF(math.Floor(s.sg*255+dither+0.5), 0, 255))
			r.img.Pix[off+2] = uint8(clampF(math.Floor(s.sb*255+dither+0.5), 0, 255))
			r.img.Pix[off+3] = s.a
		}
	}
}

type linearStop struct {
	R, G, B, A float64
	Position   float64
}

var srgbDecodeTable [256]float64
var srgbEncodeTable [4096]float64

func init() {
	for i := range 256 {
		c := float64(i) / 255
		if c <= 0.04045 {
			srgbDecodeTable[i] = c / 12.92
		} else {
			srgbDecodeTable[i] = math.Pow((c+0.055)/1.055, 2.4)
		}
	}
	for i := range 4096 {
		v := float64(i) / 4095
		if v <= 0.0031308 {
			srgbEncodeTable[i] = v * 12.92
		} else {
			srgbEncodeTable[i] = 1.055*math.Pow(v, 1.0/2.4) - 0.055
		}
	}
}

func srgbToLinear(v uint8) float64 {
	return srgbDecodeTable[v]
}

func linearToSrgbF(v float64) float64 {
	if v <= 0 {
		return 0
	}
	if v >= 1 {
		return 1
	}
	return srgbEncodeTable[int(v*4095+0.5)]
}


func toLinearStops(stops []style.ColorStop) []linearStop {
	out := make([]linearStop, len(stops))
	for i, s := range stops {
		out[i] = linearStop{
			R: srgbToLinear(s.Color.R), G: srgbToLinear(s.Color.G),
			B: srgbToLinear(s.Color.B), A: s.Color.A,
			Position: s.Position,
		}
	}
	return out
}

var bayerMatrix = [4][4]float64{
	{0.0 / 16, 8.0 / 16, 2.0 / 16, 10.0 / 16},
	{12.0 / 16, 4.0 / 16, 14.0 / 16, 6.0 / 16},
	{3.0 / 16, 11.0 / 16, 1.0 / 16, 9.0 / 16},
	{15.0 / 16, 7.0 / 16, 13.0 / 16, 5.0 / 16},
}


func fillRect(img *image.RGBA, x, y, w, h int, c color.Color) {
	rect := image.Rect(x, y, x+w, y+h).Intersect(img.Bounds())
	draw.Draw(img, rect, uniformSrc(c), image.Point{}, draw.Over)
}

func blendOver(sr, sg, sb, sa, dr, dg, db, da uint8) (uint8, uint8, uint8, uint8) {
	if sa == 255 {
		return sr, sg, sb, 255
	}
	if sa == 0 {
		return dr, dg, db, da
	}
	a := uint32(sa)
	ia := 255 - a
	oa := a + uint32(da)*ia/255
	if oa == 0 {
		return 0, 0, 0, 0
	}
	rr := (uint32(sr)*a + uint32(dr)*uint32(da)*ia/255) / oa
	gg := (uint32(sg)*a + uint32(dg)*uint32(da)*ia/255) / oa
	bb := (uint32(sb)*a + uint32(db)*uint32(da)*ia/255) / oa
	return uint8(rr), uint8(gg), uint8(bb), uint8(oa)
}

func styleToColor(c style.Color) color.RGBA {
	a := uint8(math.Round(c.A * 255))
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: a}
}

func uniformSrc(c color.Color) *image.Uniform {
	if rgba, ok := c.(color.RGBA); ok {
		return image.NewUniform(color.NRGBA{R: rgba.R, G: rgba.G, B: rgba.B, A: rgba.A})
	}
	return image.NewUniform(c)
}

func (r *PNGRenderer) renderBoxShadows(l layout.Layout, cs *style.ComputedStyle, absX, absY float64, insetPass bool) {
	if cs.BoxShadow == "" {
		return
	}
	shadows, err := style.ParseBoxShadow(cs.BoxShadow)
	if err != nil || len(shadows) == 0 {
		return
	}
	for _, s := range shadows {
		if s.Inset != insetPass {
			continue
		}
		sc := styleToColor(s.Color)
		if s.Inset {
			r.renderInsetShadow(absX, absY, l.Width, l.Height, s, sc)
		} else {
			r.renderOutsetShadow(absX, absY, l.Width, l.Height, s, sc)
		}
	}
}

func (r *PNGRenderer) renderOutsetShadow(absX, absY, w, h float64, s style.Shadow, sc color.RGBA) {
	blur := int(math.Ceil(s.Blur))
	if blur < 1 {
		blur = 1
	}
	pad := blur * 3

	sx := int(absX + s.OffsetX - s.Spread)
	sy := int(absY + s.OffsetY - s.Spread)
	sw := int(w + 2*s.Spread)
	sh := int(h + 2*s.Spread)

	tw := sw + 2*pad
	th := sh + 2*pad
	if tw <= 0 || th <= 0 {
		return
	}

	alpha := image.NewAlpha(image.Rect(0, 0, tw, th))
	for py := pad; py < pad+sh; py++ {
		for px := pad; px < pad+sw; px++ {
			alpha.SetAlpha(px, py, color.Alpha{A: sc.A})
		}
	}

	blurred := boxBlurAlpha(alpha, blur)

	ox := sx - pad
	oy := sy - pad
	bounds := r.img.Bounds()
	for py := 0; py < th; py++ {
		dy := oy + py
		if dy < bounds.Min.Y || dy >= bounds.Max.Y {
			continue
		}
		for px := 0; px < tw; px++ {
			dx := ox + px
			if dx < bounds.Min.X || dx >= bounds.Max.X {
				continue
			}
			a := blurred.AlphaAt(px, py).A
			if a == 0 {
				continue
			}
			src := color.RGBA{R: sc.R, G: sc.G, B: sc.B, A: a}
			doff := r.img.PixOffset(dx, dy)
			dr, dg, db, da := blendOver(src.R, src.G, src.B, src.A,
				r.img.Pix[doff], r.img.Pix[doff+1], r.img.Pix[doff+2], r.img.Pix[doff+3])
			r.img.Pix[doff] = dr
			r.img.Pix[doff+1] = dg
			r.img.Pix[doff+2] = db
			r.img.Pix[doff+3] = da
		}
	}
}

func boxBlurAlpha(src *image.Alpha, radius int) *image.Alpha {
	if radius <= 0 {
		return src
	}
	b := src.Bounds()
	tmp := image.NewAlpha(b)
	dst := image.NewAlpha(b)
	boxBlurH(src, tmp, b, radius)
	boxBlurV(tmp, dst, b, radius)
	tmp2 := image.NewAlpha(b)
	boxBlurH(dst, tmp2, b, radius)
	dst2 := image.NewAlpha(b)
	boxBlurV(tmp2, dst2, b, radius)
	return dst2
}

func boxBlurH(src, dst *image.Alpha, b image.Rectangle, r int) {
	h := b.Dy()
	if h > 64 {
		workers := 4
		chunk := (h + workers - 1) / workers
		var wg sync.WaitGroup
		for i := range workers {
			y0 := b.Min.Y + i*chunk
			y1 := y0 + chunk
			if y1 > b.Max.Y {
				y1 = b.Max.Y
			}
			if y0 >= y1 {
				break
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				blurHRows(src, dst, b, r, y0, y1)
			}()
		}
		wg.Wait()
		return
	}
	blurHRows(src, dst, b, r, b.Min.Y, b.Max.Y)
}

func blurHRows(src, dst *image.Alpha, b image.Rectangle, r, y0, y1 int) {
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

func boxBlurV(src, dst *image.Alpha, b image.Rectangle, r int) {
	w := b.Dx()
	if w > 64 {
		workers := 4
		chunk := (w + workers - 1) / workers
		var wg sync.WaitGroup
		for i := range workers {
			x0 := b.Min.X + i*chunk
			x1 := x0 + chunk
			if x1 > b.Max.X {
				x1 = b.Max.X
			}
			if x0 >= x1 {
				break
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				blurVCols(src, dst, b, r, x0, x1)
			}()
		}
		wg.Wait()
		return
	}
	blurVCols(src, dst, b, r, b.Min.X, b.Max.X)
}

func blurVCols(src, dst *image.Alpha, b image.Rectangle, r, x0, x1 int) {
	h := b.Dy()
	div := float64(2*r + 1)
	stride := dst.Stride
	for x := x0; x < x1; x++ {
		sum := 0.0
		for y := -r; y <= r; y++ {
			cy := clampInt(y, 0, h-1)
			sum += float64(src.Pix[cy*stride+x])
		}
		for y := b.Min.Y; y < b.Max.Y; y++ {
			dst.Pix[y*stride+x] = uint8(math.Round(sum / div))
			ny := clampInt(y+r+1, 0, h-1)
			oy := clampInt(y-r, 0, h-1)
			sum += float64(src.Pix[ny*stride+x]) - float64(src.Pix[oy*stride+x])
		}
	}
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func (r *PNGRenderer) renderInsetShadow(absX, absY, w, h float64, s style.Shadow, sc color.RGBA) {
	x := int(absX)
	y := int(absY)
	wi := int(w)
	hi := int(h)

	blur := int(s.Blur)
	if blur < 1 {
		blur = 1
	}
	spread := int(s.Spread)
	ox := int(s.OffsetX)
	oy := int(s.OffsetY)

	for ring := 0; ring < blur; ring++ {
		alpha := uint8(float64(sc.A) * float64(blur-ring) / float64(blur))
		c := color.RGBA{sc.R, sc.G, sc.B, alpha}
		inset := ring + spread

		topY := y + inset + oy
		if topY >= y && topY < y+hi {
			fillRect(r.img, x, topY, wi, 1, c)
		}
		botY := y + hi - 1 - inset + oy
		if botY >= y && botY < y+hi {
			fillRect(r.img, x, botY, wi, 1, c)
		}
		leftX := x + inset + ox
		if leftX >= x && leftX < x+wi {
			fillRect(r.img, leftX, y, 1, hi, c)
		}
		rightX := x + wi - 1 - inset + ox
		if rightX >= x && rightX < x+wi {
			fillRect(r.img, rightX, y, 1, hi, c)
		}
	}
}

func (r *PNGRenderer) cachedMask(w, h int, tl, tr, br, bl float64) *image.Alpha {
	key := maskKey{w, h, int(tl * 100), int(tr * 100), int(br * 100), int(bl * 100)}
	if m, ok := r.maskCache[key]; ok {
		return m
	}
	m := roundedMask(w, h, tl, tr, br, bl)
	if r.maskCache == nil {
		r.maskCache = make(map[maskKey]*image.Alpha)
	}
	if len(r.maskCache) < 32 {
		r.maskCache[key] = m
	}
	return m
}

func roundedMask(w, h int, tl, tr, br, bl float64) *image.Alpha {
	const s = 2
	sw, sh := w*s, h*s

	hi := image.NewAlpha(image.Rect(0, 0, sw, sh))
	for y := range sh {
		for x := range sw {
			hi.SetAlpha(x, y, color.Alpha{A: 255})
		}
	}

	type corner struct {
		r      float64
		cx, cy float64
		x0, y0 int
	}
	corners := []corner{
		{tl * s, tl * s, tl * s, 0, 0},
		{tr * s, float64(sw) - tr*s, tr * s, sw - int(tr*s), 0},
		{bl * s, bl * s, float64(sh) - bl*s, 0, sh - int(bl*s)},
		{br * s, float64(sw) - br*s, float64(sh) - br*s, sw - int(br*s), sh - int(br*s)},
	}
	for _, c := range corners {
		if c.r <= 0 {
			continue
		}
		ri := int(c.r)
		for ly := range ri {
			for lx := range ri {
				px := c.x0 + lx
				py := c.y0 + ly
				if px < 0 || px >= sw || py < 0 || py >= sh {
					continue
				}
				dx := float64(px) + 0.5 - c.cx
				dy := float64(py) + 0.5 - c.cy
				dist := math.Sqrt(dx*dx+dy*dy) - c.r
				if dist >= 0.5 {
					hi.SetAlpha(px, py, color.Alpha{A: 0})
				} else if dist > -0.5 {
					a := uint8((0.5 - dist) * 255)
					hi.SetAlpha(px, py, color.Alpha{A: a})
				}
			}
		}
	}

	mask := image.NewAlpha(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			a00 := uint32(hi.AlphaAt(x*s, y*s).A)
			a10 := uint32(hi.AlphaAt(x*s+1, y*s).A)
			a01 := uint32(hi.AlphaAt(x*s, y*s+1).A)
			a11 := uint32(hi.AlphaAt(x*s+1, y*s+1).A)
			mask.SetAlpha(x, y, color.Alpha{A: uint8((a00 + a10 + a01 + a11 + 2) / 4)})
		}
	}
	return mask
}

func (r *PNGRenderer) renderInlineSVG(pn *parse.Node, cs *style.ComputedStyle, absX, absY, w, h float64) {
	svgXML := SerializeSVGNode(pn)
	img, err := rasterizeSVG([]byte(svgXML), int(w), int(h))
	if err != nil {
		return
	}
	dstRect := image.Rect(int(absX), int(absY), int(absX+w), int(absY+h))
	hasRadius := cs.BorderTopLeftRadius > 0 || cs.BorderTopRightRadius > 0 ||
		cs.BorderBottomLeftRadius > 0 || cs.BorderBottomRightRadius > 0
	if hasRadius {
		tmp := image.NewRGBA(dstRect.Sub(dstRect.Min))
		xdraw.BiLinear.Scale(tmp, tmp.Bounds(), img, img.Bounds(), xdraw.Over, nil)
		mask := r.cachedMask(int(w), int(h),
			cs.BorderTopLeftRadius, cs.BorderTopRightRadius,
			cs.BorderBottomRightRadius, cs.BorderBottomLeftRadius)
		draw.DrawMask(r.img, dstRect, tmp, image.Point{}, mask, image.Point{}, draw.Over)
	} else {
		xdraw.BiLinear.Scale(r.img, dstRect, img, img.Bounds(), xdraw.Over, nil)
	}
}

func (r *PNGRenderer) renderImage(src string, cs *style.ComputedStyle, absX, absY, w, h float64) {
	dataURI, ok := resolveImageSource(src)
	if !ok {
		return
	}

	img, err := decodeDataURI(dataURI, int(w), int(h))
	if err != nil {
		return
	}

	dstRect := image.Rect(int(absX), int(absY), int(absX+w), int(absY+h))

	hasRadius := cs.BorderTopLeftRadius > 0 || cs.BorderTopRightRadius > 0 ||
		cs.BorderBottomLeftRadius > 0 || cs.BorderBottomRightRadius > 0

	if hasRadius {
		tmp := image.NewRGBA(dstRect.Sub(dstRect.Min))
		xdraw.BiLinear.Scale(tmp, tmp.Bounds(), img, img.Bounds(), xdraw.Over, nil)
		mask := r.cachedMask(int(w), int(h),
			cs.BorderTopLeftRadius, cs.BorderTopRightRadius,
			cs.BorderBottomRightRadius, cs.BorderBottomLeftRadius)
		draw.DrawMask(r.img, dstRect, tmp, image.Point{}, mask, image.Point{}, draw.Over)
	} else {
		xdraw.BiLinear.Scale(r.img, dstRect, img, img.Bounds(), xdraw.Over, nil)
	}
}

func decodeDataURI(uri string, targetSize ...int) (image.Image, error) {
	idx := strings.Index(uri, ",")
	if idx < 0 {
		return nil, image.ErrFormat
	}
	header := uri[:idx]
	payload := uri[idx+1:]

	if strings.Contains(header, "image/svg") {
		var data []byte
		if strings.Contains(header, "base64") {
			var err error
			data, err = base64.StdEncoding.DecodeString(payload)
			if err != nil {
				return nil, err
			}
		} else {
			decoded, err := urlDecode(payload)
			if err != nil {
				return nil, err
			}
			data = []byte(decoded)
		}
		tw, th := 0, 0
		if len(targetSize) >= 2 {
			tw, th = targetSize[0], targetSize[1]
		}
		return rasterizeSVG(data, tw, th)
	}

	data, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(data)
	if strings.Contains(header, "image/png") {
		return png.Decode(r)
	}
	if strings.Contains(header, "image/jpeg") {
		return jpeg.Decode(r)
	}
	img, _, err := image.Decode(r)
	return img, err
}

func urlDecode(s string) (string, error) {
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		if s[i] == '%' && i+2 < len(s) {
			hi := unhex(s[i+1])
			lo := unhex(s[i+2])
			if hi >= 0 && lo >= 0 {
				b.WriteByte(byte(hi<<4 | lo))
				i += 3
				continue
			}
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String(), nil
}

func unhex(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	}
	return -1
}
