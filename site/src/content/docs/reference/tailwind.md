---
title: Tailwind Classes Reference
description: Complete list of supported Tailwind v3 and v4 utility classes.
---

Ogre resolves Tailwind utility classes at render time. No Tailwind CLI required. The tables below list v3 (default) mappings; the [v4 differences](#v4-differences) section at the bottom notes what changes when you pass `TailwindVersion: TailwindV4`. Custom overrides via `TailwindConfig` (Go struct, Tailwind v4 `@theme` CSS, or JSON) are covered in the [Tailwind guide](/guides/tailwind/#custom-config).

## Layout

| Class | CSS |
|-------|-----|
| `flex` | `display: flex` |
| `flex-row` | `flex-direction: row` |
| `flex-col` | `flex-direction: column` |
| `flex-wrap` | `flex-wrap: wrap` |
| `flex-nowrap` | `flex-wrap: nowrap` |
| `flex-1` | `flex: 1 1 0%` |
| `flex-auto` | `flex: 1 1 auto` |
| `flex-initial` | `flex: 0 1 auto` |
| `flex-none` | `flex: none` |
| `flex-grow` | `flex-grow: 1` |
| `flex-grow-0` | `flex-grow: 0` |
| `flex-shrink` | `flex-shrink: 1` |
| `flex-shrink-0` | `flex-shrink: 0` |
| `hidden` | `display: none` |
| `block` | `display: block` |
| `relative` | `position: relative` |
| `absolute` | `position: absolute` |
| `grid` | `display: grid` |
| `inline-grid` | `display: grid` |

## Grid

See the [grid guide](/guides/grid/) for the layout algorithm and Phase A limitations.

| Class | CSS |
|-------|-----|
| `grid-cols-{1..12}` | `grid-template-columns: repeat(N, minmax(0, 1fr))` |
| `grid-cols-none` | `grid-template-columns: none` |
| `grid-rows-{1..6}` | `grid-template-rows: repeat(N, minmax(0, 1fr))` |
| `grid-rows-none` | `grid-template-rows: none` |
| `col-span-{1..12}` | `grid-column: span N / span N` |
| `col-span-full` | `grid-column: 1 / -1` |
| `col-start-{1..13}` | `grid-column-start: N` |
| `col-end-{1..13}` | `grid-column-end: N` |
| `row-span-{1..6}` | `grid-row: span N / span N` |
| `row-span-full` | `grid-row: 1 / -1` |
| `row-start-{1..7}` | `grid-row-start: N` |
| `row-end-{1..7}` | `grid-row-end: N` |
| `grid-flow-row` | `grid-auto-flow: row` |
| `grid-flow-col` | `grid-auto-flow: column` |
| `grid-flow-dense` | `grid-auto-flow: dense` |
| `grid-flow-row-dense` | `grid-auto-flow: row dense` |
| `grid-flow-col-dense` | `grid-auto-flow: column dense` |
| `auto-cols-{auto,min,max,fr}` | `grid-auto-columns: auto \| min-content \| max-content \| 1fr` |
| `auto-rows-{auto,min,max,fr}` | `grid-auto-rows: auto \| min-content \| max-content \| 1fr` |

## Alignment

| Class | CSS |
|-------|-----|
| `items-start` | `align-items: flex-start` |
| `items-end` | `align-items: flex-end` |
| `items-center` | `align-items: center` |
| `items-stretch` | `align-items: stretch` |
| `items-baseline` | `align-items: baseline` |
| `justify-start` | `justify-content: flex-start` |
| `justify-end` | `justify-content: flex-end` |
| `justify-center` | `justify-content: center` |
| `justify-between` | `justify-content: space-between` |
| `justify-around` | `justify-content: space-around` |
| `justify-evenly` | `justify-content: space-evenly` |
| `self-auto` | `align-self: auto` |
| `self-start` | `align-self: flex-start` |
| `self-end` | `align-self: flex-end` |
| `self-center` | `align-self: center` |
| `self-stretch` | `align-self: stretch` |

## Spacing

Pattern: `{property}-{size}` where size maps to the Tailwind spacing scale.

| Size | Value |
|------|-------|
| `0` | `0px` |
| `px` | `1px` |
| `0.5` | `2px` |
| `1` | `4px` |
| `1.5` | `6px` |
| `2` | `8px` |
| `3` | `12px` |
| `4` | `16px` |
| `5` | `20px` |
| `6` | `24px` |
| `8` | `32px` |
| `10` | `40px` |
| `12` | `48px` |
| `16` | `64px` |
| `20` | `80px` |
| `24` | `96px` |
| `32` | `128px` |
| `40` | `160px` |
| `48` | `192px` |
| `64` | `256px` |
| `80` | `320px` |
| `96` | `384px` |

Prefixes: `p`, `px`, `py`, `pt`, `pr`, `pb`, `pl`, `m`, `mx`, `my`, `mt`, `mr`, `mb`, `ml`, `gap`, `gap-x`, `gap-y`, `space-x`, `space-y`

## Sizing

| Class | CSS |
|-------|-----|
| `w-{n}` | `width: {n * 4}px` |
| `h-{n}` | `height: {n * 4}px` |
| `w-full` | `width: 100%` |
| `h-full` | `height: 100%` |
| `w-screen` | `width: 100vw` |
| `h-screen` | `height: 100vh` |
| `w-auto` | `width: auto` |
| `h-auto` | `height: auto` |
| `w-1/2` | `width: 50%` |
| `w-1/3` | `width: 33.333%` |
| `w-2/3` | `width: 66.667%` |
| `w-1/4` | `width: 25%` |
| `w-3/4` | `width: 75%` |

## Typography

| Class | CSS |
|-------|-----|
| `text-xs` | `font-size: 12px` |
| `text-sm` | `font-size: 14px` |
| `text-base` | `font-size: 16px` |
| `text-lg` | `font-size: 18px` |
| `text-xl` | `font-size: 20px` |
| `text-2xl` | `font-size: 24px` |
| `text-3xl` | `font-size: 30px` |
| `text-4xl` | `font-size: 36px` |
| `text-5xl` | `font-size: 48px` |
| `text-6xl` | `font-size: 60px` |
| `text-7xl` | `font-size: 72px` |
| `text-8xl` | `font-size: 96px` |
| `text-9xl` | `font-size: 128px` |
| `font-thin` | `font-weight: 100` |
| `font-light` | `font-weight: 300` |
| `font-normal` | `font-weight: 400` |
| `font-medium` | `font-weight: 500` |
| `font-semibold` | `font-weight: 600` |
| `font-bold` | `font-weight: 700` |
| `font-extrabold` | `font-weight: 800` |
| `font-black` | `font-weight: 900` |
| `text-left` | `text-align: left` |
| `text-center` | `text-align: center` |
| `text-right` | `text-align: right` |
| `italic` | `font-style: italic` |
| `uppercase` | `text-transform: uppercase` |
| `lowercase` | `text-transform: lowercase` |
| `capitalize` | `text-transform: capitalize` |
| `underline` | `text-decoration: underline` |
| `line-through` | `text-decoration: line-through` |
| `truncate` | `overflow: hidden; text-overflow: ellipsis; white-space: nowrap` |
| `line-clamp-{n}` | `-webkit-line-clamp: {n}` |

## Colors

Pattern: `{text|bg|border}-{palette}-{shade}`

Palettes: `slate`, `gray`, `zinc`, `neutral`, `stone`, `red`, `orange`, `amber`, `yellow`, `lime`, `green`, `emerald`, `teal`, `cyan`, `sky`, `blue`, `indigo`, `violet`, `purple`, `fuchsia`, `pink`, `rose`

Shades: `50`, `100`, `200`, `300`, `400`, `500`, `600`, `700`, `800`, `900`, `950`

Special: `text-white`, `text-black`, `text-transparent`, `bg-white`, `bg-black`, `bg-transparent`

## Borders

| Class | CSS |
|-------|-----|
| `border` | `border-width: 1px` |
| `border-0` | `border-width: 0` |
| `border-2` | `border-width: 2px` |
| `border-4` | `border-width: 4px` |
| `border-8` | `border-width: 8px` |
| `rounded-none` | `border-radius: 0` |
| `rounded-sm` | `border-radius: 2px` |
| `rounded` | `border-radius: 4px` |
| `rounded-md` | `border-radius: 6px` |
| `rounded-lg` | `border-radius: 8px` |
| `rounded-xl` | `border-radius: 12px` |
| `rounded-2xl` | `border-radius: 16px` |
| `rounded-3xl` | `border-radius: 24px` |
| `rounded-full` | `border-radius: 9999px` |

## Effects

| Class | CSS |
|-------|-----|
| `shadow-sm` | Small shadow |
| `shadow` | Default shadow |
| `shadow-md` | Medium shadow |
| `shadow-lg` | Large shadow |
| `shadow-xl` | Extra large shadow |
| `shadow-2xl` | 2XL shadow |
| `shadow-none` | No shadow |
| `opacity-{n}` | `opacity: {n/100}` (0-100) |

## Arbitrary values

## Filters

| Class | CSS |
|-------|-----|
| `blur-none` | `filter: blur(0)` |
| `blur-sm` | `filter: blur(4px)` |
| `blur` | `filter: blur(8px)` |
| `blur-md` | `filter: blur(12px)` |
| `blur-lg` | `filter: blur(16px)` |
| `blur-xl` | `filter: blur(24px)` |
| `blur-2xl` | `filter: blur(40px)` |
| `blur-3xl` | `filter: blur(64px)` |
| `brightness-{n}` | `filter: brightness({n/100})` (0, 50, 75, 90, 95, 100, 105, 110, 125, 150, 200) |
| `grayscale` | `filter: grayscale(100%)` |
| `grayscale-0` | `filter: grayscale(0)` |

## Transforms

| Class | CSS |
|-------|-----|
| `rotate-{n}` | `transform: rotate({n}deg)` (0, 1, 2, 3, 6, 12, 45, 90, 180) |
| `scale-{n}` | `transform: scale({n/100})` (0, 50, 75, 90, 95, 100, 105, 110, 125, 150) |
| `scale-x-{n}` | `transform: scaleX({n/100})` |
| `scale-y-{n}` | `transform: scaleY({n/100})` |
| `translate-x-{n}` | `transform: translateX({spacing})` |
| `translate-y-{n}` | `transform: translateY({spacing})` |
| `skew-x-{n}` | `transform: skewX({n}deg)` |
| `skew-y-{n}` | `transform: skewY({n}deg)` |

## Arbitrary values

Use bracket notation for any property:

```html
<div class="text-[32px] bg-[#ff5500] w-[200px] p-[20px] rounded-[12px] gap-[8px] leading-[1.5] tracking-[0.05em] rotate-[15deg] blur-[2px] scale-[1.2]">
</div>
```

In v4, the parentheses shortcut resolves the token as a CSS variable:

```html
<div class="bg-(--brand) w-(--card-w) p-(--card-p)">
  <!-- expands to background-color: var(--brand); width: var(--card-w); padding: var(--card-p) -->
</div>
```

## v4 differences

Set `Options.TailwindVersion = ogre.TailwindV4` (or send `"tailwindVersion": "v4"` over HTTP, or `--tailwind-version v4` on the CLI) to opt in.

**Renamed scales.** These class names carry different values under v4:

| Class         | v3 value              | v4 value              |
| ------------- | --------------------- | --------------------- |
| `shadow-sm`   | 1px subtle shadow     | 1px+3px shadow (v3's `shadow`) |
| `shadow-xs`   | *(unsupported)*       | 1px subtle shadow     |
| `shadow-2xs`  | *(unsupported)*       | new tiny shadow       |
| `rounded-sm`  | `2px`                 | `4px`                 |
| `rounded-xs`  | *(unsupported)*       | `2px`                 |
| `blur-sm`     | `blur(4px)`           | `blur(8px)`           |
| `blur-xs`     | *(unsupported)*       | `blur(4px)`           |

**New v4-only class names.**

| Class                            | CSS                                        |
| -------------------------------- | ------------------------------------------ |
| `shrink`, `shrink-0`             | `flex-shrink: 1` / `0`                     |
| `grow`, `grow-0`                 | `flex-grow: 1` / `0`                       |
| `bg-linear-to-{t,tr,r,br,b,bl,l,tl}` | `linear-gradient(...)` direction        |
| `text-ellipsis`                  | `text-overflow: ellipsis`                  |
| `bg-(--var)`, `w-(--var)`, etc.  | `background-color: var(--var)`, etc.       |

**OKLCH palette.** All `bg-*`, `text-*`, `border-*`, `from-*`, `via-*`, `to-*` colors resolve to `oklch(...)` strings taken from Tailwind v4's `theme.css`. The color parser converts them to sRGB at rasterization time.

**Legacy v3 classes still recognized in v4 mode** (harmless fallthroughs): `flex-shrink`, `flex-grow`, `bg-gradient-to-*`, bare `shadow` / `rounded` / `blur`. These aren't defined by v4, but Ogre still resolves them so mixed-source markup renders instead of silently doing nothing.
