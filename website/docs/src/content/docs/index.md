---
title: Introduction
description: 'Ogre is an open-source, pure Go alternative to Vercel Satori for generating OpenGraph images from HTML and CSS. Convert HTML to SVG, PNG, and JPEG with zero dependencies. No Chrome, no Node, no CGo.'
head:
  - tag: meta
    attrs:
      name: keywords
      content: 'ogre, satori alternative, vercel satori, golang, go, opengraph, og image, html to image, html to svg, html to png, pure go, no cgo'
  - tag: meta
    attrs:
      property: 'og:title'
      content: 'Ogre — Vercel Satori Alternative in Go'
  - tag: meta
    attrs:
      property: 'og:description'
      content: 'Open-source, pure Go alternative to Vercel Satori. Convert HTML and CSS to SVG, PNG, and JPEG images. Single static binary.'
---

Ogre is a pure Go alternative to [Vercel's Satori](https://github.com/vercel/satori) for converting HTML and CSS into SVG, PNG, and JPEG images. It is designed for generating OpenGraph images, social cards, and dynamic image content from templates.

## What is this for?

When you share a link on Twitter, Slack, Discord, or LinkedIn, a preview image appears. That image is an OpenGraph (OG) image. Instead of designing a static image for every page, you can generate them from an HTML template.

Ogre handles this. Write HTML with inline styles or Tailwind classes, pass it to Ogre, get an image back. Use it as a Go library, a standalone CLI, or a self-hosted HTTP server — no Go knowledge required for the latter two. Blog post cards, documentation pages, event banners, repo cards, or any dynamic image your application needs.

```bash
go install github.com/macawls/ogre/cmd/ogre@latest
```

## Quick look

```go
result, _ := ogre.Render(`
    <div class="flex w-full h-full bg-slate-900 p-16 items-center justify-center">
        <div class="text-5xl font-bold text-white">Hello World</div>
    </div>
`, ogre.Options{Width: 1200, Height: 630})

os.WriteFile("og.svg", result.Data, 0644)
```

## Why this exists

Most web applications need dynamic image generation at some point — OG cards, social previews, certificates, invoices. The existing options all have trade-offs:

- **Satori** requires a JavaScript runtime. If your backend is Go, you're running a separate service just for images.
- **Headless Chrome / Puppeteer** is heavy. A Chrome process uses hundreds of megabytes of RAM and takes seconds to render a single image.
- **Image manipulation libraries** (like Go's `image` package) work but you lose the ability to use HTML/CSS for layout.

Ogre was built to solve this for Go applications. Add it as a dependency and call `ogre.Render()` — no sidecar service, no runtime, no external process. Or run it standalone as a server behind an internal endpoint.

| | Ogre | Satori |
|---|---|---|
| Binary | 11 MB static binary | JavaScript runtime required |
| Simple render | 0.03–0.08 ms | 0.3–2.5 ms |
| Complex render | 3–8 ms | 4–17 ms |
| Output | SVG, PNG, JPEG | SVG only |
| Use as library | `go get`, one function | npm package |

Render times measured on AMD Ryzen 5 5600H, 1200x630 renders, both producing SVG. Full benchmark data in [Satori Comparison](/advanced/satori-comparison/).

## Design goals

- **Pure Go.** No CGo, no external binaries, no runtime dependencies. Adding Ogre won't force CGo on your project. Single static binary with `CGO_ENABLED=0`.
- **Drop-in Satori alternative.** Accepts the same HTML/CSS subset as Satori but runs natively in Go.
- **Tailwind built-in.** Resolves Tailwind v3 utility classes directly. No build step needed.
- **Production-ready server.** Includes an HTTP server with LRU caching, rate limiting, and template support.

## Dependencies

Only standard library and `golang.org/x/*` packages:

- `golang.org/x/net/html` for HTML parsing
- `golang.org/x/image/font` for font interfaces and OpenType parsing
- `golang.org/x/text/unicode/bidi` for bidirectional text

No third-party imports.

## Output formats

| Format | Content Type | Notes |
|--------|-------------|-------|
| SVG | `image/svg+xml` | Font glyphs embedded as path data. Self-contained. |
| PNG | `image/png` | Rasterized with gradient support. |
| JPEG | `image/jpeg` | Configurable quality (default 90). |
