// Package render converts layout trees to SVG, PNG, and JPEG output.
package render

import (
	"fmt"
	"strings"

	"github.com/macawls/ogre/font"
	"github.com/macawls/ogre/layout"
	"github.com/macawls/ogre/parse"
	"github.com/macawls/ogre/style"
)

type RenderContext struct {
	Width         int
	Height        int
	Styles        map[*parse.Node]*style.ComputedStyle
	NodeMap       map[*parse.Node]*layout.Node
	WrappedText   map[*parse.Node][]font.TextLine
	FontMgr       *font.Manager
	EmojiProvider *font.EmojiProvider
	reverse       map[*layout.Node]*parse.Node
	ids           idGen
}

type idGen struct{ n int }

func (g *idGen) next(prefix string) string {
	g.n++
	return fmt.Sprintf("%s%d", prefix, g.n)
}

type SVGOptions struct {
	FontMgr       *font.Manager
	EmojiProvider *font.EmojiProvider
}

// RenderSVG generates the corresponding output format.
func RenderSVG(tree *layout.LayoutTree, styles map[*parse.Node]*style.ComputedStyle, wrappedText map[*parse.Node][]font.TextLine, width, height int, fontMgr ...*font.Manager) string {
	return RenderSVGWithOptions(tree, styles, wrappedText, width, height, SVGOptions{
		FontMgr: firstMgr(fontMgr),
	})
}

func firstMgr(mgrs []*font.Manager) *font.Manager {
	if len(mgrs) > 0 {
		return mgrs[0]
	}
	return nil
}

// RenderSVGWithOptions generates the corresponding output format.
func RenderSVGWithOptions(tree *layout.LayoutTree, styles map[*parse.Node]*style.ComputedStyle, wrappedText map[*parse.Node][]font.TextLine, width, height int, opts SVGOptions) string {
	ctx := &RenderContext{
		Width:         width,
		Height:        height,
		Styles:        styles,
		NodeMap:       tree.NodeMap,
		WrappedText:   wrappedText,
		FontMgr:       opts.FontMgr,
		EmojiProvider: opts.EmojiProvider,
		reverse:       make(map[*layout.Node]*parse.Node, len(tree.NodeMap)),
	}
	for pn, ln := range tree.NodeMap {
		ctx.reverse[ln] = pn
	}

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, width, height, width, height)

	if tree.Root != nil {
		pn := ctx.reverse[tree.Root]
		cs := styles[pn]
		renderNode(&b, tree.Root, pn, cs, ctx)
	}

	b.WriteString("</svg>")
	return b.String()
}

func renderNode(b *strings.Builder, node *layout.Node, pn *parse.Node, cs *style.ComputedStyle, ctx *RenderContext) {
	renderNodeAt(b, node, pn, cs, ctx, 0, 0)
}

func renderNodeAt(b *strings.Builder, node *layout.Node, pn *parse.Node, cs *style.ComputedStyle, ctx *RenderContext, parentX, parentY float64) {
	if cs == nil {
		cs = style.NewComputedStyle()
	}

	l := node.Layout
	absX := parentX + l.X
	absY := parentY + l.Y
	l.X = absX
	l.Y = absY
	needsOpacity := cs.Opacity > 0 && cs.Opacity < 1

	defsContent, clipAttr := RenderOverflowClip(cs, l.X, l.Y, l.Width, l.Height, func(prefix string) string {
		return ctx.ids.next(prefix)
	})

	if defsContent != "" {
		fmt.Fprintf(b, "<defs>%s</defs>", defsContent)
	}

	if needsOpacity {
		fmt.Fprintf(b, `<g opacity="%.4g">`, cs.Opacity)
	}

	filterDefs, filterAttr := RenderCSSFilter(cs.Filter, func(prefix string) string {
		return ctx.ids.next(prefix)
	})
	needsFilter := filterAttr != ""
	if filterDefs != "" {
		fmt.Fprintf(b, "<defs>%s</defs>", filterDefs)
	}
	if needsFilter {
		fmt.Fprintf(b, `<g %s>`, filterAttr)
	}

	needsClipGroup := clipAttr != ""
	if needsClipGroup {
		fmt.Fprintf(b, `<g %s>`, clipAttr)
	}

	if pn != nil && pn.Type == parse.TextNode {
		if lines, ok := ctx.WrappedText[pn]; ok && len(lines) > 0 {
			result := RenderTextWithIDGen(lines, cs, absX, absY, l.Width, l.Height, func(prefix string) string {
				return ctx.ids.next(prefix)
			}, ctx.FontMgr, ctx.EmojiProvider)
			b.WriteString(result.Shadows)
			b.WriteString(result.Content)
			b.WriteString(result.Decorations)
		} else {
			renderTextAt(b, l, pn, cs, ctx.FontMgr)
		}
	} else {
		if cs.BoxShadow != "" {
			parsed, _ := style.ParseBoxShadow(cs.BoxShadow)
			shadows := RenderBoxShadow(parsed, l.X, l.Y, l.Width, l.Height,
				maxRadius(cs.BorderTopLeftRadius, cs.BorderBottomLeftRadius),
				func(prefix string) string { return ctx.ids.next(prefix) })
			if shadows != "" {
				b.WriteString(shadows)
			}
		}

		bgResult := RenderBackground(cs, l.X, l.Y, l.Width, l.Height,
			func(prefix string) string { return ctx.ids.next(prefix) })
		if bgResult.Defs != "" {
			fmt.Fprintf(b, "<defs>%s</defs>", bgResult.Defs)
		}
		if len(bgResult.Layers) > 1 {
			rx := maxRadius(cs.BorderTopLeftRadius, cs.BorderBottomLeftRadius)
			ry := maxRadius(cs.BorderTopRightRadius, cs.BorderBottomRightRadius)
			for i := len(bgResult.Layers) - 1; i >= 0; i-- {
				layer := bgResult.Layers[i]
				if layer.Fill == "" {
					continue
				}
				fmt.Fprintf(b, `<rect x="%.4g" y="%.4g" width="%.4g" height="%.4g" fill="%s"`,
					l.X, l.Y, l.Width, l.Height, layer.Fill)
				if rx > 0 {
					fmt.Fprintf(b, ` rx="%.4g"`, rx)
				}
				if ry > 0 {
					fmt.Fprintf(b, ` ry="%.4g"`, ry)
				}
				b.WriteString("/>")
			}
		} else if bgResult.Fill != "" && bgResult.Fill != "none" {
			rx := maxRadius(cs.BorderTopLeftRadius, cs.BorderBottomLeftRadius)
			ry := maxRadius(cs.BorderTopRightRadius, cs.BorderBottomRightRadius)
			fmt.Fprintf(b, `<rect x="%.4g" y="%.4g" width="%.4g" height="%.4g" fill="%s"`,
				l.X, l.Y, l.Width, l.Height, bgResult.Fill)
			if rx > 0 {
				fmt.Fprintf(b, ` rx="%.4g"`, rx)
			}
			if ry > 0 {
				fmt.Fprintf(b, ` ry="%.4g"`, ry)
			}
			b.WriteString("/>")
		}

		borders := RenderBorders(cs, l.X, l.Y, l.Width, l.Height)
		if borders != "" {
			b.WriteString(borders)
		}

		if cs.BoxShadow != "" {
			parsed, _ := style.ParseBoxShadow(cs.BoxShadow)
			insetShadows := RenderInsetBoxShadow(parsed, l.X, l.Y, l.Width, l.Height,
				maxRadius(cs.BorderTopLeftRadius, cs.BorderBottomLeftRadius),
				func(prefix string) string { return ctx.ids.next(prefix) })
			if insetShadows != "" {
				b.WriteString(insetShadows)
			}
		}

		if cs.Transform != "" {
			transform := RenderTransform(cs.Transform, cs.TransformOrigin, l.X, l.Y, l.Width, l.Height)
			if transform != "" {
				fmt.Fprintf(b, `<g transform="%s">`, transform)
			}
		}

		if pn != nil && pn.Tag == "img" {
			if src := pn.Attrs["src"]; src != "" {
				b.WriteString(RenderImage(src, cs, absX, absY, l.Width, l.Height))
			}
		} else if pn != nil && pn.Tag == "svg" {
			b.WriteString(RenderInlineSVG(pn, absX, absY, l.Width, l.Height))
		} else {
			for _, child := range node.Children {
				cpn := ctx.reverse[child]
				ccs := ctx.Styles[cpn]
				renderNodeAt(b, child, cpn, ccs, ctx, absX, absY)
			}
		}

		if cs.Transform != "" {
			b.WriteString("</g>")
		}
	}

	if needsClipGroup {
		b.WriteString("</g>")
	}

	if needsFilter {
		b.WriteString("</g>")
	}

	if needsOpacity {
		b.WriteString("</g>")
	}
}

func renderTextAt(b *strings.Builder, l layout.Layout, pn *parse.Node, cs *style.ComputedStyle, fontMgr *font.Manager) {
	family := cs.FontFamily
	if family == "" {
		family = "sans-serif"
	}
	size := cs.FontSize
	if size == 0 {
		size = 16
	}
	weight := cs.FontWeight
	if weight == 0 {
		weight = 400
	}

	fill := colorToCSS(cs.Color)
	if cs.Color.A == -1 {
		fill = "#000000"
	}

	ascent := size * 0.8
	if fontMgr != nil {
		fontStyle := cs.FontStyle
		if fontStyle == "" {
			fontStyle = "normal"
		}
		face := fontMgr.Resolve(family, weight, fontStyle)
		if face != nil {
			ff, err := fontMgr.NewFace(face, size)
			if err == nil {
				ascent = font.Ascent(ff)
			}
		}
	}
	y := l.Y + ascent

	if fontMgr != nil {
		fontStyle := cs.FontStyle
		if fontStyle == "" {
			fontStyle = "normal"
		}
		rtl := cs.Direction == "rtl"
		var pathD string
		if rtl || needsShaping(pn.Text) {
			pathD, _ = font.ShapedTextToPath(fontMgr, pn.Text, family, weight, fontStyle, size, rtl)
		} else {
			pathD, _ = font.TextToPath(fontMgr, pn.Text, family, weight, fontStyle, size)
		}
		if pathD != "" {
			fmt.Fprintf(b, `<path d="%s" fill="%s" transform="translate(%.4g,%.4g)"/>`,
				pathD, fill, l.X, y)
			return
		}
	}

	fmt.Fprintf(b, `<text x="%.4g" y="%.4g" font-family="%s" font-size="%.4g" font-weight="%d" fill="%s"`,
		l.X, y, xmlEscape(family), size, weight, fill)
	if cs.LetterSpacing > 0 {
		fmt.Fprintf(b, ` letter-spacing="%.4gpx"`, cs.LetterSpacing)
	}
	fmt.Fprintf(b, ">")
	b.WriteString(xmlEscape(pn.Text))
	b.WriteString("</text>")
}

func colorToCSS(c style.Color) string {
	if c.A == 1.0 {
		return c.Hex()
	}
	return fmt.Sprintf("rgba(%d,%d,%d,%.4g)", c.R, c.G, c.B, c.A)
}

func xmlEscape(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&apos;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func maxRadius(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
