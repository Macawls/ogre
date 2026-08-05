---
title: Satori Comparison
description: How Ogre compares to Vercel's Satori — feature matrix, render performance, and pixel accuracy.
---

Ogre and [Satori](https://github.com/vercel/satori) solve the same problem — HTML/CSS to image — in different ecosystems. Satori runs in TypeScript on a Node.js or Bun runtime. Ogre is a pure-Go binary. This page collects the concrete differences.

## Feature comparison

| Feature | Ogre | Satori |
|---------|------|--------|
| Language | Go | TypeScript |
| Output formats | SVG, PNG, JPEG | SVG only |
| Dependencies | stdlib + golang.org/x + go-text | yoga-wasm + others |
| Deployment | Single static binary | Node.js runtime |
| Tailwind support | Built-in (v3 default, v4 opt-in) | Via plugin |
| PNG output | Built-in | Requires resvg |
| Layout engine | Custom flexbox (W3C) | Yoga (via WASM) |
| HTTP server | Built-in with caching | BYO |
| Font embedding | SVG paths | SVG paths |
| Emoji | Twemoji CDN | Twemoji CDN |
| `<div>` default | `display: block` | `display: flex` |
| Pixel accuracy | 98%+ vs Satori | Reference |

## Render performance

Measured on AMD Ryzen 5 5600H, 1200×630 SVG output, 25 fixtures in `test/fixtures/`. Ogre numbers from `go test -bench=BenchmarkRenderSVG`; Satori numbers from `test/satori-reference/bench.ts` (20 iterations per fixture, warm renderer).

**Summary:** Ogre is faster than Satori on every fixture. Median ratio ~5×, range 2.2× to 26×.

| Fixture | Ogre (ms) | Satori (ms) | Ogre vs Satori |
|---------|----------:|------------:|---------------:|
| 01-solid-bg | 0.027 | 0.7 | 26× |
| 02-flex-row-grow | 0.060 | 0.9 | 15× |
| 03-flex-column | 0.057 | 0.5 | 8.8× |
| 04-text-basic | 0.160 | 1.2 | 7.5× |
| 05-text-multiline | 0.700 | 2.6 | 3.7× |
| 06-linear-gradient | 0.030 | 0.3 | 10× |
| 07-border-radius | 0.042 | 0.4 | 9.5× |
| 08-nested-flex | 0.065 | 0.5 | 7.7× |
| 09-justify-center | 0.034 | 0.3 | 8.8× |
| 10-space-between | 0.052 | 0.4 | 7.7× |
| 11-padding-margin | 0.040 | 0.2 | 5.0× |
| 12-opacity | 0.036 | 0.2 | 5.6× |
| 13-text-colors | 0.317 | 1.4 | 4.4× |
| 14-border-styles | 0.538 | 1.6 | 3.0× |
| 15-box-shadow | 0.044 | 0.2 | 4.5× |
| 16-absolute-position | 0.045 | 0.3 | 6.7× |
| 17-overflow-hidden | 0.049 | 0.3 | 6.1× |
| 18-transform | 0.038 | 0.2 | 5.3× |
| 19-radial-gradient | 0.032 | 0.2 | 6.3× |
| 20-og-card-real | 1.73 | 4.7 | 2.7× |
| 21-github-repo-card | 1.68 | 4.2 | 2.5× |
| 22-blog-minimal | 1.80 | 4.3 | 2.4× |
| 23-product-pricing | 0.84 | 2.7 | 3.2× |
| 24-event-banner | 1.53 | 3.4 | 2.2× |
| 25-dashboard-stat | 1.08 | 2.7 | 2.5× |

Fixtures 01–19 are atomic layout/typography/effect tests. Fixtures 20–25 are realistic OG-card compositions. At OG-image scale (a few requests per second) both tools are fast enough that speed rarely decides the tool choice — the more useful signal is that Ogre stays consistent on composed layouts (~1–2 ms), where JavaScript runtimes tend to show more variance.

The Ogre benchmark calls `ogre.Render()` (one-shot API — a fresh font manager per iteration). Satori is measured with fonts pre-loaded. Ogre's numbers therefore include per-call setup that Satori's don't, so the ratio above is a conservative view of Ogre's throughput. Using `ogre.NewRenderer()` for a shared, thread-safe renderer removes that overhead.

## Pixel accuracy

All 25 fixtures pass a pixel-diff check against Satori's SVG output. Worst-case fixture is `23-product-pricing` at 98.6% match; 13 of 25 fixtures are exact matches. Scores are stored in `test/output/scores.json` and updated as the fixtures change.

The remaining sub-2% differences are sub-pixel positioning that comes from the different layout engines: Satori delegates to Yoga (C++ compiled to WASM), Ogre implements the W3C flexbox spec directly in Go.

## When to use Ogre

- You are building a Go backend and want HTML→image without a Node process
- You want PNG or JPEG directly, without piping SVG through resvg
- You want one binary you can drop into a container, a CLI, or a Go `net/http` mux
- You want built-in Tailwind (v3 or v4) with no build step
- You want an HTTP server with caching, rate limiting, and templates that runs out of the box

## When to use Satori

- You are already running Node.js, Bun, or a Next.js stack
- You are using `@vercel/og` in Next.js and want to stay on the reference path
- You need a specific Satori-only behavior your templates depend on

## Compatibility

Ogre accepts the same HTML/CSS subset as Satori. Templates written for Satori should render in Ogre with no changes — the 98%+ pixel-accuracy figure above is verified against the exact fixture set.
