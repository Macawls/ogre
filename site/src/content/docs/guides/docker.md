---
title: Docker
description: Running Ogre in a container.
---

## Building the image

```bash
docker build -t ogre .
```

The Dockerfile uses a multi-stage build:
1. Build stage compiles the Go binary with CGO disabled
2. Final stage uses Google's distroless image

## Running

From a local build:

```bash
docker run -p 3000:3000 ogre
```

From GitHub Container Registry:

```bash
docker run -p 3000:3000 ghcr.io/macawls/ogre:latest
```

The container starts in server mode by default, listening on port 3000.

## Custom port

```bash
docker run -p 8080:8080 -e ADDR=:8080 ghcr.io/macawls/ogre:latest
```

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `ADDR` | `:3000` | Listen address |
| `CORS_ORIGIN` | `*` | Allowed CORS origin(s), comma-separated, supports wildcards |
| `CACHE_MB` | `64` | LRU cache size in MB |
| `RATE_LIMIT` | `0` | Requests per second per IP (0 = unlimited) |
| `TIMEOUT` | `10` | Render timeout in seconds |
| `MAX_ELEMENTS` | `1000` | Max HTML elements per render |

```bash
docker run -p 3000:3000 \
  -e CORS_ORIGIN=https://example.com \
  -e CACHE_MB=128 \
  -e RATE_LIMIT=10 \
  ghcr.io/macawls/ogre:latest
```

## Docker Compose

```yaml
services:
  ogre:
    image: ghcr.io/macawls/ogre:latest
    ports:
      - "3000:3000"
    environment:
      - CORS_ORIGIN=https://example.com
      - CACHE_MB=128
    restart: unless-stopped
```

## Image size

The final image is minimal — distroless base with a single static binary. Typically under 10 MB.
