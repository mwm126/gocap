package watt

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func SpiceCmd(localPort uint) (*exec.Cmd, error) {
	log.Println("Looking for spicy in PATH")
	virt_viewer, err := exec.LookPath("spicy")
	if err != nil {
		return nil, errors.New("Could not find spicy in PATH. If using homebrew, install with brew install spicy-gtk.")
	}

	return exec.Command(
		virt_viewer,
		fmt.Sprintf("spice://127.0.0.1::%d", localPort),
	), nil
}

func RunSpice(localPort uint) {
	cmd, err := SpiceCmd(localPort)
	if err != nil {
		log.Println("Could not run Spice")
	}
	log.Println("\n\n\nRunSpice: ", cmd)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Println("spiceviewer output: ", string(output))
		log.Println("spiceviewer error: ", err)
	}

}
