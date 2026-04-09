---
title: CSS Properties
description: All CSS properties supported by Ogre.
---

Ogre supports a subset of CSS designed for image generation. All properties can be used as inline styles.

## Layout

| Property | Values |
|----------|--------|
| `display` | `flex`, `none`, `block`, `contents` |
| `position` | `static`, `relative`, `absolute` |
| `top`, `right`, `bottom`, `left` | Length values |
| `width`, `height` | Length, percentage, `auto` |
| `min-width`, `min-height` | Length values |
| `max-width`, `max-height` | Length values |
| `aspect-ratio` | Number (e.g. `16/9`) |
| `overflow` | `visible`, `hidden` |
| `box-sizing` | `content-box`, `border-box` |

## Flexbox

| Property | Values |
|----------|--------|
| `flex-direction` | `row`, `row-reverse`, `column`, `column-reverse` |
| `flex-wrap` | `nowrap`, `wrap`, `wrap-reverse` |
| `flex-grow` | Number |
| `flex-shrink` | Number |
| `flex-basis` | Length, percentage, `auto` |
| `align-items` | `flex-start`, `flex-end`, `center`, `stretch`, `baseline` |
| `align-self` | `auto`, `flex-start`, `flex-end`, `center`, `stretch`, `baseline` |
| `align-content` | `flex-start`, `flex-end`, `center`, `stretch`, `space-between`, `space-around` |
| `justify-content` | `flex-start`, `flex-end`, `center`, `space-between`, `space-around`, `space-evenly` |
| `gap`, `row-gap`, `column-gap` | Length values |

> [!NOTE]
> `<div>` defaults to `display: flex`, matching Satori behavior (not browser behavior).

## Box model

| Property | Values |
|----------|--------|
| `margin` | Length (all sides, shorthand) |
| `padding` | Length (all sides, shorthand) |
| `border-width` | Length (all sides, shorthand) |
| `border-style` | `solid`, `dashed`, `dotted` |
| `border-color` | Color value |
| `border-radius` | Length (all corners, shorthand) |

## Typography

| Property | Values |
|----------|--------|
| `font-family` | Font name |
| `font-size` | Length |
| `font-weight` | `100`-`900`, `normal`, `bold` |
| `font-style` | `normal`, `italic` |
| `color` | Color value |
| `line-height` | Number or length |
| `letter-spacing` | Length |
| `text-align` | `left`, `right`, `center`, `justify`, `start`, `end` |
| `text-transform` | `none`, `uppercase`, `lowercase`, `capitalize` |
| `text-decoration-line` | `none`, `underline`, `overline`, `line-through` |
| `text-decoration-color` | Color value |
| `text-decoration-style` | Style value |
| `text-shadow` | Shadow value |
| `white-space` | `normal`, `nowrap`, `pre`, `pre-wrap`, `pre-line` |
| `word-break` | `normal`, `break-all`, `break-word`, `keep-all` |
| `text-overflow` | `ellipsis` |
| `-webkit-line-clamp` | Number |

## Background

| Property | Values |
|----------|--------|
| `background-color` | Color value |
| `background-image` | `linear-gradient()`, `radial-gradient()`, `url()` |
| `background-size` | Length, percentage, `cover`, `contain` |
| `background-position` | Position value |
| `background-repeat` | Repeat value |

## Visual

| Property | Values |
|----------|--------|
| `opacity` | `0`-`1` |
| `box-shadow` | Shadow value |
| `transform` | Transform functions |
| `transform-origin` | Position value |
| `object-fit` | `fill`, `contain`, `cover`, `scale-down`, `none` |
| `object-position` | Position value |
| `filter` | `blur()`, `grayscale()`, `brightness()` |
| `clip-path` | Clip path value |

## Shorthands

These shorthand properties expand to their individual components:

`margin`, `padding`, `border`, `border-radius`, `flex`, `gap`, `background`, `font`, `text-decoration`, `overflow`, `border-top`, `border-right`, `border-bottom`, `border-left`, `border-width`, `border-style`, `border-color`

## Color values

Ogre accepts:
- Named colors: `white`, `black`, `red`, etc.
- Hex: `#ff0000`, `#f00`
- RGB: `rgb(255, 0, 0)`
- RGBA: `rgba(255, 0, 0, 0.5)`
- HSL: `hsl(0, 100%, 50%)`

## Length values

- Pixels: `16px`
- Percentages: `50%`
- Em: `1.5em`
- Rem: `1rem`
- Viewport: `100vw`, `100vh`
