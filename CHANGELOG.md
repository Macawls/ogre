# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-04-09

### Added

- HTML + inline CSS to SVG rendering with embedded font paths
- HTML + inline CSS to PNG rendering with gradient support
- HTML + inline CSS to JPEG rendering
- Tailwind CSS v3 utility class support (colors, spacing, typography, flexbox, borders, effects, arbitrary values)
- Custom flexbox layout engine (W3C spec compliant with gap support)
- Font embedding as SVG path elements for self-contained output
- Built-in Go fonts (Regular + Bold) for zero-config usage
- Google Fonts auto-resolution (40+ popular fonts)
- CDN font loading by URL with disk and memory caching
- Custom font loading (TTF, OTF, WOFF)
- Emoji detection and Twemoji CDN integration
- Text overflow ellipsis and line-clamp support
- CSS variable support (--var and var())
- CSS filter support (blur, grayscale, brightness, contrast, saturate, sepia, hue-rotate, invert, drop-shadow)
- Box shadow rendering (outset and inset) in SVG and PNG
- Linear and radial gradient support in SVG and PNG
- Border rendering with radius, styles (solid, dashed, dotted, double)
- CSS transform support (translate, scale, rotate, skew)
- Overflow hidden with clip-path
- Opacity support
- HTTP server with LRU caching, rate limiting, CORS, and structured logging
- Go template support in server (POST /render/template)
- Per-request font loading via URL in server API
- CLI tool (--serve, --render, --html, --output, --format)
- Prometheus-compatible /metrics endpoint
- Dockerfile for deployment
- Shared Renderer struct for concurrent-safe reuse across requests
- Glyph path caching for text-heavy renders
- Font face caching
- Image fetch caching with timeout
- 25 Satori comparison test fixtures with pixel-by-pixel PNG verification
- 25 Takumi comparison references
- 28 edge case tests (empty input, deep nesting, concurrent rendering, malformed HTML)
- Custom font integration tests
- Performance benchmarks (Ogre vs Satori vs Takumi)
