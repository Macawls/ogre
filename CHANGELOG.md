# Changelog

All notable changes to this project will be documented in this file.
## [3.2.1] - 2026-08-08

### Bug Fixes

- Use v3 module path (#17)

### CI

- Chain changelog job after release in a single workflow
- Fetch ogre.wasm from GitHub releases when source is absent
- Add Go to Nixpacks setup phase

### Documentation

- Update changelog for v3.2.0

## [3.2.0] - 2026-08-07

### Documentation

- Add WebAssembly guide with live demo
- Make install page easier to skim
- Per-OS install sections and soften intro identity
- Add dedicated Performance page under Reference
- Refresh Satori comparison with re-measured benchmarks
- Sync API contracts with code
- Reposition around Ogre's own strengths and add slides use case
- Update changelog for v3.1.0

### Features

- Ship ogre.wasm as release artifact

## [3.1.0] - 2026-07-30

### Documentation

- Update changelog for v3.1.0
- Document font-variation-settings + variable font support
- Add Images guide covering <img> and background-image usage
- Update changelog for v3.0.0

### Features

- Respect font-variation-settings on variable fonts (#12)
- Add Tailwind v4 opt-in via Options.TailwindVersion

### Performance

- Call GlyphDataOutline directly, add variation benches

## [3.0.0] - 2026-07-30

### Documentation

- Update changelog for v2.1.0

### Refactor

- Remove JSX-style builder API (#13)

## [2.1.0] - 2026-07-30

### Documentation

- Auto-sync changelog from CHANGELOG.md, remove Examples page (#10)
- Update changelog for v2.0.1

### Performance

- V2.1.0 batch — 9 measured perf improvements across render, font, and style (#11)

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

