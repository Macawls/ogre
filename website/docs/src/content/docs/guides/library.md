---
title: Go Library
description: Embedding Ogre in your Go application.
---

## One-shot render

The simplest way to render. Creates a new font manager each call.

```go
result, err := ogre.Render(`
    <div class="flex flex-col w-full h-full bg-slate-900 p-16 justify-center">
        <div class="text-5xl font-bold text-white">Hello World</div>
    </div>
`, ogre.Options{
    Width:  1200,
    Height: 630,
})
if err != nil {
    log.Fatal(err)
}

os.WriteFile("og.svg", result.Data, 0644)
```

## Shared renderer

For servers and applications that render multiple images, use a shared renderer. It reuses the font manager across renders and is thread-safe.

```go
r := ogre.NewRenderer()

// Use across multiple goroutines
result, err := r.Render(html, ogre.Options{
    Width:  1200,
    Height: 630,
    Format: ogre.FormatPNG,
})
```

## Options

```go
type Options struct {
    Width         int          // Canvas width (default 1200)
    Height        int          // Canvas height (default 630)
    Format        Format       // "svg" (default), "png", or "jpeg"
    Quality       int          // JPEG quality 1-100 (default 90)
    Fonts         []FontSource // Custom fonts to load for this render
    Debug         bool         // Enable debug output
    EmojiProvider string       // "twemoji" (default) or "none"
    MaxElements   int          // Maximum HTML element count (0 = unlimited)
}
```

## Result

```go
type Result struct {
    Data        []byte // Rendered image bytes
    ContentType string // MIME type ("image/svg+xml", "image/png", "image/jpeg")
    Width       int
    Height      int
}
```

## Format constants

```go
ogre.FormatSVG  // "svg"
ogre.FormatPNG  // "png"
ogre.FormatJPEG // "jpeg"
```

## Integration with net/http

The `Handler` method returns an `http.Handler` that accepts JSON POST requests and returns rendered images. It reuses the renderer's font manager across requests.

```go
r := ogre.NewRenderer()

mux := http.NewServeMux()
mux.Handle("POST /og", r.Handler(ogre.HandlerConfig{
    Width:  1200,
    Height: 630,
    Format: ogre.FormatPNG,
}))

log.Fatal(http.ListenAndServe(":8080", mux))
```

The handler accepts the same JSON body as the [HTTP server](/guides/server/):

```bash
curl -X POST http://localhost:8080/og \
  -H "Content-Type: application/json" \
  -d '{
    "html": "<div class=\"flex w-full h-full bg-slate-900 items-center justify-center\"><div class=\"text-5xl font-bold text-white\">Hello</div></div>"
  }' \
  -o og.png
```

Request fields (`width`, `height`, `format`) override the defaults from `HandlerConfig` when provided. The handler also supports `template` + `data` fields for Go template rendering.

### HandlerConfig

```go
type HandlerConfig struct {
    Width  int    // Default canvas width (default 1200)
    Height int    // Default canvas height (default 630)
    Format Format // Default output format (default FormatPNG)
}
```
