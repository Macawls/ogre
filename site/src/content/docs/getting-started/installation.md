---
title: Installation
description: Install Ogre from a prebuilt binary, via `go install`, or as a Docker container.
---

## Download prebuilt binary

Every release ships static binaries for Linux, macOS, and Windows (amd64 and arm64) at [github.com/Macawls/ogre/releases/latest](https://github.com/Macawls/ogre/releases/latest). Archives are named `ogre_{version}_{os}_{arch}.{ext}` — `.tar.gz` for Linux and macOS, `.zip` for Windows.

### Linux

```bash
VERSION=3.1.0
ARCH=amd64  # or arm64
curl -L "https://github.com/Macawls/ogre/releases/download/v${VERSION}/ogre_${VERSION}_linux_${ARCH}.tar.gz" | tar -xz
sudo mv ogre /usr/local/bin/
```

### macOS

```bash
VERSION=3.1.0
ARCH=arm64  # or amd64 for Intel Macs
curl -L "https://github.com/Macawls/ogre/releases/download/v${VERSION}/ogre_${VERSION}_darwin_${ARCH}.tar.gz" | tar -xz
sudo mv ogre /usr/local/bin/
```

macOS may quarantine the downloaded binary. If you see "cannot be opened because the developer cannot be verified," run:

```bash
xattr -d com.apple.quarantine /usr/local/bin/ogre
```

### Windows

Download `ogre_<version>_windows_amd64.zip` (or `_arm64.zip`) from the [releases page](https://github.com/Macawls/ogre/releases/latest), extract, and place `ogre.exe` somewhere on your `PATH`.

### Verify checksum

Every release publishes a `checksums.txt` alongside the archives:

```bash
curl -LO "https://github.com/Macawls/ogre/releases/download/v${VERSION}/checksums.txt"
sha256sum -c checksums.txt --ignore-missing
```

## Go install

If you have Go 1.25 or later:

```bash
go install github.com/macawls/ogre/cmd/ogre@latest
```

This installs the `ogre` binary to `$GOPATH/bin` (usually `~/go/bin`).

## Go library

```bash
go get github.com/macawls/ogre@latest
```

Then import it:

```go
import "github.com/macawls/ogre"
```

## Docker

Pull and run the prebuilt image:

```bash
docker pull ghcr.io/macawls/ogre:latest
docker run -p 3000:3000 ghcr.io/macawls/ogre:latest
```

The image uses Google's distroless base and contains only the static binary. The container starts in server mode on port 3000. See the [Docker guide](/guides/docker/) for environment variables and configuration.

### Build from source

If you'd rather build the image yourself:

```bash
docker build -t ogre .
docker run -p 3000:3000 ogre
```

## Requirements

- **Prebuilt binaries and Docker:** none — the binary is static.
- **`go install` or library import:** Go 1.25 or later.
- **Building from source:** Go 1.25 or later.
