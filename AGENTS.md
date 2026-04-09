# Ogre — Agent Guidelines

Instructions for AI agents and contributors working on this codebase.

## Code style

- **Idiomatic Go.** Follow the [Google Go Style Guide](https://google.github.io/styleguide/go/) and [Effective Go](https://go.dev/doc/effective_go).
- **No comments** unless the logic is non-obvious to an AI agent. No docstrings on obvious functions. No commented-out code.
- **No stuttering names.** A type in package `parse` is `Node`, not `ParseNode`.
- **No unnecessary interfaces.** Concrete types unless you need polymorphism.
- **No C-style APIs.** Use Go conventions: multiple return values for errors, zero values as defaults.
- **Error handling:** Return errors, don't panic. Wrap with `fmt.Errorf("context: %w", err)`.
- **Naming:** Short, clear names. Receivers are one or two letters. Local variables are short. Exported names are descriptive.

## Architecture rules

- **Minimal external imports.** Prefer stdlib and `golang.org/x/*`. The only allowed `github.com` import is `github.com/go-text/typesetting` for text shaping.
- **No CGo.** Single static binary.
- **`<div>` defaults to `display: flex`** (matching Satori, not browser behavior).
- Tests go next to the code they test, in the same package.

## Package structure

```
ogre.go          # Public API: Render(), NewRenderer()
handler.go       # HTTP handler: Renderer.Handler()
jsx.go           # JSX-style builder: Div(), Span(), etc.
cmd/ogre/        # CLI entry point
parse/           # HTML parsing → node tree
style/           # CSS properties, Tailwind resolver, inheritance
layout/          # Flexbox layout engine (W3C spec)
font/            # Font loading, text measurement, glyph paths
render/          # SVG generation, PNG/JPEG rasterization
server/          # HTTP server, caching, templates, rate limiting
```

## Testing

- `go test ./...` must pass before any PR.
- Test fixtures in `testdata/fixtures/` for e2e tests.
- Comparison fixtures in `test/fixtures/` for cross-tool accuracy tests.
- No mocks for internal packages. Test real behavior.

## Commits

- Use [Conventional Commits](https://www.conventionalcommits.org/): `feat:`, `fix:`, `perf:`, `docs:`, `test:`, `ci:`, `chore:`.
- Commits that only change docs site content must use `docs(site):` scope.
- GoReleaser generates changelogs from these prefixes on release.

## Docs site

The docs site is at `site/` (Astro Starlight). The WASM playground binary is gitignored and must be built before deploying:

```bash
GOOS=js GOARCH=wasm go build -ldflags="-s -w" \
  -o site/public/playground/ogre.wasm \
  ./cmd/playground/
```

## What not to do

- Don't add features without a concrete use case.
- Don't add abstractions for one-time operations.
- Don't add error handling for scenarios that can't happen.
- Don't add backwards-compatibility shims. Just change the code.
- Don't add `github.com` dependencies. Find a way with stdlib + `golang.org/x/*`.
