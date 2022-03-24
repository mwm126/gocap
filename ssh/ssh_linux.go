package ssh

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"aeolustec.com/capclient/cap"
)

func find_term() []string {
	path, err := exec.LookPath("gnome-terminal")
	if err == nil {
		return []string{path, "--"}
	}
	log.Println("gnome-terminal not found; trying xterm")

	path, err = exec.LookPath("xterm")
	if err == nil {
		return []string{path, "-e"}
	}
	log.Println("xterm not found.")

	return []string{}
}

func run_ssh(conn *cap.Connection) {
	args := find_term()
	if len(args) == 0 {
		log.Println("Could not start SSH session")
		return
	}
	sshpass, err := exec.LookPath("sshpass")
	if err != nil {
		log.Println("Could not find sshpass; re-enter password")
	} else {
		args = append(args, sshpass, "-e")
	}

	args = append(args,
		"ssh",
		"-o",
		"UserKnownHostsFile=/dev/null",
		"-o",
		"StrictHostKeyChecking=no",
		"-l",
		conn.GetUsername(),
		"-p",
		fmt.Sprint(cap.SSH_LOCAL_PORT),
		"127.0.0.1",
	)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("SSHPASS=%s", conn.GetPassword()),
	)

	err = cmd.Run()
	if err != nil {
		log.Println("Error: could not start SSH session in terminal: ", err)
	}
}
