<p align="center">
  <img src="mascot.png" width="180" alt="Ogre">
</p>

<h1 align="center">Ogre</h1>

<p align="center">
  Pure Go HTML/CSS to SVG/PNG/JPEG renderer for OpenGraph images.<br>
  Vercel Satori alternative. No CGo.
</p>

## Features

- HTML + inline CSS to SVG, PNG, or JPEG
- Tailwind CSS v3 utility classes (no build step)
- Flexbox layout engine (W3C spec)
- Complex script rendering via pure Go HarfBuzz port (Arabic, Hebrew, Devanagari, Thai)
- RTL text support with Unicode bidi algorithm
- Emoji rendering (Twemoji, OpenMoji, Noto) in SVG and PNG
- Font embedding as SVG paths (self-contained SVGs)
- Built-in Go fonts, Google Fonts auto-fetch, CDN loading, WOFF decompression
- HTTP server with LRU caching, rate limiting, templates
- Go library with `http.Handler` integration
- JSX-style Go builder API
- Tailwind filter and transform classes (blur, scale, rotate, etc.)
- 95%+ pixel accuracy vs Satori across 25 test fixtures

## Quick Start

### As a CLI

```bash
go install github.com/macawls/ogre/cmd/ogre@latest

# Render HTML file to SVG
ogre --render template.html --output og.svg

# Render inline HTML to PNG
ogre --html '<div class="flex w-full h-full bg-blue-500 items-center justify-center"><div class="text-4xl font-bold text-white">Hello</div></div>' --output og.png --format png

# Start HTTP server
ogre --serve --port 3000
```

CLI flags:

| Flag       | Default | Description                                                    |
| ---------- | ------- | -------------------------------------------------------------- |
| `--serve`  | false   | Start HTTP server mode                                         |
| `--port`   | 3000    | Server port                                                    |
| `--render` |         | Path to HTML file to render                                    |
| `--html`   |         | Inline HTML string to render                                   |
| `--output` |         | Output file path (stdout for SVG if omitted, required for PNG/JPEG) |
| `--width`  | 1200    | Canvas width in pixels                                         |
| `--height` | 630     | Canvas height in pixels                                        |
| `--format` | svg     | Output format: `svg`, `png`, or `jpeg`                         |

### As a Go Library

```go
package main

import (
    "os"

    "github.com/macawls/ogre"
)

func main() {
    result, err := ogre.Render(`
        <div class="flex flex-col w-full h-full bg-slate-900 p-16 justify-center">
            <div class="text-5xl font-bold text-white">My Blog Post</div>
            <div class="text-xl text-slate-400 mt-4">A subtitle here</div>
        </div>
    `, ogre.Options{Width: 1200, Height: 630})
    if err != nil {
        panic(err)
    }

    os.WriteFile("og.svg", result.Data, 0644)
}
```

### Shared Renderer (recommended for servers)

```go
r := ogre.NewRenderer()
// Reuses font manager across renders — thread-safe
result, _ := r.Render(html, ogre.Options{Width: 1200, Height: 630})
```

### Custom Fonts

```go
// From file data
fontData, _ := os.ReadFile("Inter-Regular.ttf")
result, _ := ogre.Render(html, ogre.Options{
    Fonts: []ogre.FontSource{{
        Name: "Inter", Weight: 400, Style: "normal", Data: fontData,
    }},
})

// From URL (fetched and cached automatically)
result, _ := ogre.Render(html, ogre.Options{
    Fonts: []ogre.FontSource{{
        Name: "Inter", Weight: 400, URL: "https://example.com/Inter-Regular.ttf",
    }},
})
```

## API Reference

### Types

```go
type Options struct {
    Width         int          // Canvas width (default 1200)
    Height        int          // Canvas height (default 630)
    Format        Format       // "svg" (default), "png", or "jpeg"
    Quality       int          // JPEG quality 1-100 (default 90)
    Fonts         []FontSource // Custom fonts to load
    Debug         bool
    EmojiProvider string       // "twemoji" (default), "none"
}

type FontSource struct {
    Name   string // Font family name
    Weight int    // Font weight (100-900)
    Style  string // "normal" or "italic"
    Data   []byte // Raw font data (TTF/OTF/WOFF)
    URL    string // URL to fetch font from (alternative to Data)
}

type Result struct {
    Data        []byte // Rendered output bytes
    ContentType string // "image/svg+xml", "image/png", or "image/jpeg"
    Width       int
    Height      int
}
```

### Functions

- `ogre.Render(html string, opts Options) (*Result, error)` -- One-shot render. Creates a new font manager each call.
- `ogre.NewRenderer() *Renderer` -- Creates a shared renderer with pre-loaded default fonts.
- `(*Renderer).Render(html string, opts Options) (*Result, error)` -- Render with shared font manager. Thread-safe.
- `(*Renderer).LoadFont(src FontSource) error` -- Pre-load a font into the shared manager.

## HTTP API

### POST /render

Render HTML to SVG, PNG, or JPEG.

Request body (JSON):

```json
{
  "html": "<div class=\"flex w-full h-full bg-blue-500\">...</div>",
  "width": 1200,
  "height": 630,
  "format": "svg",
  "quality": 90
}
```

The `quality` field (1-100) controls JPEG compression. Ignored for SVG and PNG. Default 90.

Response: image bytes with appropriate `Content-Type` header (`image/svg+xml`, `image/png`, or `image/jpeg`).

Response headers include `ETag`, `Cache-Control`, and `X-Cache` (HIT/MISS).

### POST /render/template

Render a Go `html/template` with data substitution.

Request body (JSON):

```json
{
  "template": "<div class=\"text-4xl\">{{.Title}}</div>",
  "data": { "Title": "Hello World" },
  "width": 1200,
  "height": 630,
  "format": "svg"
}
```

Response: same as `/render`.

### GET /health

Returns `{"status":"ok"}`.

### CORS

All endpoints return `Access-Control-Allow-Origin: *`.

### Limits

- Max request body size: 10 MB
- Default cache size: 64 MB (LRU, keyed by SHA-256 of input)

## Supported CSS Properties

### Layout

- `display`: flex, none, block, contents
- `position`: static, relative, absolute
- `top`, `right`, `bottom`, `left`
- `width`, `height`, `min-width`, `min-height`, `max-width`, `max-height`
- `aspect-ratio`
- `overflow`: visible, hidden
- `box-sizing`: content-box, border-box

### Flexbox

- `flex-direction`: row, row-reverse, column, column-reverse
- `flex-wrap`: nowrap, wrap, wrap-reverse
- `flex-grow`, `flex-shrink`, `flex-basis`
- `align-items`: auto, flex-start, flex-end, center, stretch, baseline, space-between, space-around
- `align-self`: auto, flex-start, flex-end, center, stretch, baseline
- `align-content`: auto, flex-start, flex-end, center, stretch, space-between, space-around
- `justify-content`: flex-start, flex-end, center, space-between, space-around, space-evenly
- `gap`, `row-gap`, `column-gap`

### Box Model

- `margin` (all sides, shorthand)
- `padding` (all sides, shorthand)
- `border-width`, `border-style`, `border-color` (all sides, shorthand)
- `border-radius` (all corners, shorthand)

### Typography

- `font-family`, `font-size`, `font-weight`, `font-style`
- `color`
- `line-height`, `letter-spacing`
- `text-align`: left, right, center, justify, start, end
- `text-transform`: none, uppercase, lowercase, capitalize
- `text-decoration-line`: none, underline, overline, line-through
- `text-decoration-color`, `text-decoration-style`
- `text-shadow`
- `white-space`: normal, nowrap, pre, pre-wrap, pre-line
- `word-break`: normal, break-all, break-word, keep-all
- `text-overflow`: ellipsis
- `-webkit-line-clamp`

### Background

- `background-color`
- `background-image` (linear-gradient, radial-gradient, url())
- `background-size`, `background-position`, `background-repeat`

### Visual

- `opacity`
- `box-shadow`
- `transform`, `transform-origin`
- `object-fit`: fill, contain, cover, scale-down, none
- `object-position`
- `filter`: blur, grayscale, brightness
- `clip-path`

### Shorthands

- `margin`, `padding`, `border`, `border-radius`
- `flex`, `gap`, `background`, `font`
- `text-decoration`, `overflow`
- `border-top`, `border-right`, `border-bottom`, `border-left`
- `border-width`, `border-style`, `border-color`

## Tailwind Support

Ogre resolves Tailwind CSS v3 utility classes directly. No build step or Tailwind CLI needed.

### Supported Utility Categories

**Layout**: `flex`, `flex-row`, `flex-col`, `flex-wrap`, `flex-nowrap`, `flex-1`, `flex-auto`, `flex-initial`, `flex-none`, `flex-grow`, `flex-grow-0`, `flex-shrink`, `flex-shrink-0`, `hidden`, `block`, `relative`, `absolute`

**Alignment**: `items-start`, `items-end`, `items-center`, `items-stretch`, `items-baseline`, `justify-start`, `justify-end`, `justify-center`, `justify-between`, `justify-around`, `justify-evenly`, `self-auto`, `self-start`, `self-end`, `self-center`, `self-stretch`, `content-start`, `content-end`, `content-center`, `content-between`, `content-around`, `content-stretch`

**Spacing**: `p-{n}`, `px-{n}`, `py-{n}`, `pt-{n}`, `pr-{n}`, `pb-{n}`, `pl-{n}`, `m-{n}`, `mx-{n}`, `my-{n}`, `mt-{n}`, `mr-{n}`, `mb-{n}`, `ml-{n}`, `gap-{n}`, `gap-x-{n}`, `gap-y-{n}`, `space-x-{n}`, `space-y-{n}`

Spacing scale: `0` = 0px, `px` = 1px, `0.5` = 2px, `1` = 4px, `1.5` = 6px, `2` = 8px, `2.5` = 10px, `3` = 12px, `3.5` = 14px, then `{n}` = n\*4px up to 96.

**Sizing**: `w-{n}`, `h-{n}`, `size-{n}`, `w-full`, `h-full`, `w-screen`, `h-screen`, `w-auto`, `h-auto`, `w-fit`, `h-fit`, fraction values (`w-1/2`, `w-1/3`, `w-2/3`, `w-1/4`, `w-3/4`, etc.), `min-w-0`, `min-w-full`, `min-h-0`, `min-h-full`, `min-h-screen`, `max-w-sm` through `max-w-2xl`, `max-w-full`, `max-w-none`, `max-h-full`, `max-h-screen`, `max-h-none`

**Typography**: `text-xs` through `text-9xl`, `font-thin` through `font-black`, `text-left`, `text-center`, `text-right`, `text-justify`, `italic`, `not-italic`, `uppercase`, `lowercase`, `capitalize`, `normal-case`, `underline`, `overline`, `line-through`, `no-underline`, `leading-none`, `leading-tight`, `leading-normal`, `leading-loose`, `leading-{n}`, `tracking-tighter` through `tracking-widest`, `truncate`, `whitespace-normal`, `whitespace-nowrap`, `whitespace-pre`, `whitespace-pre-wrap`, `line-clamp-{1-6}`

**Colors**: `text-{color}-{shade}`, `bg-{color}-{shade}`, `border-{color}-{shade}`, plus `text-white`, `text-black`, `text-transparent`, `bg-white`, `bg-black`, `bg-transparent`

Available color palettes: slate, gray, zinc, neutral, stone, red, orange, amber, yellow, lime, green, emerald, teal, cyan, sky, blue, indigo, violet, purple, fuchsia, pink, rose. Shades: 50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950.

**Borders**: `border`, `border-0`, `border-2`, `border-4`, `border-8`, `border-t-{n}`, `border-r-{n}`, `border-b-{n}`, `border-l-{n}`, `border-solid`, `border-dashed`, `border-dotted`, `rounded-none` through `rounded-full`

**Effects**: `shadow-sm` through `shadow-2xl`, `shadow-none`, `opacity-{0-100}`

**Position**: `z-0` through `z-50`, `z-auto`, `top-{n}`, `right-{n}`, `bottom-{n}`, `left-{n}`, `inset-{n}`, `aspect-square`, `aspect-video`, `aspect-auto`

**Overflow**: `overflow-hidden`, `overflow-visible`

### Arbitrary Values

Use bracket notation for custom values:

```html
<div
  class="text-[32px] bg-[#ff5500] w-[200px] p-[20px] rounded-[12px] gap-[8px] leading-[1.5] tracking-[0.05em]"
></div>
```

## Comparison with Satori

| Feature          | Ogre                      | Satori            |
| ---------------- | ------------------------- | ----------------- |
| Language         | Go                        | TypeScript        |
| Output formats   | SVG, PNG, JPEG            | SVG only          |
| Dependencies     | stdlib + golang.org/x     | yoga-wasm, others |
| Binary size      | Single static binary      | Node.js runtime   |
| Tailwind support | Built-in (v3)             | Via plugin        |
| Font embedding   | SVG paths                 | SVG paths         |
| Layout engine    | Custom flexbox (W3C spec) | Yoga (via WASM)   |
| Emoji            | Twemoji CDN               | Twemoji CDN       |
| `<div>` default  | `display: flex`           | `display: flex`   |
| PNG output       | Built-in                  | Requires resvg    |
| Pixel accuracy   | 95%+ vs Satori            | Reference         |

## Docker

```bash
docker build -t ogre .
docker run -p 3000:3000 ogre
```

The image uses a multi-stage build. The final image uses Google's distroless base with only the static binary.

The container starts in server mode by default, listening on port 3000.

## Architecture

The rendering pipeline has four stages:

1. **Parse** (`parse/`) -- HTML string is parsed into a node tree. Inline styles and class attributes are extracted.
2. **Style** (`style/`) -- Tailwind classes are resolved to CSS properties. Shorthands are expanded. CSS values are parsed and computed. Properties inherit where appropriate.
3. **Layout** (`layout/`) -- A custom flexbox layout engine computes the position and size of every node. Text nodes are measured using the font manager to determine line breaks and dimensions.
4. **Render** (`render/`) -- The layout tree is rendered to SVG (with font glyphs converted to path data) or rasterized to PNG.

### Package Structure

- `cmd/ogre/` -- CLI entry point and HTTP server startup
- `parse/` -- HTML parsing, node tree
- `style/` -- CSS property definitions, shorthand expansion, Tailwind resolver, inheritance, computed values
- `layout/` -- Flexbox layout engine (W3C spec, not a Yoga port)
- `font/` -- Font loading, text measurement, glyph path extraction, WOFF decompression, emoji support
- `render/` -- SVG generation and PNG rasterization
- `server/` -- HTTP API, LRU caching, template rendering

### Dependencies

Only `golang.org/x/*` packages are used. No third-party imports.

- `golang.org/x/net/html` -- HTML parsing
- `golang.org/x/image/font` -- Font interfaces
- `golang.org/x/image/font/opentype` -- OTF/TTF parsing
- `golang.org/x/image/math/fixed` -- Fixed-point math for font metrics
- `golang.org/x/text/unicode/bidi` -- Bidirectional text
