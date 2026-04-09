#!/bin/bash
set -e
cd "$(dirname "$0")/../.."
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o website/docs/public/playground/ogre.wasm ./cmd/playground/
echo "Built ogre.wasm ($(du -h website/docs/public/playground/ogre.wasm | cut -f1))"
