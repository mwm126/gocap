#!/usr/bin/env bash
set -euo pipefail

TAG=v21.9.23
TITLE="CAP Client Release ${TAG}"
NOTES="Link in libyubikey and libykpers"
ASSETS=""

function build_linux {
    docker build .fyne-cross/linux/ -t capclient-linux
    fyne-cross linux -image capclient-linux:latest
    ASSET="fyne-cross/gocap.${TAG}_Linux.tar.gz"
    ln -f fyne-cross/dist/linux-amd64/gocap.tar.gz ${ASSET}
    ASSETS="${ASSETS} ${ASSET}"
}

function build_mac {
    ASSET="fyne-cross/Gocap.${TAG}_Mac.zip"
    ASSETS="${ASSETS} ${ASSET}"
    # fyne-cross darwin --app-id "com.aeolustec.capclient"
    zip -r -j ${ASSET} fyne-cross/dist/darwin-amd64/gocap.app
}

function build_windows {
    docker build .fyne-cross/windows/ -t capclient-windows
    fyne-cross windows -image capclient-windows:latest
    ASSET="fyne-cross/gocap.${TAG}_Windows.zip"
    ln -f fyne-cross/dist/windows-amd64/gocap.exe.zip ${ASSET}
    ASSETS="${ASSETS} ${ASSET}"
}

git clean -fdx
go generate ./...

build_linux
# build_mac
build_windows

gh release create ${TAG} \
    --draft \
    -n "${NOTES}" -t "${TITLE}" \
    "${ASSETS}"
