---
title: Quick Start
description: Render your first image with Ogre in under a minute.
---

## CLI

Render an HTML file to SVG:

```bash
ogre --render template.html --output og.svg
```

Render inline HTML to PNG:

```bash
ogre --html '<div class="flex w-full h-full bg-blue-500 items-center justify-center"><div class="text-4xl font-bold text-white">Hello</div></div>' --output og.png --format png
```

Start the HTTP server:

```bash
ogre --serve --port 3000
```

## Go library

```go
package main

import (
    "os"

    "github.com/macawls/ogre"
)

func main() {
    result, err := ogre.Render(`
        <div class="flex flex-col w-full h-full bg-slate-900 p-16 justify-center">
            <div class="text-5xl font-bold text-white">My Blog Post</div>
            <div class="text-xl text-slate-400 mt-4">A subtitle here</div>
        </div>
    `, ogre.Options{Width: 1200, Height: 630})
    if err != nil {
        panic(err)
    }

    os.WriteFile("og.svg", result.Data, 0644)
}
```

## HTTP API

With the server running:

```bash
curl -X POST http://localhost:3000/render \
  -H "Content-Type: application/json" \
  -d '{
    "html": "<div class=\"flex w-full h-full bg-blue-500 items-center justify-center\"><div class=\"text-4xl font-bold text-white\">Hello</div></div>",
    "width": 1200,
    "height": 630,
    "format": "png"
  }' \
  -o og.png
```
