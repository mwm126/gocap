#!/usr/bin/env bash
set -euo pipefail

TAG=v21.12.30
TITLE="CAP Client Release ${TAG}"
NOTES="VNC functionality"
UNAME_S=$(uname -s)
ASSETS=""


function main {
    echo "Building on ${UNAME_S}...."

    build_linux
    build_mac
    build_windows

    if [[ $# -eq 0 || "$1" != "upload" ]]; then
        echo
        echo "Build successful. Run again as './release.sh upload' to upload assets to Github."
        exit
    fi

    gh release create ${TAG} \
       --draft \
       -n "${NOTES}" -t "${TITLE}" \
       "${ASSETS}"
}


function build_linux {
    echo
    echo "Building for Linux..."
    if [ "$UNAME_S" != "Linux" ]; then
        echo "Skipping Linux build...(must build Linux on Linux to have TurboVNC build dependencies)"
        return
    fi
    echo
    echo "Building for Linux..."
    echo
    env GOOS=linux go generate ./...
    docker build .fyne-cross/linux/ -t capclient-linux
    fyne-cross linux -name capclient -image capclient-linux:latest -env CGO_CFLAGS="-I/usr/include/ykpers-1/"
    ASSET_LINUX="fyne-cross/capclient.${TAG}_Linux.tar.xz"
    mv fyne-cross/dist/linux-amd64/capclient.tar.xz ${ASSET_LINUX}
    ASSETS="${ASSETS} ${ASSET_LINUX}"
}


function build_mac {
    echo
    echo "Building for Mac..."
    TURBO_HOME=/Applications/TurboVNC
    echo
    echo "Note:  This script will run sudo to DELETE your ${TURBO_HOME} directory, and then (re)install TurboVNC-2.2.7 to ${TURBO_HOME}."
    echo
    read -p "If you don't want this, Ctrl-C to cancel.  Otherwise, Enter to continue." -r
    echo

    if [ "$UNAME_S" != "Darwin" ]; then
        echo
        echo "Skipping Mac build...(must build Mac on Mac)"
        echo
        return
    fi
    env GOOS=darwin go generate ./...
    fyne-cross darwin -name capclient --app-id "com.aeolustec.capclient" -env CGO_CFLAGS="-I/usr/local/include/ykpers-1 -I/usr/local/include" -env CGO_LDFLAGS="/usr/local/lib/libykpers-1.a /usr/local/lib/libyubikey.a"
    ASSET_MAC="fyne-cross/Gocap.${TAG}_Mac.zip"
    zip -r -j ${ASSET_MAC} fyne-cross/dist/darwin-amd64/capclient.app
    ASSETS="${ASSETS} ${ASSET_MAC}"
}


function build_windows {
    echo
    echo "Building for Windows..."
    env GOOS=windows go generate ./...
    docker build .fyne-cross/windows/ -t capclient-windows
    fyne-cross windows -name capclient.exe -image capclient-windows:latest -env CGO_CFLAGS="-I/usr/include/ykpers-1/ -I/usr/share/mingw-w64/include/" -env CGO_LDFLAGS=-L/usr/x86_64-w64-mingw32/lib
    ASSET_WINDOWS="fyne-cross/capclient.${TAG}_Windows.zip"
    mv fyne-cross/dist/windows-amd64/capclient.exe.zip ${ASSET_WINDOWS}
    ASSETS="${ASSETS} ${ASSET_WINDOWS}"
}


main "$@"
