#!/usr/bin/env bash
set -euo pipefail

TAG=v21.9.20
TITLE="CAP Client Release ${TAG}"
NOTES="changes xyz"

ASSET_MAC=fyne-cross/Gocap.${TAG}_Mac.zip
ASSET_LINUX=fyne-cross/gocap.${TAG}_Linux.tar.gz
ASSET_WINDOWS=fyne-cross/gocap.${TAG}_Windows.zip

git clean -fdx
go generate ./...
fyne-cross windows
fyne-cross linux
fyne-cross darwin  --app-id "com.aeolustec.capclient"

zip -r -j ${ASSET_MAC} fyne-cross/dist/darwin-amd64/gocap.app
ln fyne-cross/dist/linux-amd64/gocap.tar.gz ${ASSET_LINUX}
ln fyne-cross/dist/windows-amd64/gocap.exe.zip ${ASSET_WINDOWS}

gh release create ${TAG} --draft \
    ${ASSET_MAC} \
    ${ASSET_LINUX} \
    ${ASSET_WINDOWS} \
    -n "${NOTES}" -t "${TITLE}"
