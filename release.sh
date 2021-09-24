#!/usr/bin/env bash
set -euo pipefail

TAG=v21.9.21
TITLE="CAP Client Release ${TAG}"
NOTES="Build with fyne-cross"

ASSET_MAC=fyne-cross/Gocap.${TAG}_Mac.zip
ASSET_LINUX=fyne-cross/gocap.${TAG}_Linux.tar.gz
ASSET_WINDOWS=fyne-cross/gocap.${TAG}_Windows.zip

docker build .fyne-cross/linux/ -t capclient-linux
docker build .fyne-cross/windows/ -t capclient-windows

git clean -fdx
go generate ./...
fyne-cross linux -image capclient-linux:latest
fyne-cross windows -image capclient-windows:latest
fyne-cross darwin --app-id "com.aeolustec.capclient"

zip -r -j ${ASSET_MAC} fyne-cross/dist/darwin-amd64/gocap.app
ln -f fyne-cross/dist/linux-amd64/gocap.tar.gz ${ASSET_LINUX}
ln -f fyne-cross/dist/windows-amd64/gocap.exe.zip ${ASSET_WINDOWS}

gh release create ${TAG} \
    --draft \
    ${ASSET_MAC} \
    ${ASSET_LINUX} \
    ${ASSET_WINDOWS} \
    -n "${NOTES}" -t "${TITLE}"
