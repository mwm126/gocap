#!/usr/bin/env bash
set -euo pipefail

TAG=v21.11.5
TITLE="CAP Client Release ${TAG}"
NOTES="Password change dialog"
ASSETS=""

function build_linux {
    docker build .fyne-cross/linux/ -t capclient-linux
    fyne-cross linux -image capclient-linux:latest
    ASSET="fyne-cross/gocap.${TAG}_Linux.tar.xz"
    ln -f fyne-cross/dist/linux-amd64/gocap.tar.xz ${ASSET}
    ASSETS="${ASSETS} ${ASSET}"
}

function build_mac {
    fyne-cross darwin --app-id "com.aeolustec.capclient"
    ASSET="fyne-cross/Gocap.${TAG}_Mac.zip"
    zip -r -j ${ASSET} fyne-cross/dist/darwin-amd64/gocap.app
    ASSETS="${ASSETS} ${ASSET}"
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

build_mac
build_windows
build_linux

gh release create ${TAG} \
    --draft \
    -n "${NOTES}" -t "${TITLE}" \
    "${ASSETS}"
