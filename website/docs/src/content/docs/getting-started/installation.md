---
title: Installation
description: How to install Ogre as a CLI or Go library.
---

## CLI

```bash
go install github.com/macawls/ogre/cmd/ogre@latest
```

This installs the `ogre` binary to your `$GOPATH/bin`.

## Go library

```bash
go get github.com/macawls/ogre@latest
```

Then import it:

```go
import "github.com/macawls/ogre"
```

## Docker

```bash
docker build -t ogre .
docker run -p 3000:3000 ogre
```

The Dockerfile uses a multi-stage build. The final image uses Google's distroless base and contains only the static binary. The container starts in server mode on port 3000.

## Requirements

- Go 1.25 or later (for CLI and library)
- Docker (optional, for containerized deployment)
