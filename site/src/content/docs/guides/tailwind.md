---
title: Tailwind CSS
description: Using Tailwind v3 and v4 utility classes in Ogre templates.
---

Ogre resolves Tailwind CSS utility classes directly at render time. No Tailwind CLI required. Both Tailwind v3 (default) and v4 are supported — pick one via `Options.TailwindVersion`.

## Basic usage

```html
<div class="flex flex-col w-full h-full bg-slate-900 p-16 justify-center">
  <div class="text-5xl font-bold text-white">Hello World</div>
  <div class="text-xl text-slate-400 mt-4">A subtitle here</div>
</div>
```

## Selecting v3 or v4

v3 is the default; passing an empty string preserves that. To opt into v4 semantics (renamed shadow/radius/blur scales, OKLCH palette, `bg-linear-to-*` gradients, `shrink`/`grow` utilities, `bg-(--var)` shortcut):

```go
result, _ := ogre.Render(html, ogre.Options{
    Width:           1200,
    Height:          630,
    TailwindVersion: ogre.TailwindV4,
})
```

Over HTTP:

```json
{
  "html": "<div class=\"bg-blue-500 shadow-sm rounded-sm\">...</div>",
  "tailwindVersion": "v4"
}
```

CLI:

```bash
ogre --tailwind-version v4 --html '...' --output og.png --format png
```

## Custom config

Extend the built-in palette, spacing scale, font families, or add arbitrary utility classes with a `TailwindConfig`. All fields are optional; nil means defaults.

### As a Go struct

```go
cfg := &ogre.TailwindConfig{
    Colors: map[string]map[int]string{
        "brand": {500: "#ff00aa", 600: "#ee0099"},
    },
    Spacing: map[string]string{
        "13":      "3.25rem",
        "sidebar": "220px",
    },
    Fonts: map[string]string{
        "display": `"Inter", sans-serif`,
    },
    Extend: map[string]map[string]string{
        "card": {"background-color": "#123456", "border-radius": "8px", "padding": "16px"},
    },
}

result, _ := ogre.Render(html, ogre.Options{
    Width:          1200,
    Height:         630,
    TailwindConfig: cfg,
})
```

Any class that would normally resolve against the built-in palette or spacing scale now consults `cfg` first. Missing keys fall back to defaults.

### As Tailwind v4 CSS

Ogre reads the same `@theme` block Tailwind v4 uses. Point `--tailwind-config` at a `.css` file:

```css
@theme {
  --color-brand-500: #ff00aa;
  --spacing-sidebar: 220px;
  --font-display: "Inter", sans-serif;
}
```

```bash
ogre --tailwind-config theme.css --html '...' --output og.png --format png
```

Supported namespaces: `--color-<name>-<shade>`, `--spacing-<key>`, `--font-<key>`, `--breakpoint-<key>`. Multiple `@theme` blocks merge. Unnamespaced variables are ignored.

### As JSON (HTTP or CLI)

The JSON schema mirrors the struct:

```json
{
  "html": "<div class=\"bg-brand-500 p-13\">...</div>",
  "tailwindConfig": {
    "colors": {"brand": {"500": "#ff00aa"}},
    "spacing": {"13": "3.25rem"}
  }
}
```

Or send the CSS body directly:

```json
{
  "html": "<div class=\"bg-brand-500\">...</div>",
  "tailwindConfigCSS": "@theme { --color-brand-500: #ff00aa; }"
}
```

CLI `.json` file works too — the extension picks the parser.

### Rules and limits

- **Precedence.** `Extend` is checked first (before arbitrary values and built-in classes), then colors/spacing/fonts overrides fall through to defaults on miss. Static utility classes (`font-bold`, `flex`, etc.) that don't take parameters still resolve through the built-in table and can only be overridden via `Extend`.
- **Colors.** Palette name is everything before the last `-<integer>` in the class. `--color-brand-primary-500` → `Colors["brand-primary"][500]`; `bg-brand-primary-500` reads it back. Shade must be an integer; keyword names like `--color-brand-DEFAULT` are not currently parsed.
- **Immutability.** A `TailwindConfig` must not be mutated after it has been passed to `Render`. The resolver caches lookups by content hash (`Key()`); mutating maps in place produces stale hits. Build the config fully before use, or call `Clone()`.
- **Trust.** `TailwindConfig.Extend` copies property values verbatim into the computed style. Only accept configs from trusted sources; do not expose the config surface to end users without validation. Ogre's HTML sanitiser does not apply to CSS emitted this way.
- **HTTP body limit.** The convenience `Renderer.Handler` caps request bodies at 10 MB. The full `server` package uses the same limit; set your own via middleware if you need more.

### Which name means what?

Some class names moved tiers in v4. Pick the version that matches your source.

| Class         | v3 output           | v4 output           |
| ------------- | ------------------- | ------------------- |
| `shadow-sm`   | subtle 1px shadow   | v3's `shadow` value |
| `shadow-xs`   | (not supported)     | v3's `shadow-sm`    |
| `shadow-2xs`  | (not supported)     | new tiny shadow     |
| `rounded-sm`  | 2px                 | 4px                 |
| `rounded-xs`  | (not supported)     | 2px                 |
| `blur-sm`     | 4px                 | 8px                 |
| `blur-xs`     | (not supported)     | 4px                 |
| `shrink`, `grow`               | (v3 uses `flex-shrink`/`flex-grow`) | supported |
| `bg-linear-to-*`               | (v3 uses `bg-gradient-to-*`)        | supported |
| `text-ellipsis`                | (v3 uses `overflow-ellipsis`)       | supported |

v4 also uses the OKLCH color palette, so `bg-blue-500` (and every other palette shade) resolves to a slightly different color than in v3.

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

In v4, the parentheses shortcut for CSS variables works too:

```html
<div class="bg-(--brand) w-(--card-w) p-(--card-p)">
  Uses var(--brand), var(--card-w), var(--card-p)
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
