#!/usr/bin/env python3

""" This script build releasable executables for the Go CAP client """

import os
import platform
import shutil
import subprocess
import sys

from pathlib import Path
from typing import Optional

DOCKER = "docker"
TAG = "v22.4.21"
TITLE = "CAP Client Release ${TAG}"
NOTES = "Watt web interface"
UNAME_S = str(platform.uname())


def main(args: list[str]) -> None:
    """Main routine"""
    print(f"Building on {UNAME_S}....")

    assets = [build_windows(), build_linux(), build_mac()]
    assets = [asset for asset in assets if asset is not None]

    if args and args[0] != "upload":
        print()
        print(
            "Build successful. Run again as './release.sh upload' to upload assets to Github."
        )
        sys.exit(0)

        subprocess.run(
            [
                "gh",
                "release",
                "create",
                TAG,
                "--draft",
                "-n",
                NOTES,
                "-t",
                TITLE,
                assets,
            ],
            check=True,
        )


def build_linux() -> Optional[Path]:
    """Build on Linux (if running on Linux)"""
    if UNAME_S != "Linux":
        print(
            "Skipping Linux build...(must build Linux on Linux to have TurboVNC build dependencies)"
        )
        return None
    print()
    print("Building for Linux...")
    print()
    env = os.environ
    env["GOOS"] = "linux"
    subprocess.run(["go", "generate", "./..."], env=env, check=True)
    subprocess.run(
        [DOCKER, "build", ".fyne-cross/linux/", "-t", "capclient-linux"], check=True
    )
    subprocess.run(
        [
            "fyne-cross",
            "linux",
            "-name",
            "capclient",
            "-image",
            "capclient-linux:latest",
            "-env",
            "CGO_CFLAGS=-I/usr/include/ykpers-1/",
        ],
        check=True,
    )
    asset_linux = Path("fyne-cross/capclient.${TAG}_Linux.tar.xz")
    subprocess.run(
        ["mv", "fyne-cross/dist/linux-amd64/capclient.tar.xz", asset_linux], check=True
    )
    return asset_linux


def build_mac() -> Optional[Path]:
    """Build on Mac (if running on Mac)"""
    if UNAME_S != "Darwin":
        print("\nSkipping Mac build...(must build Mac on Mac)\n")
        return None
    print()
    print("\nBuilding for Mac...\n")
    print()
    turbo_home = "/Applications/TurboVNC"
    print(
        f"Note:  This script will run sudo to DELETE your {turbo_home} directory,"
        " and then (re)install TurboVNC-2.2.7 to {turbo_home}."
    )
    input(
        "\nIf you don't want this, Ctrl-C to cancel.  Otherwise, Enter to continue.\n"
    )

    env = os.environ
    env["GOOS"] = "darwin"
    subprocess.run(["go", "generate", "./..."], env=env, check=True)
    subprocess.run(
        [
            "fyne-cross",
            "darwin",
            "-name",
            "capclient",
            "--app-id",
            "com.aeolustec.capclient",
            "-env",
            "CGO_CFLAGS=-I/usr/local/include/ykpers-1",
            "-I/usr/local/include",
            "-env",
            "CGO_LDFLAGS=/usr/local/lib/libykpers-1.a",
            "/usr/local/lib/libyubikey.a",
        ],
        check=True,
    )
    asset_mac = Path("fyne-cross/Gocap.${TAG}_Mac.zip")
    subprocess.run(
        ["zip", "-r", "-j", asset_mac, "fyne-cross/dist/darwin-amd64/capclient.app"],
        check=True,
    )
    return asset_mac


def build_windows() -> Optional[Path]:
    """Build on Windows"""
    print()
    print("Building for Windows...")
    print()
    env = os.environ
    env["GOOS"] = "windows"
    subprocess.run(["go", "generate", "./..."], env=env, check=True)
    subprocess.run(
        [DOCKER, "build", ".fyne-cross/windows/", "-t", "capclient-windows"],
        check=True,
    )
    subprocess.run(
        [
            "fyne-cross",
            "windows",
            "-name",
            "capclient.exe",
            "-image",
            "capclient-windows:latest",
            "-env",
            "CGO_CFLAGS=-I/usr/include/ykpers-1/ -I/usr/share/mingw-w64/include/",
            "-env",
            "CGO_LDFLAGS=-L/usr/x86_64-w64-mingw32/lib",
        ],
        check=True,
    )
    asset_windows = Path("fyne-cross/capclient.${TAG}_Windows.zip")
    shutil.move("fyne-cross/dist/windows-amd64/capclient.exe.zip", asset_windows)
    return asset_windows


main(sys.argv)
