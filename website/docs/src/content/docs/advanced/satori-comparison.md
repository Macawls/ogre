---
title: Satori Comparison
description: How Ogre compares to Vercel's Satori.
---

Ogre is designed as a Go-native alternative to [Satori](https://github.com/vercel/satori), Vercel's HTML/CSS to SVG converter written in TypeScript.

## Feature comparison

| Feature | Ogre | Satori |
|---------|------|--------|
| Language | Go | TypeScript |
| Output formats | SVG, PNG, JPEG | SVG only |
| Dependencies | stdlib + golang.org/x | yoga-wasm + others |
| Deployment | Single static binary | Node.js runtime |
| Tailwind support | Built-in (v3) | Via plugin |
| PNG output | Built-in | Requires resvg |
| Layout engine | Custom flexbox (W3C) | Yoga (via WASM) |
| HTTP server | Built-in with caching | BYO |
| Font embedding | SVG paths | SVG paths |
| Emoji | Twemoji CDN | Twemoji CDN |
| `<div>` default | `display: flex` | `display: flex` |
| Pixel accuracy | 95%+ vs Satori | Reference |

## When to use Ogre

- You are building a Go backend and want OG image generation without Node.js
- You need PNG/JPEG output without a separate rasterization step
- You want a single binary deployment
- You want built-in Tailwind support without a build step
- You want a production-ready HTTP server included

## When to use Satori

- You are already running a Node.js or Next.js stack
- You need the exact rendering behavior of the reference implementation
- You are using `@vercel/og` in Next.js

## Compatibility

Ogre accepts the same HTML/CSS subset as Satori. Templates written for Satori should work in Ogre with no changes. The 95%+ pixel accuracy metric is verified across 25 test fixtures that cover layout, typography, gradients, borders, shadows, and positioning.

The main behavioral difference is in edge cases around flex layout calculations. Satori uses Yoga (a C++ layout engine compiled to WASM), while Ogre uses a from-scratch Go implementation of the W3C flexbox spec. In the 5% of cases where output differs, it is typically sub-pixel positioning differences.
