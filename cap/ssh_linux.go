package cap

import (
	"log"
	"os/exec"
	"strconv"
)

//go:generate go run gen.go

func run_ssh(conn_man *CapConnectionManager) {
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
