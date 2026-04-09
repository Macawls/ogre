---
title: Emoji & RTL
description: Emoji rendering with Twemoji/OpenMoji providers and RTL text support.
---

## Emoji

Ogre renders emoji in both SVG and PNG output by fetching images from an emoji CDN. Text is automatically split into emoji and non-emoji segments — emoji segments render as images, text segments render as normal font glyphs.

```html
<div style="display:flex;flex-direction:column;width:100%;height:100%;
  background-color:#09090b;padding:64px;justify-content:center">
  <div style="color:white;font-size:48px;font-weight:700">🚀 Ogre v1.0</div>
  <div style="color:#a1a1aa;font-size:24px;margin-top:16px">
    Pure Go HTML to Image ✨ No Chrome 🔥 No Node 💪</div>
</div>
```

| SVG | PNG | JPEG |
|-----|-----|------|
| ![SVG](/examples/emoji-card.svg) | ![PNG](/examples/emoji-card.png) | ![JPEG](/examples/emoji-card.jpg) |

### Providers

Ogre supports multiple emoji providers:

| Provider | Style | Format |
|----------|-------|--------|
| `twemoji` (default) | Twitter/X style | SVG + PNG |
| `openmoji` | OpenMoji style | SVG + PNG |
| `noto` | Google Noto Emoji | SVG only |

In SVG output, emoji render as `<image href="cdn-url">`. In PNG/JPEG output, emoji PNGs are fetched from the CDN, scaled to match the font size, and composited into the image.

### Configuring the provider

```go
// Default (Twemoji)
result, _ := ogre.Render(html, ogre.Options{})

// Disable emoji rendering
result, _ := ogre.Render(html, ogre.Options{
    EmojiProvider: "none",
})
```

### How it works

1. `font.SplitEmoji()` splits text into alternating text/emoji runs
2. Text runs render via the normal font glyph path
3. Emoji runs are resolved to CDN URLs via `font.EmojiSVGURL()` or `font.EmojiPNGURL()`
4. SVG: emoji become `<image href="...">` elements
5. PNG/JPEG: emoji PNGs are fetched, cached, scaled with bilinear interpolation, and composited

Emoji images are cached in memory for the lifetime of the `EmojiProvider` instance.

## RTL text

Use `direction: rtl` for right-to-left text. Ogre applies the Unicode bidirectional algorithm (`golang.org/x/text/unicode/bidi`) to reorder characters and automatically right-aligns text.

```html
<div style="direction:rtl;color:white;font-size:32px">
  RTL text renders right-aligned
</div>
```

### What works

- Bidi character reordering via the Unicode bidi algorithm
- Automatic right-alignment when `direction: rtl` is set
- Mixed LTR/RTL content in the same line
- `direction` property inherits to child elements

### Text shaping

Ogre uses [`go-text/typesetting`](https://github.com/go-text/typesetting) — a pure Go port of HarfBuzz — for complex script rendering. When text contains Arabic, Hebrew, Devanagari, Thai, or other complex scripts, Ogre automatically uses the shaper for proper glyph substitution (GSUB) and positioning (GPOS).

In SVG output, shaped text renders as glyph paths with correct ligatures and connected letterforms. In PNG/JPEG output, individual glyphs render correctly but connected forms may vary depending on the font.

Load an appropriate font for the script you're rendering. Ogre auto-fetches from Google Fonts:

```html
<div style="font-family: Noto Sans Arabic; direction: rtl; color: white; font-size: 48px">
  مرحبا بالعالم
</div>
```
