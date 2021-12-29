#!/bin/bash -lex

# Remove previous installation
sudo rm -rf /Applications/TurboVNC
ls /Applications

TURBO_DMG=TurboVNC-2.2.7.dmg
TURBO_APP=joule/TurboVNC-Mac

curl -L --insecure "https://sourceforge.net/projects/turbovnc/files/2.2.7/TurboVNC-2.2.7.dmg/download" --output ${TURBO_DMG}

listing=$(sudo hdiutil attach ${TURBO_DMG} | grep Volumes)

volume=$(echo "$listing" | cut -f 3)
package=$(ls -1 "$volume" | grep .pkg | head -1)

sudo installer -pkg "$volume"/"$package" -target LocalSystem

mkdir -p ${TURBO_APP}
rsync -av /Applications/TurboVNC/TurboVNC\ Viewer.app/ ${TURBO_APP}
