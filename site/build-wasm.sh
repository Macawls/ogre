#!/usr/bin/env bash
set -e
cd "$(dirname "$0")"
mkdir -p public/playground
DEST=public/playground/ogre.wasm

if [ -d ../cmd/wasm ] && command -v go >/dev/null 2>&1; then
    echo "building ogre.wasm from ../cmd/wasm..."
    (cd .. && GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o "site/$DEST" ./cmd/wasm/)
else
    echo "fetching ogre.wasm from GitHub releases..."
    TAG=$(curl -sSfL https://api.github.com/repos/Macawls/ogre/releases/latest \
        | grep -m1 '"tag_name"' \
        | sed -E 's/.*"tag_name": *"(v[0-9.]+)".*/\1/')
    VER=${TAG#v}
    URL="https://github.com/Macawls/ogre/releases/download/${TAG}/ogre_${VER}_wasm.tar.gz"
    echo "  $URL"
    curl -sSfL "$URL" | tar -xz -C public/playground ogre.wasm
fi
echo "ogre.wasm: $(du -h "$DEST" | cut -f1)"
