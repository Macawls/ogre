// Package ogre renders HTML and CSS to SVG, PNG, or JPEG images.
package ogre

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/macawls/ogre/font"
	"github.com/macawls/ogre/layout"
	"github.com/macawls/ogre/parse"
	"github.com/macawls/ogre/render"
	"github.com/macawls/ogre/style"
)

// Format specifies the output image format.
type Format string

const (
	// FormatSVG renders output as SVG.
	FormatSVG Format = "svg"
	// FormatPNG renders output as PNG.
	FormatPNG Format = "png"
	// FormatJPEG renders output as JPEG.
	FormatJPEG Format = "jpeg"
)

// Options configures a render operation including dimensions, format, and fonts.
type Options struct {
	Width         int
	Height        int
	Format        Format
	Quality       int
	Fonts         []FontSource
	Debug         bool
	EmojiProvider string // "twemoji" (default), "none"
	MaxElements   int
}

// FontSource describes a font to load, either from raw bytes or a URL.
type FontSource struct {
	Name   string
	Weight int
	Style  string
	Data   []byte
	URL    string
}

// Result holds the rendered image data and metadata.
type Result struct {
	Data        []byte
	ContentType string
	Width       int
	Height      int
}

// Renderer manages fonts and converts HTML to images.
type Renderer struct {
	fonts     *font.Manager
	fontCache *font.FontCache
}

// NewRenderer creates a Renderer with default fonts loaded.
func NewRenderer() *Renderer {
	fm := font.NewManager()
	_ = fm.LoadDefaults()
	return &Renderer{fonts: fm}
}

// LoadFont registers a font for use in subsequent renders.
func (r *Renderer) LoadFont(src FontSource) error {
	data := src.Data
	if len(data) == 0 && src.URL != "" {
		if r.fontCache == nil {
			cacheDir := filepath.Join(os.TempDir(), "ogre-font-cache")
			r.fontCache = font.NewFontCache(cacheDir)
		}
		var err error
		data, err = r.fontCache.Fetch(src.URL)
		if err != nil {
			return fmt.Errorf("fetch font %q: %w", src.URL, err)
		}
	}
	return r.fonts.LoadFont(font.FontSource{
		Name:   src.Name,
		Weight: src.Weight,
		Style:  src.Style,
		Data:   data,
	})
}

// Render converts an HTML string to an image using the given options.
func (r *Renderer) Render(html string, opts Options) (*Result, error) {
	if opts.Width <= 0 {
		opts.Width = 1200
	}
	if opts.Height <= 0 {
		opts.Height = 630
	}
	if opts.Format == "" {
		opts.Format = FormatSVG
	}

	for _, src := range opts.Fonts {
		if err := r.LoadFont(src); err != nil {
			return nil, err
		}
	}

	root, err := parse.Parse(html)
	if err != nil {
		return nil, err
	}

	if opts.MaxElements > 0 {
		count := parse.CountNodes(root)
		if count > opts.MaxElements {
			return nil, fmt.Errorf("HTML exceeds maximum element count: %d > %d", count, opts.MaxElements)
		}
	}

	w := float64(opts.Width)
	h := float64(opts.Height)

	styles := style.Resolve(root, w, h)

	if r.fontCache == nil {
		cacheDir := filepath.Join(os.TempDir(), "ogre-font-cache")
		r.fontCache = font.NewFontCache(cacheDir)
	}
	seen := make(map[string]bool)
	for _, cs := range styles {
		fam := cs.FontFamily
		if fam == "" || fam == "sans-serif" || fam == "default" || seen[fam] {
			continue
		}
		seen[fam] = true
		if r.fonts.HasFamily(fam) {
			continue
		}
		data, err := font.FetchGoogleFont(fam, 400, r.fontCache)
		if err == nil {
			_ = r.fonts.LoadFont(font.FontSource{Name: fam, Weight: 400, Data: data})
		}
		data700, err := font.FetchGoogleFont(fam, 700, r.fontCache)
		if err == nil {
			_ = r.fonts.LoadFont(font.FontSource{Name: fam, Weight: 700, Data: data700})
		}
	}

	wrappedText := make(map[*parse.Node][]font.TextLine)

	measureText := func(pn *parse.Node, text string, cs *style.ComputedStyle, maxWidth float64) (float64, float64) {
		face := r.fonts.Resolve(cs.FontFamily, cs.FontWeight, cs.FontStyle)
		if face == nil {
			return 0, 0
		}
		ff, err := r.fonts.NewFace(face, cs.FontSize)
		if err != nil {
			return 0, 0
		}

		cfg := font.WrapConfig{
			MaxWidth:      maxWidth,
			FontFace:      ff,
			FontSize:      cs.FontSize,
			LineHeight:    cs.LineHeight,
			LetterSpacing: cs.LetterSpacing,
			WhiteSpace:    int(cs.WhiteSpace),
			WordBreak:     int(cs.WordBreak),
			LineClamp:     cs.LineClamp,
			TextOverflow:  cs.TextOverflow,
		}

		lines := font.WrapText(text, cfg)
		if len(lines) == 0 {
			return 0, 0
		}

		wrappedText[pn] = lines

		var maxW float64
		for _, l := range lines {
			if l.Width > maxW {
				maxW = l.Width
			}
		}
		totalH := float64(len(lines)) * cs.LineHeight
		return maxW, totalH
	}

	tree := layout.ComputeLayout(root, styles, w, h, measureText)

	var emojiProvider *font.EmojiProvider
	if opts.EmojiProvider != "none" {
		emojiProvider = font.NewEmojiProvider()
	}

	if opts.Format == FormatPNG {
		pngData, err := render.RenderPNG(tree, styles, r.fonts, opts.Width, opts.Height, render.PNGOptions{
			WrappedText:   wrappedText,
			EmojiProvider: emojiProvider,
		})
		if err != nil {
			return nil, err
		}
		return &Result{Data: pngData, ContentType: "image/png", Width: opts.Width, Height: opts.Height}, nil
	}

	if opts.Format == FormatJPEG {
		quality := opts.Quality
		if quality <= 0 {
			quality = 90
		}
		jpegData, err := render.RenderJPEG(tree, styles, r.fonts, opts.Width, opts.Height, quality, render.PNGOptions{
			WrappedText:   wrappedText,
			EmojiProvider: emojiProvider,
		})
		if err != nil {
			return nil, err
		}
		return &Result{Data: jpegData, ContentType: "image/jpeg", Width: opts.Width, Height: opts.Height}, nil
	}

	svgOpts := render.SVGOptions{FontMgr: r.fonts}
	if emojiProvider != nil {
		svgOpts.EmojiProvider = emojiProvider
	}
	svg := render.RenderSVGWithOptions(tree, styles, wrappedText, opts.Width, opts.Height, svgOpts)

	return &Result{
		Data:        []byte(svg),
		ContentType: "image/svg+xml",
		Width:       opts.Width,
		Height:      opts.Height,
	}, nil
}

// Render is a convenience function that creates a temporary Renderer and renders the HTML.
func Render(html string, opts Options) (*Result, error) {
	r := NewRenderer()
	return r.Render(html, opts)
}
