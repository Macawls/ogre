---
title: Custom Fonts
description: Loading and using custom fonts in Ogre.
---

Ogre supports multiple font loading strategies. It ships with built-in Go fonts that work out of the box.

## Built-in fonts

Ogre loads Go's default fonts on startup. These render immediately with no configuration.

## Google Fonts

When Ogre encounters a `font-family` in your HTML that it doesn't recognize, it automatically fetches it from Google Fonts. Both regular (400) and bold (700) weights are loaded.

```html
<div style="font-family: Roboto; font-size: 48px; color: white">
  This text uses Roboto from Google Fonts
</div>
```

Fetched fonts are cached on disk in a temporary directory.

## Loading fonts from files

```go
fontData, _ := os.ReadFile("Inter-Regular.ttf")

result, _ := ogre.Render(html, ogre.Options{
    Fonts: []ogre.FontSource{{
        Name:   "Inter",
        Weight: 400,
        Style:  "normal",
        Data:   fontData,
    }},
})
```

## Loading fonts from URLs

```go
result, _ := ogre.Render(html, ogre.Options{
    Fonts: []ogre.FontSource{{
        Name:   "Inter",
        Weight: 400,
        Style:  "normal",
        URL:    "https://example.com/Inter-Regular.ttf",
    }},
})
```

URL fonts are fetched and cached automatically.

## Pre-loading fonts on a shared renderer

```go
r := ogre.NewRenderer()

r.LoadFont(ogre.FontSource{
    Name:   "Inter",
    Weight: 400,
    Style:  "normal",
    Data:   regularData,
})
r.LoadFont(ogre.FontSource{
    Name:   "Inter",
    Weight: 700,
    Style:  "normal",
    Data:   boldData,
})

// All subsequent renders use these fonts
result, _ := r.Render(html, ogre.Options{Width: 1200, Height: 630})
```

## Supported formats

- TTF (TrueType)
- OTF (OpenType)
- WOFF (Web Open Font Format, automatically decompressed)

WOFF2 is not supported. Convert WOFF2 fonts to TTF or OTF before loading them.

## Font rendering

In SVG output, font glyphs are converted to path data. This makes SVGs self-contained — they render correctly without the font installed. In PNG/JPEG output, fonts are rasterized directly.
