---
title: Architecture
description: How Ogre's rendering pipeline works.
---

## Pipeline

Ogre renders HTML to images in four stages:

```
HTML string → Parse → Style → Layout → Render → SVG/PNG/JPEG
```

### 1. Parse

The `parse/` package takes an HTML string and produces a node tree. It uses `golang.org/x/net/html` for parsing. Inline `style` attributes and `class` attributes are extracted from each element.

### 2. Style

The `style/` package resolves all styling:

- Tailwind utility classes are mapped to CSS property values
- CSS shorthand properties are expanded (e.g. `margin: 10px 20px` becomes four individual margin values)
- CSS values are parsed into computed values (colors, lengths, etc.)
- Property inheritance is applied (e.g. `color`, `font-family` inherit from parent)
- Relative units are resolved against the viewport and parent dimensions

### 3. Layout

The `layout/` package implements a custom flexbox layout engine based on the W3C specification. It is not a port of Yoga or any other layout library.

The layout engine:
- Computes the position and size of every node
- Handles flex direction, wrapping, grow/shrink, alignment, and gaps
- Calls back into the font manager to measure text nodes for line breaking

### 4. Render

The `render/` package produces the final output:

- **SVG**: Generates SVG elements with font glyphs converted to `<path>` data. The output is self-contained and renders correctly without fonts installed.
- **PNG**: Rasterizes the layout tree to a pixel buffer with gradient support.
- **JPEG**: Same as PNG, then encoded as JPEG with configurable quality.

## Package structure

```
ogre.go          # Public API: Render(), NewRenderer()
├── parse/       # HTML parsing → node tree
├── style/       # CSS resolution, Tailwind, shorthands, inheritance
├── layout/      # Flexbox layout engine (W3C spec)
├── font/        # Font loading, text measurement, glyph paths
├── render/      # SVG generation, PNG/JPEG rasterization
├── server/      # HTTP API, LRU cache, templates, rate limiting
└── cmd/ogre/    # CLI entry point
```

## Key design decisions

**`<div>` defaults to `display: flex`.** This matches Satori's behavior, not browser behavior. It means every `<div>` is a flex container by default.

**Font glyphs as SVG paths.** In SVG output, text is not rendered as `<text>` elements with font references. Instead, each glyph is converted to path data. This makes SVGs self-contained but larger.

**No Yoga dependency.** The flexbox engine is a from-scratch implementation based on the W3C flexbox specification. This avoids WASM or CGo dependencies.

**Stdlib-only imports.** All external dependencies come from `golang.org/x/*`. No third-party packages.
