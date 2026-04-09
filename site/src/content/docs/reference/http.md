---
title: HTTP Endpoints
description: Complete reference for the Ogre HTTP API.
---

## POST /render

Render HTML to an image.

### Request

```json
{
  "html": "<div class=\"flex w-full h-full bg-blue-500\">Hello</div>",
  "width": 1200,
  "height": 630,
  "format": "svg",
  "fonts": []
}
```

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `html` | string | yes | | HTML to render |
| `width` | int | no | 1200 | Canvas width |
| `height` | int | no | 630 | Canvas height |
| `format` | string | no | `"svg"` | `"svg"`, `"png"`, or `"jpeg"` |
| `quality` | int | no | 90 | JPEG compression quality (1-100). Ignored for SVG/PNG. |
| `fonts` | array | no | `[]` | Custom fonts |

### Font objects

```json
{
  "name": "Inter",
  "weight": 400,
  "style": "normal",
  "url": "https://example.com/Inter.ttf"
}
```

Or with base64 data:

```json
{
  "name": "Inter",
  "weight": 400,
  "style": "normal",
  "data": "base64-encoded-font-data"
}
```

### Response

Image bytes with headers:

| Header | Value |
|--------|-------|
| `Content-Type` | `image/svg+xml`, `image/png`, or `image/jpeg` |
| `ETag` | SHA-256 hash of the input |
| `X-Cache` | `HIT` or `MISS` |
| `Cache-Control` | `public, max-age=86400` (hit) or `public, max-age=3600` (miss) |

## POST /render/template

Render a Go `html/template` with data substitution.

### Request

```json
{
  "template": "<div class=\"text-4xl text-white\">{{.Title}}</div>",
  "data": { "Title": "Hello World" },
  "width": 1200,
  "height": 630,
  "format": "svg"
}
```

| Field | Type | Required | Default |
|-------|------|----------|---------|
| `template` | string | yes | |
| `data` | object | no | `{}` |
| `width` | int | no | 1200 |
| `height` | int | no | 630 |
| `format` | string | no | `"svg"` |
| `quality` | int | no | 90 |

### Response

Same as `/render`.

## GET /health

```json
{"status": "ok"}
```

## GET /metrics

```json
{
  "render_total": 150,
  "render_errors": 2,
  "cache_hits": 98,
  "cache_misses": 52,
  "total_duration_ms": 4200
}
```

## Error responses

All errors return JSON:

```json
{
  "error": "description of what went wrong"
}
```

| Status | Meaning |
|--------|---------|
| 400 | Bad request (missing fields, invalid template, invalid font) |
| 429 | Rate limit exceeded |
| 500 | Render failed |
| 504 | Render timeout |

## CORS

All endpoints return `Access-Control-Allow-Origin: *`.

## Limits

| Limit | Value |
|-------|-------|
| Max request body | 10 MB |
| Max fonts per request | 5 |
| Max font size | 5 MB |
| Default cache size | 64 MB |
| Default render timeout | 10 seconds |
| Default max elements | 1000 |
