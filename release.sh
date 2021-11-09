#!/usr/bin/env bash
set -euo pipefail

TAG=v21.11.5
TITLE="CAP Client Release ${TAG}"
NOTES="Password change dialog"
ASSETS=""

function build_linux {
    docker build .fyne-cross/linux/ -t capclient-linux
    fyne-cross linux -name capclient -image capclient-linux:latest -env CGO_CFLAGS="-I/usr/include/ykpers-1/"
    ASSET="fyne-cross/capclient.${TAG}_Linux.tar.xz"
    ln -f fyne-cross/dist/linux-amd64/capclient.tar.xz ${ASSET}
    ASSETS="${ASSETS} ${ASSET}"
}

function build_mac {
    fyne-cross darwin -name capclient --app-id "com.aeolustec.capclient"
    ASSET="fyne-cross/Gocap.${TAG}_Mac.zip"
    zip -r -j ${ASSET} fyne-cross/dist/darwin-amd64/capclient.app
    ASSETS="${ASSETS} ${ASSET}"
}

function build_windows {
    docker build .fyne-cross/windows/ -t capclient-windows
    fyne-cross windows -name capclient.exe -image capclient-windows:latest -env CGO_CFLAGS="-I/usr/include/ykpers-1/ -I/usr/share/mingw-w64/include/" -env CGO_LDFLAGS=-L/usr/x86_64-w64-mingw32/lib
    ASSET="fyne-cross/capclient.${TAG}_Windows.zip"
    ln -f fyne-cross/dist/windows-amd64/capclient.exe.zip ${ASSET}
    ASSETS="${ASSETS} ${ASSET}"
}

git clean -fdx
go generate ./...

build_windows
build_linux
build_mac

gh release create ${TAG} \
    --draft \
    -n "${NOTES}" -t "${TITLE}" \
    "${ASSETS}"
