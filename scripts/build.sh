
#!/usr/bin/env bash
set -e

APP_NAME=pgtool
VERSION=${VERSION:-v0.1.0}
OUT_DIR=dist

rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

platforms=(
  "darwin/arm64"
  "darwin/amd64"
  "linux/amd64"
  "linux/arm64"
)

for p in "${platforms[@]}"; do
  GOOS=${p%/*}
  GOARCH=${p#*/}

  bin="${APP_NAME}-${VERSION}-${GOOS}-${GOARCH}"
  echo "Building $bin"

  GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUT_DIR/$bin" .
done
