package ssh

import (
	"fmt"

	"aeolustec.com/capclient/cap"

	"log"
	"os/exec"
	"strconv"
)

func run_ssh(conn *cap.Connection) {
	cmd := exec.Command(
		"x-terminal-emulator",
		"--",
		"ssh",
		fmt.Sprintf("%s@%s", conn.GetUsername(), "localhost"),
		"-p",
		strconv.Itoa(cap.SSH_LOCAL_PORT),
	)
	err := cmd.Run()
	if err != nil {
		log.Println("gnome-terminal FAIL: ", err)
	}
}
