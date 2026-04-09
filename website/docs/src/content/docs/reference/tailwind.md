---
title: Tailwind Classes Reference
description: Complete list of supported Tailwind v3 utility classes.
---

Ogre resolves these Tailwind v3 classes at render time with no build step.

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
