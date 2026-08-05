---
title: Introduction
description: 'Ogre is an open-source, pure Go HTML/CSS renderer for OpenGraph images and presentation slides. Convert HTML and CSS to SVG, PNG, and JPEG from a single static binary — usable as a Go library, CLI, or HTTP server. No Chrome, no Node, no CGo.'
head:
  - tag: meta
    attrs:
      name: keywords
      content: 'ogre, golang, go, opengraph, og image, html to image, html to svg, html to png, html to slides, presentation slides, pure go, no cgo, satori alternative'
  - tag: meta
    attrs:
      property: 'og:title'
      content: 'Ogre — Pure Go HTML/CSS to SVG, PNG, and JPEG'
  - tag: meta
    attrs:
      property: 'og:description'
      content: 'Open-source, pure Go renderer for HTML and CSS. Generate OpenGraph images and presentation slides as SVG, PNG, or JPEG from a single static binary.'
---

Ogre is a pure Go renderer that converts HTML and CSS into SVG, PNG, and JPEG images. It compiles to a single static binary and works as a Go library, CLI, or HTTP server — designed for generating OpenGraph images, social cards, presentation slides, and dynamic image content from templates.

## What is this for?

When you share a link on Twitter, Slack, Discord, or LinkedIn, a preview image appears. That image is an OpenGraph (OG) image. Instead of designing a static image for every page, you can generate them from an HTML template.

Ogre handles this. Write HTML with inline styles or Tailwind classes, pass it to Ogre, get an image back. Use it as a Go library, a standalone CLI, or a self-hosted HTTP server. Blog post cards, documentation pages, event banners, repo cards, presentation slides, or any dynamic image your application needs.

The canvas size is arbitrary — defaults are 1200×630 (the OG standard), but set `Width` and `Height` to whatever you need. For example, `1920×1080` for HD presentation slides, or a custom size for a certificate or invoice.

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

Dynamic image generation (OG cards, social previews, certificates, invoices, presentation slides) usually pulls in a separate runtime — a JavaScript service or a headless browser. Ogre keeps it inside one Go binary:

- **One binary, three surfaces.** The same code is a Go library, a CLI, and an HTTP server. `go get` it, drop it into an existing `net/http` mux, run it as a subprocess, or deploy the Docker image.
- **Low overhead.** No Chrome, no Node, no CGo — `CGO_ENABLED=0` produces a static binary that runs in single-digit megabytes of RAM instead of the hundreds a headless browser needs.
- **Agent-friendly.** One-command install, JSON HTTP API, `llms.txt` + `llms-full.txt` served alongside the docs, and a [public instance](https://ogre-api.macawls.dev) for experiments. Feed HTML in, get bytes out — no browser fleet to manage.
- **Tailwind built-in.** Utility classes resolve directly (v3 default, v4 opt-in). No build step, no Tailwind CLI.

For a side-by-side with Satori — feature matrix, benchmark numbers, and pixel-accuracy scores — see the [Satori Comparison](/advanced/satori-comparison/) page.

## Design goals

- **Pure Go.** No CGo, no external binaries. Single static binary with `CGO_ENABLED=0`.
- **Output quality first.** The priority is correct, complete rendering — PNG/JPEG output, inline SVGs, box shadows, CSS filters, transforms, RTL text. Performance optimizations follow once output fidelity is solid.
- **Tailwind built-in.** Resolves Tailwind v3 (default) or v4 utility classes directly. No build step needed.
- **Production-ready server.** Includes an HTTP server with LRU caching, rate limiting, and template support.

## Dependencies

Standard library, `golang.org/x/*`, and one external package:

- `golang.org/x/net/html` for HTML parsing
- `golang.org/x/image/font` for font interfaces and rasterization
- `golang.org/x/image/vector` for 2D vector path rasterization
- `golang.org/x/text/unicode/bidi` for bidirectional text
- `github.com/go-text/typesetting` for text shaping (kerning, ligatures, RTL)

## Output formats

| Format | Content Type | Notes |
|--------|-------------|-------|
| SVG | `image/svg+xml` | Font glyphs embedded as path data. Self-contained. |
| PNG | `image/png` | Rasterized with gradient support. |
| JPEG | `image/jpeg` | Configurable quality (default 90). |
