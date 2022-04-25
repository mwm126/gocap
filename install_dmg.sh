#!/bin/bash -lex

TURBO_HOME=/Applications/TurboVNC

# Remove previous installation
sudo rm -rf ${TURBO_HOME}
ls /Applications

TURBO_DMG=joule/embeds/TurboVNC-2.2.7.dmg
TURBO_APP=joule/embeds/TurboVNC-Mac

curl -L --insecure "https://sourceforge.net/projects/turbovnc/files/2.2.7/TurboVNC-2.2.7.dmg/download" --output ${TURBO_DMG}

listing=$(sudo hdiutil attach ${TURBO_DMG} | grep Volumes)

volume=$(echo "$listing" | cut -f 3)
package=TurboVNC.pkg

sudo installer -pkg "$volume"/"$package" -target LocalSystem

mkdir -p ${TURBO_APP}
rsync -av ${TURBO_HOME}/TurboVNC\ Viewer.app/ ${TURBO_APP}
