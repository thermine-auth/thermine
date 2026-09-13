#!/usr/bin/env bash
#
# Build xermess for release.
#
#   scripts/build.sh            build the server into bin/
#   scripts/build.sh --web      build the admin panel as well
#
# The version and the commit are compiled in, so a running binary can say
# which one it is.

set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

binary="bin/xermess"
version="${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo dev)}"
commit="$(git rev-parse --short HEAD 2>/dev/null || echo unknown)"

echo "building $binary ($version, $commit)"

CGO_ENABLED=0 go build \
	-trimpath \
	-ldflags "-s -w -X main.version=$version -X main.commit=$commit" \
	-o "$binary" ./cmd/xermess

if [[ "${1:-}" == "--web" ]]; then
	echo "building the admin panel"
	cd web
	bun install --frozen-lockfile
	bun run build
fi

echo "done"
