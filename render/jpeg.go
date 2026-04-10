package render

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"

	fontpkg "github.com/macawls/ogre/font"
	"github.com/macawls/ogre/layout"
	"github.com/macawls/ogre/parse"
	"github.com/macawls/ogre/style"
)

// RenderJPEG generates the corresponding output format.
func RenderJPEG(tree *layout.LayoutTree, styles map[*parse.Node]*style.ComputedStyle, fonts *fontpkg.Manager, width, height, quality int, opts ...PNGOptions) ([]byte, error) {
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

	if quality <= 0 || quality > 100 {
		quality = 90
	}

	var buf bytes.Buffer
	if err := encodeJPEG444(&buf, img, quality); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
