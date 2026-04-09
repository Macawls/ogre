# Ogre

See [AGENTS.md](AGENTS.md) for full project guidelines.

Pure Go HTML/CSS to SVG/PNG renderer for OpenGraph image generation.

## Rules

- Minimal external imports. Prefer stdlib and `golang.org/x/*`. The only allowed `github.com` import is `github.com/go-text/typesetting` for text shaping.
- Idiomatic Go. No C-style APIs, no stuttering names, no unnecessary interfaces.
- No comments in code unless the logic is non-obvious to an AI agent.
- No CGo. Single static binary.
- `<div>` defaults to `display: flex` (matching Satori behavior, not browser behavior).

## Architecture

- `cmd/ogre/` — CLI + HTTP server
- `parse/` — HTML parsing, CSS inline style parsing, node tree
- `style/` — CSS property definitions, shorthand expansion, inheritance, computed values
- `layout/` — Custom flexbox layout engine (W3C spec, not a Yoga port)
- `font/` — Font loading and text measurement using `golang.org/x/image/font`
- `render/` — SVG generation and PNG rasterization
- `server/` — HTTP API, caching, templates

## Dependencies (allowed)

- `golang.org/x/net/html` — HTML parsing
- `golang.org/x/image/font` — Font interfaces
- `golang.org/x/image/font/opentype` — OTF/TTF parsing
- `golang.org/x/image/math/fixed` — Fixed-point math for font metrics
- `golang.org/x/text/unicode/bidi` — Bidirectional text
- `image`, `image/png`, `image/color` — Standard image packages
- `encoding/xml` — SVG output
