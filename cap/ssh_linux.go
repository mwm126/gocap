package cap

import (
	"exec"
	"fmt"
	"log"
	"strconv"
)

func run_ssh() {
	cmd := exec.Command(
		"x-terminal-emulator",
		"--",
		"ssh",
		"localhost",
		"-p",
		strconv.Itoa(SSH_LOCAL_PORT),
	)
	err := cmd.Run()
	if err != nil {
		log.Println("gnome-terminal FAIL: ", err)
	}
}
