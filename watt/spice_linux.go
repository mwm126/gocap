package watt

import (
	"fmt"
	"os/exec"
)

func get_spice_cmd(port int) (string, error) {
	remote_viewer, err := exec.LookPath("remote-viewer")

	if err != nil {
		fmt.Println(" Could not find remote-viewer in PATH\n" + "You can download the installer here:" + "<a href='https://virt-manager.org/download/'>\n" + "https://virt-manager.org/download/</a>")
		return "", err
	}

	return fmt.Sprintf("env -u LD_LIBRARY_PATH %s spice://127.0.0.1:%d", remote_viewer, port), nil
}
