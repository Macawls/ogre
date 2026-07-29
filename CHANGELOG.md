# Changelog

All notable changes to this project will be documented in this file.
## [2.1.0] - 2026-07-29

### Performance

- Tune PNG encoder with BestSpeed + BufferPool (encode -22%, -98% bytes/op, 29 -> 0 allocs)
- Drop fmt.Sscanf and []rune conversion in glyph path parser (TextToPath 4.3x faster)
- Cache resolved Tailwind class expansions in a sync.Map (ResolveTailwind -17%, -56% allocs)
- Cache parsed go-text Face on Manager for reuse across shape calls (i18n shape 23x faster)
- Pool alpha buffers used by box-shadow blur (-75% mem per blur)
- LRU-cap the glyph path cache (bounded memory for long-lived servers)
- Use strings.Builder in wrapCollapsed to kill O(n) string concat (wrap -14%, -71% allocs)
- Integer fast path in box blur for radius=1 (~28% faster for thin shadows)
- Direct Pix writes in the opacity composite loop (~27% faster)

### Documentation

- Update changelog for v2.0.1

## [2.0.1] - 2026-07-29

### Bug Fixes

- Percent widths, box-shadow, and img border-radius on flex children (#9)

### Documentation

- Update changelog for v2.0.0

## [2.0.0] - 2026-07-29

### Documentation

- Update changelog for v1.5.1

### Features

- Default <div> to display: block instead of display: flex

## [1.5.1] - 2026-07-22

### Bug Fixes

- Styling, layout, and PNG rendering fixes (#6)

### Documentation

- Update changelog for v1.5.0

## [1.5.0] - 2026-07-14

### Documentation

- Update changelog for v1.4.0

### Features

- Gradient text via background-clip (#5)

## [1.4.0] - 2026-04-12

### Bug Fixes

- Rename changelog to .mdx for import support

### Features

- Automated changelog with git-cliff

### Performance

- Gradient strip precomputation, buffer pooling, parallel blur

## [1.3.0] - 2026-04-12

### Bug Fixes

- Containers grow to fit content when no explicit height
- Text overflow in column-direction flex containers
- Increase emoji-to-text spacing in ship faster cards
- Emoji card layout in ship faster template
- Emoji size in ship faster template, product card layout
- Cache busting for example images, dark mode default, docs cleanup

### Documentation

- Update design goals and add comparison tool commands
- Add Tailwind tab back to playground

### Performance

- SRGB lookup tables for gradient rendering

## [1.2.0] - 2026-04-10

### Documentation

- Fix README dependency claims, add inline SVG and gradient features
- Template redesign with shared data, image viewer, RTL and Tailwind examples
- Fix dependency claims, add rate limit warning, rendering bug guidelines

### Features

- 4:4:4 chroma JPEG encoder
- Tailwind gradient utilities
- SVG rasterization, inline SVGs, and PNG rendering fixes
- Wildcard and multi-origin CORS support

## [1.1.0] - 2026-04-09

### Bug Fixes

- Light mode accent color, simplify CLAUDE.md, docs commit scope
- Remove redundant Docs link from header, reorder nav
- Add site URL for sitemap generation

### Documentation

- Add docs/playground/examples links to README
- Mention hot reload in playground description
- Remove unnecessary qualifier from intro
- Bigger mascot in README

### Features

- Environment variable configuration for server
- Add client-side navigation + external links open in new tab

### Performance

- Aggressive prefetch + client prerendering + fast transitions

### Reverted

- Remove client-side navigation, prefetch, and view transitions

## [1.0.0] - 2026-04-09

### Bug Fixes

- Gitignore was blocking cmd/ogre/ — change 'ogre' to '/ogre'
- Add .node-version for Nixpacks (Node 22 required by Astro)
- Remove gosimple linter (merged into staticcheck in v2)
- Install golangci-lint from source for Go 1.25 compat
- Use golangci-lint v2.11.4 for Go 1.25 support
- Add version field to golangci-lint config for v2 compatibility
- CI — use golangci-lint-action v7, relax perf test timeout for race detector
- Pin golangci-lint to v2.1 for Go 1.25 compatibility

### Documentation

- Add badges to README (tests, release, go reference, GHCR, license)

### Features

- Initial release v1.0.0

