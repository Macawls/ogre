# Changelog

All notable changes to this project will be documented in this file.
## [Unreleased]

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

