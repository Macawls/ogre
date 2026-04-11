# Ogre Compare

Interactive comparison tool for Ogre, Satori, and Takumi renderers.

## Prerequisites

- Go
- [Bun](https://bun.sh) (for Satori and Takumi rendering)

## Setup

```bash
cd test/satori-reference
bun install
```

## Run

```bash
cd test/compare
go run .
```

Open http://localhost:4444 in your browser.

## How it works

- Ogre renders natively (in-process, Go)
- Satori and Takumi render via Bun subprocess (`render-one.ts`)
- Timing is measured internally by each renderer (not affected by subprocess startup)
- Pixel diff uses per-channel Euclidean distance with a threshold of 50
