package main

// For fyne-cross on Linux to build the Linux distribution (or for development)
//go:generate cmake -S subprojects/turbovnc/ -B build-turbo -DCMAKE_INSTALL_PREFIX=joule/turbovnc/ -GNinja -DTVNC_BUILDSERVER=0
//go:generate cmake --build build-turbo --target install
