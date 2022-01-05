package main

// For Windows development (not fyne-cross)
//go:generate curl -L --insecure "https://sourceforge.net/projects/turbovnc/files/2.2.7/TurboVNC-2.2.7-x64.exe/download" --output joule/embeds/TurboVNC-2.2.7-x64.exe
//go:generate innoextract -e joule/embeds/TurboVNC-2.2.7-x64.exe --output-dir joule/embeds/TurboVNC-2.2.7
//go:generate curl --insecure "https://the.earth.li/~sgtatham/putty/latest/w64/putty.exe" --output ssh/embeds/putty.exe
