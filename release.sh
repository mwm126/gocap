#!/usr/bin/env bash
set -euo pipefail

TAG=v21.12.29
TITLE="CAP Client Release ${TAG}"
NOTES="Password change dialog"

function build_linux {
    docker build .fyne-cross/linux/ -t capclient-linux
    fyne-cross linux -name capclient -image capclient-linux:latest -env CGO_CFLAGS="-I/usr/include/ykpers-1/"
    ASSET_LINUX="fyne-cross/capclient.${TAG}_Linux.tar.xz"
    ln -f fyne-cross/dist/linux-amd64/capclient.tar.xz ${ASSET_LINUX}
}

function build_mac {
    fyne-cross darwin -name capclient --app-id "com.aeolustec.capclient" -env CGO_CFLAGS="-I/usr/local/include/ykpers-1 -I/usr/local/include" -env CGO_LDFLAGS="/usr/local/lib/libykpers-1.a /usr/local/lib/libyubikey.a"
    ASSET_MAC="fyne-cross/Gocap.${TAG}_Mac.zip"
    zip -r -j ${ASSET_MAC} fyne-cross/dist/darwin-amd64/capclient.app
}

function build_windows {
    docker build .fyne-cross/windows/ -t capclient-windows
    fyne-cross windows -name capclient.exe -image capclient-windows:latest -env CGO_CFLAGS="-I/usr/include/ykpers-1/ -I/usr/share/mingw-w64/include/" -env CGO_LDFLAGS=-L/usr/x86_64-w64-mingw32/lib
    ASSET_WINDOWS="fyne-cross/capclient.${TAG}_Windows.zip"
    ln -f fyne-cross/dist/windows-amd64/capclient.exe.zip ${ASSET_WINDOWS}
}

git clean -fdx
go generate ./...

build_windows
build_linux
build_mac

gh release create ${TAG} \
    --draft \
    -n "${NOTES}" -t "${TITLE}" \
    ${ASSET_LINUX} ${ASSET_MAC} ${ASSET_WINDOWS}
