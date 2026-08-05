---
title: Performance
description: Ogre's render times, memory characteristics, and how to reproduce the numbers on your own hardware.
---

Every render time and memory number on this page comes from `go test -bench` on the fixtures in `test/fixtures/`. The commands are at the bottom — if the numbers matter for your decision, run them on the hardware you actually plan to deploy on.

## Methodology

- Hardware: AMD Ryzen 5 5600H
- Canvas: 1200×630 (the OpenGraph standard)
- Fixture set: the same 25 files used for pixel-accuracy verification, ranging from atomic layout tests (`01-solid-bg`) through realistic composed cards (`20-og-card-real`, `21-github-repo-card`, etc.)
- Warm-up: Go's `testing` framework runs each benchmark `b.N` times to convergence
- Renderer: one-shot `ogre.Render()` per iteration (creates a new font manager each call). Real applications using `ogre.NewRenderer()` amortize font-manager setup across calls and are faster in aggregate

## Render times

| Fixture | SVG (ms) | PNG (ms) |
|---------|---------:|---------:|
| 01-solid-bg          | 0.027 | 5.6  |
| 02-flex-row-grow     | 0.060 | 18.6 |
| 03-flex-column       | 0.057 | 18.5 |
| 04-text-basic        | 0.160 | 6.8  |
| 05-text-multiline    | 0.700 | 8.2  |
| 06-linear-gradient   | 0.030 | 40.7 |
| 07-border-radius     | 0.042 | 10.2 |
| 08-nested-flex       | 0.065 | 19.8 |
| 09-justify-center    | 0.034 | 7.6  |
| 10-space-between     | 0.052 | 7.1  |
| 11-padding-margin    | 0.040 | 17.0 |
| 12-opacity           | 0.036 | 13.8 |
| 13-text-colors       | 0.317 | 7.6  |
| 14-border-styles     | 0.538 | 10.8 |
| 15-box-shadow        | 0.044 | 18.0 |
| 16-absolute-position | 0.045 | 6.9  |
| 17-overflow-hidden   | 0.049 | 11.9 |
| 18-transform         | 0.038 | 7.4  |
| 19-radial-gradient   | 0.032 | 46.0 |
| 20-og-card-real      | 1.73  | 42.8 |
| 21-github-repo-card  | 1.68  | 42.4 |
| 22-blog-minimal      | 1.80  | 28.4 |
| 23-product-pricing   | 0.84  | 40.9 |
| 24-event-banner      | 1.53  | 9.7  |
| 25-dashboard-stat    | 1.08  | 20.9 |

**SVG summary:** atomic fixtures render in 0.03–0.70 ms, composed OG-card fixtures in 0.84–1.80 ms. The outliers on the low end (`06-linear-gradient`, `19-radial-gradient` at 0.03 ms) are just background rectangles with a gradient fill; the outliers on the high end within the atomic set (`05-text-multiline`, `13-text-colors`, `14-border-styles`) do more text shaping.

**PNG summary:** 5.6–46 ms per fixture. PNG is dominated by rasterization — every gradient sample, every glyph outline, every shadow blur is computed per-pixel. Composed cards land in the 20–43 ms range. If you need PNG at higher throughput than one render per ~25 ms, cache the output (the built-in HTTP server does this automatically via LRU).

## Memory

- SVG renders allocate roughly 32 KB (simple fixtures) to 890 KB (composed OG cards) per call, driven by the glyph-path buffers.
- PNG renders allocate 3–6 MB per call — a 1200×630 RGBA framebuffer alone is 3 MB.
- The one-shot `ogre.Render()` path allocates a fresh font manager per call. Switching to `ogre.NewRenderer()` reuses fonts across calls and eliminates that overhead.

Full per-benchmark `B/op` and `allocs/op` numbers are printed by `go test -bench=. -benchmem` — see the commands below.

## Reproducing these numbers

```bash
cd test
go test -bench=BenchmarkRenderSVG -benchmem -benchtime=3s -run=^$
go test -bench=BenchmarkRenderPNG -benchmem -benchtime=3s -run=^$
```

The `-benchtime=3s` flag gives each fixture 3 seconds of iteration for stable numbers on a mixed-load machine; drop it if you want faster (noisier) results.

Additional benchmark suites live in the sub-packages:

```bash
go test -bench=. ./font/    # text shaping and glyph extraction
go test -bench=. ./render/  # SVG serialization and PNG rasterization
go test -bench=. ./style/   # Tailwind resolution and shorthand expansion
```

## Comparison with other tools

For side-by-side numbers against Satori (and pixel-accuracy scores across the same 25 fixtures), see the [Satori Comparison](/advanced/satori-comparison/) page.
