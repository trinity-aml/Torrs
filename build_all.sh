#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DIST_DIR="$ROOT_DIR/Dist"
APP_NAME="torrs"
PACKAGE="./cmd/main"

cd "$ROOT_DIR"

if ! command -v zip >/dev/null 2>&1; then
	echo "zip is required" >&2
	exit 1
fi

TAG="$(git describe --tags --abbrev=0 2>/dev/null || true)"
if [[ -z "$TAG" ]]; then
	TAG="untagged"
fi
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || true)"
if [[ -z "$COMMIT" ]]; then
	COMMIT="unknown"
fi

ARCHIVE_NAME="${APP_NAME}-${TAG}-${COMMIT}.zip"

mkdir -p "$DIST_DIR"
rm -f "$DIST_DIR"/"${APP_NAME}"-* "$DIST_DIR"/"${APP_NAME}"-*.zip

# Keep the matrix aligned with runtime dependencies.
# blevesearch/mmap-go supports these GOOS values for normal binaries.
# Linux is intentionally limited to the common server targets.
DEFAULT_PLATFORMS="$(go tool dist list | awk '
	/^(darwin|freebsd|windows)\// { print; next }
	/^linux\/(386|amd64|arm|arm64)$/ { print }
')"
PLATFORMS="${PLATFORMS:-$DEFAULT_PLATFORMS}"

echo "Build tag: $TAG"
echo "Build commit: $COMMIT"
echo "Output: $DIST_DIR"

builds=()

while IFS= read -r platform; do
	[[ -z "$platform" ]] && continue
	goos="${platform%/*}"
	goarch="${platform#*/}"
	output="${DIST_DIR}/${APP_NAME}-${goos}-${goarch}"
	if [[ "$goos" == "windows" ]]; then
		output="${output}.exe"
	fi

	echo "Building ${goos}/${goarch} -> ${output#$ROOT_DIR/}"
	CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" go build -trimpath -ldflags='-s -w' -o "$output" "$PACKAGE"
	builds+=("$(basename "$output")")
done <<< "$PLATFORMS"

if [[ "${#builds[@]}" -eq 0 ]]; then
	echo "No binaries were built" >&2
	exit 1
fi

(
	cd "$DIST_DIR"
	zip -9 -q "$ARCHIVE_NAME" "${builds[@]}"
)

echo "Archive: Dist/$ARCHIVE_NAME"
