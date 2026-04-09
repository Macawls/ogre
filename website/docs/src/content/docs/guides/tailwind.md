---
title: Tailwind CSS
description: Using Tailwind v3 utility classes in Ogre templates.
---

Ogre resolves Tailwind CSS v3 utility classes directly at render time. No Tailwind CLI, no build step, no configuration.

## Basic usage

```html
<div class="flex flex-col w-full h-full bg-slate-900 p-16 justify-center">
  <div class="text-5xl font-bold text-white">Hello World</div>
  <div class="text-xl text-slate-400 mt-4">A subtitle here</div>
</div>
```

## Mixing with inline styles

Tailwind classes and inline styles can be combined. Inline styles take precedence.

```html
<div class="flex w-full h-full bg-slate-900" style="padding: 60px">
  <div class="text-5xl font-bold" style="color: #a78bfa">Custom color</div>
</div>
```

## Arbitrary values

Use bracket notation for values outside the default scale:

```html
<div class="text-[32px] bg-[#ff5500] w-[200px] p-[20px] rounded-[12px] gap-[8px] leading-[1.5] tracking-[0.05em]">
  Custom values
</div>
```

## Supported categories

### Layout
`flex`, `flex-row`, `flex-col`, `flex-wrap`, `flex-nowrap`, `flex-1`, `flex-auto`, `flex-initial`, `flex-none`, `flex-grow`, `flex-shrink`, `hidden`, `block`, `relative`, `absolute`

### Alignment
`items-start`, `items-end`, `items-center`, `items-stretch`, `items-baseline`, `justify-start`, `justify-end`, `justify-center`, `justify-between`, `justify-around`, `justify-evenly`, `self-auto`, `self-start`, `self-end`, `self-center`, `self-stretch`

### Spacing
`p-{n}`, `px-{n}`, `py-{n}`, `pt-{n}`, `pr-{n}`, `pb-{n}`, `pl-{n}`, `m-{n}`, `mx-{n}`, `my-{n}`, `mt-{n}`, `mr-{n}`, `mb-{n}`, `ml-{n}`, `gap-{n}`, `space-x-{n}`, `space-y-{n}`

Scale: `0` = 0px, `px` = 1px, `0.5` = 2px, `1` = 4px ... `96` = 384px

### Sizing
`w-{n}`, `h-{n}`, `size-{n}`, `w-full`, `h-full`, `w-screen`, `h-screen`, `w-auto`, `h-auto`, fractions (`w-1/2`, `w-1/3`, `w-2/3`, etc.), `min-w-*`, `max-w-*`, `min-h-*`, `max-h-*`

### Typography
`text-xs` through `text-9xl`, `font-thin` through `font-black`, `text-left`, `text-center`, `text-right`, `italic`, `uppercase`, `lowercase`, `capitalize`, `underline`, `line-through`, `leading-*`, `tracking-*`, `truncate`, `line-clamp-{1-6}`

### Colors
`text-{color}-{shade}`, `bg-{color}-{shade}`, `border-{color}-{shade}`

Available palettes: slate, gray, zinc, neutral, stone, red, orange, amber, yellow, lime, green, emerald, teal, cyan, sky, blue, indigo, violet, purple, fuchsia, pink, rose. Shades: 50-950.

### Borders
`border`, `border-{0,2,4,8}`, `border-{t,r,b,l}-{n}`, `border-solid`, `border-dashed`, `border-dotted`, `rounded-none` through `rounded-full`

### Effects
`shadow-sm` through `shadow-2xl`, `opacity-{0-100}`

### Filters
`blur-none`, `blur-sm`, `blur`, `blur-md`, `blur-lg`, `blur-xl`, `blur-2xl`, `blur-3xl`, `brightness-{0-200}`, `grayscale`, `grayscale-0`

### Transforms
`rotate-{0,1,2,3,6,12,45,90,180}`, `scale-{0,50,75,90,95,100,105,110,125,150}`, `scale-x-{n}`, `scale-y-{n}`, `translate-x-{n}`, `translate-y-{n}`, `skew-x-{n}`, `skew-y-{n}`

### Position
`z-{0-50}`, `top-{n}`, `right-{n}`, `bottom-{n}`, `left-{n}`, `inset-{n}`, `aspect-square`, `aspect-video`
