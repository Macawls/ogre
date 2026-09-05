#!/usr/bin/env bash
# Regenerate the OpenGraph card at site/public/og.png by rendering
# og-template.html with the Ogre CLI. Commit the PNG output.
#
# Requires: ogre on PATH (`go install github.com/macawls/ogre/v3/cmd/ogre@latest`).

set -euo pipefail

cd "$(dirname "$0")"

ogre \
  --render og-template.html \
  --output ../public/og.png \
  --format png \
  --width 1200 \
  --height 630

echo "Wrote $(realpath ../public/og.png)"
