#!/bin/bash

PLATFORMS=(
  'linux/amd64'
  'linux/arm64'
  'linux/arm7'
  'linux/arm5'
  'linux/386'
  'windows/amd64'
  'windows/386'
  'darwin/amd64'
  'darwin/arm64'
  'freebsd/amd64'
  'freebsd/arm7'
  'linux/mips'
  'linux/mipsle'
  'linux/mips64'
  'linux/mips64le'
  'linux/riscv64'
)

type setopt >/dev/null 2>&1

set_goarm() {
  if [[ "$1" =~ arm([5,7]) ]]; then
    GOARCH="arm"
    GOARM="${BASH_REMATCH[1]}"
    GO_ARM="GOARM=${GOARM}"
  else
    GOARM=""
    GO_ARM=""
  fi
}
# use softfloat for mips builds
set_gomips() {
  if [[ "$1" =~ mips ]]; then
    if [[ "$1" =~ mips(64) ]]; then MIPS64="${BASH_REMATCH[1]}"; fi
    GO_MIPS="GOMIPS${MIPS64}=softfloat"
  else
    GO_MIPS=""
  fi
}

GOBIN="go"

$GOBIN version

LDFLAGS="'-s -w'"
FAILURES=""
ROOT=${PWD}
OUTPUT="${ROOT}/dist/torrs"

#### Build web
echo "Generate static"
go run ./cmd/genpages/gen_pages.go

#### Build
rm -fr ./dist/*
$GOBIN clean -i -r -cache --modcache
$GOBIN mod tidy

BUILD_FLAGS="-ldflags=${LDFLAGS}"

#####################################
### X86 build section
#####

for PLATFORM in "${PLATFORMS[@]}"; do
  GOOS=${PLATFORM%/*}
  GOARCH=${PLATFORM#*/}
  set_goarm "$GOARCH"
  set_gomips "$GOARCH"
  BIN_FILENAME="${OUTPUT}-${GOOS}-${GOARCH}${GOARM}"
  if [[ "${GOOS}" == "windows" ]]; then BIN_FILENAME="${BIN_FILENAME}.exe"; fi
  CMD="GOOS=${GOOS} GOARCH=${GOARCH} ${GO_ARM} ${GO_MIPS} CGO_ENABLED=0 ${GOBIN} build ${BUILD_FLAGS} -o ${BIN_FILENAME} ./cmd/main"
  echo "${CMD}"
  eval "$CMD" || FAILURES="${FAILURES} ${GOOS}/${GOARCH}${GOARM}"
done

#####################################
### Android build section
#####

declare -a COMPILERS=(
  "arm7:armv7a-linux-androideabi21-clang"
  "arm64:aarch64-linux-android21-clang"
  "386:i686-linux-android21-clang"
  "amd64:x86_64-linux-android21-clang"
)

export NDK_VERSION="27.2.12479018" # 27.2.12479018
export NDK_TOOLCHAIN="${PWD}/../../android-ndk-r27c/toolchains/llvm/prebuilt/linux-x86_64"

GOOS=android

for V in "${COMPILERS[@]}"; do
  GOARCH=${V%:*}
  COMPILER=${V#*:}
  export CC="$NDK_TOOLCHAIN/bin/$COMPILER"
  export CXX="$NDK_TOOLCHAIN/bin/$COMPILER++"
  set_goarm "$GOARCH"
  BIN_FILENAME="${OUTPUT}-${GOOS}-${GOARCH}${GOARM}"
  CMD="GOOS=${GOOS} GOARCH=${GOARCH} ${GO_ARM} CGO_ENABLED=1 ${GOBIN} build ${BUILD_FLAGS} -o ${BIN_FILENAME} ./cmd/main"
  echo "${CMD}"
  eval "${CMD}" || FAILURES="${FAILURES} ${GOOS}/${GOARCH}${GOARM}"
done

#####################################
### Windows build without GUI #######
#####################################
GOOS="windows"
LDFLAGS="'-s -w -H=windowsgui'"
BUILD_FLAGS="-ldflags=${LDFLAGS}"
BIN_FILENAME="${OUTPUT}-${GOOS}-amd64-nogui.exe"
CMD="GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 ${GOBIN} build ${BUILD_FLAGS} -o ${BIN_FILENAME} ./cmd/main"
echo "${CMD}"
eval "$CMD" || FAILURES="${FAILURES} windows/amd64 NOGUI"

cp -fr ./views ./dist

# eval errors
if [[ "${FAILURES}" != "" ]]; then
  echo ""
  echo "failed on: ${FAILURES}"
  exit 1
fi
