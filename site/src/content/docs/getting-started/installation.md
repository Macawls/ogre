---
title: Installation
description: Install Ogre from a prebuilt binary, via `go install`, or as a Docker container.
---

## Download prebuilt binary

Every release ships static binaries at [github.com/Macawls/ogre/releases/latest](https://github.com/Macawls/ogre/releases/latest). Pick the archive for your OS and architecture:

| OS | File pattern |
|---|---|
| Linux | `ogre_<version>_linux_amd64.tar.gz` (or `_arm64`) |
| macOS (Apple Silicon) | `ogre_<version>_darwin_arm64.tar.gz` |
| macOS (Intel) | `ogre_<version>_darwin_amd64.tar.gz` |
| Windows | `ogre_<version>_windows_amd64.zip` (or `_arm64`) |

Each archive contains a single `ogre` binary (or `ogre.exe` on Windows). Extract it, put it on your `PATH`, and you're done.

### Linux and macOS

Extract the downloaded archive:

```bash
tar -xzf ogre_*.tar.gz
```

Move the binary somewhere on your `PATH`. `/usr/local/bin` is the usual choice for a system-wide install:

```bash
sudo mv ogre /usr/local/bin/
```

Check it works:

```bash
ogre --help
```

**macOS only:** the first run may show "cannot be opened because the developer cannot be verified." Clear the quarantine flag once:

```bash
xattr -d com.apple.quarantine /usr/local/bin/ogre
```

### Windows

1. Download the `.zip` from the releases page and unzip it — you'll get a single `ogre.exe`.
2. Move `ogre.exe` to a folder that's already on your `PATH` (e.g. `C:\Windows\System32`), or create a new folder like `C:\Users\<you>\bin` and add it to `PATH` via **System Properties → Environment Variables**.
3. Open a new PowerShell or Command Prompt window and run `ogre --help` to verify.

### Verifying the download (optional)

Every release ships a `checksums.txt`. Download it from the same release page and compare:

```bash
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
