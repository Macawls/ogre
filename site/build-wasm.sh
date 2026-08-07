#!/usr/bin/env bash
set -e
cd "$(dirname "$0")/.."
mkdir -p site/public/playground
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o site/public/playground/ogre.wasm ./cmd/wasm/
echo "Built ogre.wasm ($(du -h site/public/playground/ogre.wasm | cut -f1))"
