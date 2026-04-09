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
2. Final stage uses Google's distroless image — just the binary, nothing else

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
docker run -p 8080:8080 ogre --serve --port 8080
```

## Docker Compose

```yaml
services:
  ogre:
    build: .
    ports:
      - "3000:3000"
    restart: unless-stopped
```

## Image size

The final image is minimal — distroless base with a single static binary. Typically under 10 MB.
