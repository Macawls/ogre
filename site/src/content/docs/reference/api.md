---
title: Go API Reference
description: Complete reference for the Ogre Go package.
---

## Functions

### Render

```go
func Render(html string, opts Options) (*Result, error)
```

One-shot render. Creates a new font manager, loads default fonts, renders the HTML, and returns the result. Convenient for single renders. For repeated renders, use `NewRenderer`.

### NewRenderer

```go
func NewRenderer() *Renderer
```

Creates a `Renderer` with default fonts loaded. The renderer is thread-safe and should be reused across renders in server applications.

## Renderer methods

### Render

```go
func (r *Renderer) Render(html string, opts Options) (*Result, error)
```

Renders HTML to an image using the renderer's shared font manager. Thread-safe.

### LoadFont

```go
func (r *Renderer) LoadFont(src FontSource) error
```

Registers a font for use in subsequent renders. If `FontSource.URL` is set and `Data` is empty, the font is fetched from the URL and cached on disk.

### Handler

```go
func (r *Renderer) Handler(cfg HandlerConfig) http.Handler
```

Returns an `http.Handler` that accepts JSON POST requests and renders images. Uses the renderer's shared font manager. The handler accepts `html`, `template`/`data`, `width`, `height`, and `format` fields in the JSON body. Fields in the request override the defaults from `HandlerConfig`.

```go
type HandlerConfig struct {
    Width   int    // Default canvas width. Default: 1200.
    Height  int    // Default canvas height. Default: 630.
    Format  Format // Default output format. Default: FormatPNG.
    Quality int    // Default JPEG quality 1-100. Default: 90.
}
```

## Types

### Options

```go
type Options struct {
    Width         int          // Canvas width in pixels. Default: 1200.
    Height        int          // Canvas height in pixels. Default: 630.
    Format        Format       // Output format. Default: FormatSVG.
    Quality       int          // JPEG quality 1-100. Default: 90. Only used with FormatJPEG.
    Fonts         []FontSource // Fonts to load for this render.
    Debug         bool         // Enable debug logging.
    EmojiProvider string       // "twemoji" (default) or "none".
    MaxElements   int          // Max HTML elements allowed. 0 = unlimited.
}
```

### FontSource

```go
type FontSource struct {
    Name   string // Font family name (e.g. "Inter")
    Weight int    // Font weight: 100-900
    Style  string // "normal" or "italic"
    Data   []byte // Raw font bytes (TTF, OTF, or WOFF)
    URL    string // URL to fetch font from. Used if Data is empty.
}
```

### Result

```go
type Result struct {
    Data        []byte // Rendered image bytes
    ContentType string // MIME type: "image/svg+xml", "image/png", or "image/jpeg"
    Width       int    // Canvas width used
    Height      int    // Canvas height used
}
```

### Format

```go
type Format string

const (
    FormatSVG  Format = "svg"
    FormatPNG  Format = "png"
    FormatJPEG Format = "jpeg"
)
```

## Server package

### server.New

```go
func New(cfg Config) *Server
```

Creates a server with routes registered.

### server.Config

```go
type Config struct {
    Addr          string        // Listen address. Default: ":8080".
    CacheBytes    int64         // LRU cache size in bytes. Default: 64 MB.
    Fonts         []ogre.FontSource // Pre-loaded fonts for all renders.
    RateLimit     float64       // Requests per second per IP. 0 = no limit.
    RenderTimeout time.Duration // Per-render timeout. Default: 10s.
    MaxElements   int           // Max HTML elements per render. Default: 1000.
}
```

### Server methods

```go
func (s *Server) Start() error    // Listen and serve until SIGINT/SIGTERM.
func (s *Server) Handler() http.Handler // Returns the HTTP handler with CORS.
```
