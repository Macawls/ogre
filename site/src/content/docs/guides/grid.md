---
title: CSS Grid
description: Using CSS Grid in Ogre templates.
---

Ogre lays out `display: grid` containers with a two-axis track sizing algorithm. Use raw CSS or the Tailwind grid utilities — both drive the same engine.

## Basic usage

```html
<div style="display:grid; grid-template-columns: 1fr 2fr 1fr; gap: 16px; width: 800px">
  <div>A</div>
  <div>B</div>
  <div>C</div>
</div>
```

Or with Tailwind:

```html
<div class="grid grid-cols-3 gap-4 w-[800px]">
  <div>A</div>
  <div>B</div>
  <div>C</div>
</div>
```

## Track sizing

| Track | Meaning |
|-------|---------|
| `100px` | Fixed length. |
| `25%` | Percentage of the container's content-box width (or height for rows). |
| `1fr` | One share of remaining space after fixed and percentage tracks. Ratios honored (`1fr 2fr` → 1:2 split). |
| `auto` | Sized to content when the container is auto-sized; otherwise treated as a flex slot. |
| `repeat(N, ...)` | Expanded at parse time. `repeat(3, 1fr)` is equivalent to `1fr 1fr 1fr`. |
| `minmax(a, b)` | Currently treated as the max value. `minmax(0, 1fr)` behaves as `1fr`. |

## Placement

- `grid-column: 2 / 4` — occupy tracks 2..3 (line 4 exclusive).
- `grid-column: 1 / -1` — full-width span (last line of the explicit grid).
- `grid-column: span 3` — cover three tracks starting at the auto cursor.
- `grid-row: 2` — pin to row 2; column resolved by auto-placement.
- `grid-column: 3 / 1` — end < start is swapped per spec.

The auto-placement algorithm follows CSS Grid §8.5 in three buckets: both-axes-definite first, then major-axis-definite (row for row-flow), then minor-axis-definite, then fully-auto items. Row flow is the default; `grid-auto-flow: column` switches to column-major.

## Gap

`gap`, `row-gap`, and `column-gap` behave identically to flex. Tailwind `gap-N`, `gap-x-N`, `gap-y-N` all work.

## Phase A limitations

The engine ships CSS Grid Level 1 with these known deferrals:

- `minmax(a, b)` uses the max value; the min floor is not honored.
- `min-content`, `max-content`, and `fit-content()` degrade to `auto`.
- `repeat(auto-fill, ...)` and `repeat(auto-fit, ...)` are not expanded.
- `grid-template-areas` and named grid lines are not parsed.
- `subgrid` is not supported.
- `justify-self` and `align-self` on cell items degrade to stretch.
- `align-items: baseline` degrades to `start`.
- Dense auto-placement (`grid-auto-flow: dense`) is not smarter than row-flow yet.
- Auto-placement will bail after 10,000 implicit tracks to prevent runaway growth from pathological input.

If you need one of the deferred features, open an issue and file a use case.

## Debugging

Set `OGRE_DEBUG=1` to have grid render with per-track debug outlines (planned; not yet implemented). For now, the compare tool in `test/compare/` is the fastest way to see track positions.
