package render

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	fontpkg "github.com/macawls/ogre/font"
	"github.com/macawls/ogre/layout"
	"github.com/macawls/ogre/parse"
	"github.com/macawls/ogre/style"
)

type PNGRenderer struct {
	img           *image.RGBA
	styles        map[*parse.Node]*style.ComputedStyle
	fonts         *fontpkg.Manager
	reverse       map[*layout.Node]*parse.Node
	wrappedText   map[*parse.Node][]fontpkg.TextLine
	emojiProvider *fontpkg.EmojiProvider
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
		tmp := image.NewRGBA(r.img.Bounds())
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

	r.renderBoxShadows(l, cs, absX, absY, false)

	hasRadius := cs.BorderTopLeftRadius > 0 || cs.BorderTopRightRadius > 0 ||
		cs.BorderBottomLeftRadius > 0 || cs.BorderBottomRightRadius > 0

	if cs.BackgroundImage != "" {
		if hasRadius {
			tmp := image.NewRGBA(r.img.Bounds())
			sub := &PNGRenderer{img: tmp, styles: r.styles, fonts: r.fonts, reverse: r.reverse, wrappedText: r.wrappedText, emojiProvider: r.emojiProvider}
			sub.renderGradient(absX, absY, l.Width, l.Height, cs)
			rect := image.Rect(int(absX), int(absY), int(absX+l.Width), int(absY+l.Height)).Intersect(r.img.Bounds())
			mask := roundedMask(int(l.Width), int(l.Height), cs.BorderTopLeftRadius, cs.BorderTopRightRadius, cs.BorderBottomRightRadius, cs.BorderBottomLeftRadius)
			draw.DrawMask(r.img, rect, tmp, rect.Min, mask, image.Point{}, draw.Over)
		} else {
			r.renderGradient(absX, absY, l.Width, l.Height, cs)
		}
	} else if !cs.BackgroundColor.IsTransparent() {
		c := styleToColor(cs.BackgroundColor)
		rect := image.Rect(int(absX), int(absY), int(absX+l.Width), int(absY+l.Height)).Intersect(r.img.Bounds())
		if hasRadius {
			mask := roundedMask(int(l.Width), int(l.Height), cs.BorderTopLeftRadius, cs.BorderTopRightRadius, cs.BorderBottomRightRadius, cs.BorderBottomLeftRadius)
			draw.DrawMask(r.img, rect, image.NewUniform(c), image.Point{}, mask, image.Point{}, draw.Over)
		} else {
			draw.Draw(r.img, rect, image.NewUniform(c), image.Point{}, draw.Over)
		}
	}

	r.renderBorders(absX, absY, l.Width, l.Height, cs)

	r.renderBoxShadows(l, cs, absX, absY, true)

	if cs.Overflow == style.OverflowHidden {
		clip := image.Rect(int(absX), int(absY), int(absX+l.Width), int(absY+l.Height))
		tmp := image.NewRGBA(r.img.Bounds())
		sub := &PNGRenderer{img: tmp, styles: r.styles, fonts: r.fonts, reverse: r.reverse, wrappedText: r.wrappedText, emojiProvider: r.emojiProvider}
		for _, child := range node.Children {
			cpn := sub.reverse[child]
			ccs := sub.styles[cpn]
			sub.renderNode(child, cpn, ccs, absX, absY)
		}
		hasRadius := cs.BorderTopLeftRadius > 0 || cs.BorderTopRightRadius > 0 ||
			cs.BorderBottomLeftRadius > 0 || cs.BorderBottomRightRadius > 0
		if hasRadius {
			mask := roundedMask(int(l.Width), int(l.Height),
				cs.BorderTopLeftRadius, cs.BorderTopRightRadius,
				cs.BorderBottomRightRadius, cs.BorderBottomLeftRadius)
			draw.DrawMask(r.img, clip, tmp, clip.Min, mask, image.Point{}, draw.Over)
		} else {
			draw.Draw(r.img, clip, tmp, clip.Min, draw.Over)
		}
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

func (r *PNGRenderer) drawTextWithEmoji(text string, x, y, ascent, size float64, tc color.RGBA, ff font.Face, cs *style.ComputedStyle) {
	if r.emojiProvider == nil || !containsEmoji(text) {
		drawer := &font.Drawer{
			Dst:  r.img,
			Src:  image.NewUniform(tc),
			Face: ff,
			Dot:  fixed.Point26_6{X: fixed.I(int(x)), Y: fixed.I(int(y))},
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
				Dot:  fixed.Point26_6{X: fixed.I(int(cx)), Y: fixed.I(int(y))},
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

	bounds := r.img.Bounds()
	for py := ry; py < ry+rh; py++ {
		if py < bounds.Min.Y || py >= bounds.Max.Y {
			continue
		}
		for px := rx; px < rx+rw; px++ {
			if px < bounds.Min.X || px >= bounds.Max.X {
				continue
			}
			dx := float64(px-rx) - cx
			dy := float64(py-ry) - cy
			t := (dx*sinA + dy*cosA) / length
			t += 0.5
			if g.Repeating {
				t = t - math.Floor(t)
			} else {
				t = math.Max(0, math.Min(1, t))
			}
			c := interpolateStops(g.Stops, t)
			r.img.SetRGBA(px, py, c)
		}
	}
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

	bounds := r.img.Bounds()
	for py := ry; py < ry+rh; py++ {
		if py < bounds.Min.Y || py >= bounds.Max.Y {
			continue
		}
		for px := rx; px < rx+rw; px++ {
			if px < bounds.Min.X || px >= bounds.Max.X {
				continue
			}
			dx := float64(px-rx) - cx
			dy := float64(py-ry) - cy
			dist := math.Sqrt(dx*dx + dy*dy)
			t := dist / maxDist
			if g.Repeating {
				t = t - math.Floor(t)
			} else {
				t = math.Max(0, math.Min(1, t))
			}
			c := interpolateStops(g.Stops, t)
			r.img.SetRGBA(px, py, c)
		}
	}
}

func interpolateStops(stops []style.ColorStop, t float64) color.RGBA {
	if len(stops) == 0 {
		return color.RGBA{0, 0, 0, 255}
	}
	if t <= stops[0].Position {
		return styleToColor(stops[0].Color)
	}
	if t >= stops[len(stops)-1].Position {
		return styleToColor(stops[len(stops)-1].Color)
	}
	for i := 1; i < len(stops); i++ {
		if t <= stops[i].Position {
			prev := stops[i-1]
			curr := stops[i]
			span := curr.Position - prev.Position
			if span <= 0 {
				return styleToColor(curr.Color)
			}
			f := (t - prev.Position) / span
			return color.RGBA{
				R: uint8(math.Round(float64(prev.Color.R) + f*(float64(curr.Color.R)-float64(prev.Color.R)))),
				G: uint8(math.Round(float64(prev.Color.G) + f*(float64(curr.Color.G)-float64(prev.Color.G)))),
				B: uint8(math.Round(float64(prev.Color.B) + f*(float64(curr.Color.B)-float64(prev.Color.B)))),
				A: uint8(math.Round((prev.Color.A + f*(curr.Color.A-prev.Color.A)) * 255)),
			}
		}
	}
	return styleToColor(stops[len(stops)-1].Color)
}

func fillRect(img *image.RGBA, x, y, w, h int, c color.Color) {
	rect := image.Rect(x, y, x+w, y+h).Intersect(img.Bounds())
	draw.Draw(img, rect, image.NewUniform(c), image.Point{}, draw.Over)
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
	w := b.Dx()
	div := float64(2*r + 1)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		sum := 0.0
		for x := -r; x <= r; x++ {
			cx := clampInt(x, 0, w-1)
			sum += float64(src.AlphaAt(cx, y).A)
		}
		for x := b.Min.X; x < b.Max.X; x++ {
			dst.SetAlpha(x, y, color.Alpha{A: uint8(math.Round(sum / div))})
			nx := clampInt(x+r+1, 0, w-1)
			ox := clampInt(x-r, 0, w-1)
			sum += float64(src.AlphaAt(nx, y).A) - float64(src.AlphaAt(ox, y).A)
		}
	}
}

func boxBlurV(src, dst *image.Alpha, b image.Rectangle, r int) {
	h := b.Dy()
	div := float64(2*r + 1)
	for x := b.Min.X; x < b.Max.X; x++ {
		sum := 0.0
		for y := -r; y <= r; y++ {
			cy := clampInt(y, 0, h-1)
			sum += float64(src.AlphaAt(x, cy).A)
		}
		for y := b.Min.Y; y < b.Max.Y; y++ {
			dst.SetAlpha(x, y, color.Alpha{A: uint8(math.Round(sum / div))})
			ny := clampInt(y+r+1, 0, h-1)
			oy := clampInt(y-r, 0, h-1)
			sum += float64(src.AlphaAt(x, ny).A) - float64(src.AlphaAt(x, oy).A)
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

func roundedMask(w, h int, tl, tr, br, bl float64) *image.Alpha {
	mask := image.NewAlpha(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			mask.SetAlpha(x, y, color.Alpha{A: 255})
		}
	}
	type corner struct {
		r      float64
		x0, y0 int
		flipX  bool
		flipY  bool
	}
	corners := []corner{
		{tl, 0, 0, false, false},
		{tr, w - int(tr), 0, true, false},
		{bl, 0, h - int(bl), false, true},
		{br, w - int(br), h - int(br), true, true},
	}
	for _, c := range corners {
		if c.r <= 0 {
			continue
		}
		ri := int(c.r)
		for ly := 0; ly < ri; ly++ {
			for lx := 0; lx < ri; lx++ {
				var dx, dy float64
				if c.flipX {
					dx = float64(lx) + 0.5
				} else {
					dx = float64(ri-lx) - 0.5
				}
				if c.flipY {
					dy = float64(ly) + 0.5
				} else {
					dy = float64(ri-ly) - 0.5
				}
				if dx*dx+dy*dy > c.r*c.r {
					px := c.x0 + lx
					py := c.y0 + ly
					if px >= 0 && px < w && py >= 0 && py < h {
						mask.SetAlpha(px, py, color.Alpha{A: 0})
					}
				}
			}
		}
	}
	return mask
}
