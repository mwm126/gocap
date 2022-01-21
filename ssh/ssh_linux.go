package ssh

import (
	"fmt"
	"log"
	"os/exec"

	"aeolustec.com/capclient/cap"
)

func find_term() string {
	path, err := exec.LookPath("x-terminal-emulator")
	if err == nil {
		return path
	}
	log.Println("x-terminal-emulator not found; assuming gnome-terminal")

	path, err = exec.LookPath("gnome-terminal")
	if err == nil {
		return path
	}
	log.Println("Warning: gnome-terminal not found.  Falling back to xterm")

	return "xterm"
}

func run_ssh(conn *cap.Connection) {
	terminal := find_term()
	cmd := exec.Command(
		terminal,
		"-e",
		fmt.Sprintf("ssh %s@%s -p %d\n",
			conn.GetUsername(),
			"localhost",
			cap.SSH_LOCAL_PORT),
	)
	err := cmd.Run()
	if err != nil {
		log.Println("Error: could not start SSH session in terminal: ", err)
	}
}
