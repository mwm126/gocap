package watt

import (
	"fmt"
	"os/exec"
)

func get_spice_cmd(port int) str {
	remote_viewer, err := find_remote_viewer("/", port)
	if err == nil {
		return remote_viewer, nil
	}
	remote_viewer, err = find_remote_viewer(os.UserHomeDir(), port)
	if err == nil {
		return remote_viewer, nil
	}
	return "", err
}

func find_remote_viewer(base string, port int) (string, Error) {
	full_path = base + "/Applications/RemoteViewer.app/Contents/MacOS/RemoteViewer"
	_, err := os.Stat(full_path)
	if os.IsNotExist(err) {
		return "", Error("Could not find RemoteViewer.app under /Applications or ~/Applications folder. You can download the bundle from: https://www.spice-space.org/osx-client.html")
	}
	return fmt.Sprintf("%s spice://127.0.0.1:%d", remote_viewer, port), nil
}
