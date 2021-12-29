package main

// For fyne-cross on Linux to build the Windows distribution
//go:generate curl -L --insecure "https://sourceforge.net/projects/turbovnc/files/2.2.7/TurboVNC-2.2.7-x64.exe/download" --output TurboVNC-2.2.7-x64.exe
//go:generate innoextract -e TurboVNC-2.2.7-x64.exe --output-dir joule/TurboVNC-2.2.7
//go:generate curl --insecure "https://the.earth.li/~sgtatham/putty/latest/w64/putty.exe" --output ssh/embeds/putty.exe

// For fyne-cross on Linux to build the Linux distribution (or for development)
//go:generate cmake -S subprojects/turbovnc/ -B build-turbo -DCMAKE_INSTALL_PREFIX=joule/turbovnc/ -GNinja -DTVNC_BUILDSERVER=0
//go:generate cmake --build build-turbo --target install
