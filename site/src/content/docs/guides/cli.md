---
title: CLI Usage
description: Using the Ogre command-line tool.
---

The `ogre` CLI renders HTML files or inline HTML strings to images.

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--serve` | `false` | Start HTTP server mode |
| `--port` | `3000` | Server port (only with `--serve`) |
| `--render` | | Path to an HTML file to render |
| `--html` | | Inline HTML string to render |
| `--output` | | Output file path. Required for PNG/JPEG. If omitted for SVG, writes to stdout. |
| `--width` | `1200` | Canvas width in pixels |
| `--height` | `630` | Canvas height in pixels |
| `--format` | `svg` | Output format: `svg`, `png`, or `jpeg` |

## Rendering a file

```bash
ogre --render card.html --output card.svg
```

The HTML file should contain a single root `<div>` with inline styles or Tailwind classes.

## Rendering inline HTML

```bash
ogre --html '<div class="flex w-full h-full bg-blue-500 items-center justify-center"><div class="text-4xl font-bold text-white">Hello</div></div>' --output hello.svg
```

## Output formats

SVG output goes to stdout by default:

```bash
ogre --render card.html > card.svg
```

PNG and JPEG require `--output`:

```bash
ogre --render card.html --output card.png --format png
ogre --render card.html --output card.jpg --format jpeg
```

## Custom dimensions

```bash
ogre --render card.html --output card.png --format png --width 800 --height 400
```

## Server mode

```bash
ogre --serve --port 8080
```

See the [HTTP Server guide](/guides/server) for details.
