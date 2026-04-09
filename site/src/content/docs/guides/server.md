---
title: HTTP Server
description: Running Ogre as an HTTP service.
---

Ogre includes a production-ready HTTP server with caching, rate limiting, and template support.

## Starting the server

### Via CLI

```bash
ogre --serve --port 3000
```

### Programmatically

```go
srv := server.New(server.Config{
    Addr:          ":3000",
    CacheBytes:    64 << 20, // 64 MB LRU cache
    RateLimit:     10,       // 10 req/s per IP
    RenderTimeout: 10 * time.Second,
    MaxElements:   1000,
})

if err := srv.Start(); err != nil {
    log.Fatal(err)
}
```

## Endpoints

### POST /render

Render HTML to an image.

```bash
curl -X POST http://localhost:3000/render \
  -H "Content-Type: application/json" \
  -d '{
    "html": "<div class=\"flex w-full h-full bg-blue-500 items-center justify-center\"><div class=\"text-4xl font-bold text-white\">Hello</div></div>",
    "width": 1200,
    "height": 630,
    "format": "svg"
  }' \
  -o output.svg
```

Request fields:

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `html` | string | required | HTML to render |
| `width` | int | 1200 | Canvas width |
| `height` | int | 630 | Canvas height |
| `format` | string | `"svg"` | `"svg"`, `"png"`, or `"jpeg"` |
| `fonts` | array | `[]` | Custom fonts (see below) |

### POST /render/template

Render a Go `html/template` with data substitution.

```bash
curl -X POST http://localhost:3000/render/template \
  -H "Content-Type: application/json" \
  -d '{
    "template": "<div class=\"flex flex-col w-full h-full bg-slate-900 p-16 justify-center\"><div class=\"text-5xl font-bold text-white\">{{.Title}}</div><div class=\"text-xl text-slate-400 mt-4\">{{.Subtitle}}</div></div>",
    "data": {"Title": "Hello World", "Subtitle": "From a template"},
    "width": 1200,
    "height": 630,
    "format": "png"
  }' \
  -o output.png
```

### GET /health

Returns `{"status":"ok"}`.

### GET /metrics

Returns render statistics:

```json
{
  "render_total": 150,
  "render_errors": 2,
  "cache_hits": 98,
  "cache_misses": 52,
  "total_duration_ms": 4200
}
```

## Custom fonts via API

Include fonts in the render request:

```json
{
  "html": "<div style=\"font-family: Inter\">Hello</div>",
  "fonts": [
    {
      "name": "Inter",
      "weight": 400,
      "style": "normal",
      "url": "https://example.com/Inter-Regular.ttf"
    }
  ]
}
```

Fonts can be provided as a URL or base64-encoded data. Limits: 5 fonts per request, 5 MB per font.

## Caching

Responses are cached in an LRU cache keyed by SHA-256 of the input (HTML + dimensions + format). Cache headers:

- `ETag`: SHA-256 hash of the input
- `X-Cache`: `HIT` or `MISS`
- `Cache-Control`: `public, max-age=86400` for cache hits, `public, max-age=3600` for misses

## CORS

All endpoints return `Access-Control-Allow-Origin: *`.

## Limits

- Max request body: 10 MB
- Max fonts per request: 5
- Max font size: 5 MB
- Default cache: 64 MB
- Default render timeout: 10 seconds
- Default max elements: 1000
