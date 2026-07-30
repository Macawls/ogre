---
title: Images
description: Embedding images in Ogre templates.
---

Ogre supports `<img>` tags and CSS `background-image: url(...)`. In SVG output, images are embedded as `<image>` elements with a data URI. In PNG and JPEG output, they are decoded and composited into the bitmap.

## Image sources

An `<img src>` can be a data URI or an HTTP(S) URL. Anything else falls back to a placeholder.

### Data URIs

Base64 PNG/JPEG or URL-encoded SVG. This is the most portable form — the render is fully self-contained with no network calls.

```html
<img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Ccircle cx='12' cy='12' r='10' fill='%23475569'/%3E%3C/svg%3E" style="width:24px;height:24px" />
```

### HTTPS URLs

Fetched at render time and cached in memory for the lifetime of the process. The fetch has a 5-second timeout — slow origins render as the fallback placeholder.

```html
<img src="https://api.dicebear.com/9.x/lorelei/svg?seed=Sarah" style="width:48px;height:48px;border-radius:24px" />
```

For repeatable renders (CI, air-gapped environments, batch jobs), inline images as data URIs instead of relying on network fetches.

## Sizing

Use `width`, `height`, and `aspect-ratio` like any other element. Images participate in flex layout normally.

```html
<img src="..." style="width:100%;aspect-ratio:16/9" />
```

## object-fit and object-position

`object-fit` controls how the image scales inside its box. `object-position` shifts it within that box.

| Value | Behavior |
|-------|----------|
| `fill` | Stretch to fill the box (default; ignores aspect ratio) |
| `contain` | Fit entirely inside the box; may leave empty space |
| `cover` | Fill the box while preserving aspect ratio; may crop |
| `scale-down` | `contain`, but never scale up beyond the image's natural size |
| `none` | No resizing |

```html
<img src="..." style="width:400px;height:200px;object-fit:cover;object-position:center" />
```

## Rounded and circular images

`border-radius` on an `<img>` clips the image via an SVG `<clipPath>`. Set the radius to half the width/height for a circle.

```html
<img src="..." style="width:48px;height:48px;border-radius:24px" />
```

## background-image

`background-image: url(...)` applies an image to a container's background. Unlike `<img>`, it composes with `background-size`, `background-position`, and `background-repeat` and doesn't add a layout box of its own.

```html
<div style="width:600px;height:300px;background-image:url('https://example.com/hero.jpg');background-size:cover;background-position:center">
  <div style="color:white;font-size:48px;padding:32px">Overlay text</div>
</div>
```

The same source types (data URI or HTTP(S) URL) are supported.

## Fallback behavior

When a `src` is missing, unreachable, or fails to decode, Ogre renders a gray placeholder with a 🖼 icon in place of the image. Broken images are visible rather than silent.

## Supported formats

- **PNG** — full support in SVG, PNG, and JPEG output.
- **JPEG** — full support in SVG, PNG, and JPEG output.
- **SVG** (as a data URI or fetched URL) — inlined for SVG output; rasterized internally for PNG and JPEG output.

`<picture>`, `srcset`, and `sizes` are not supported. Only `<img src>` is read.
