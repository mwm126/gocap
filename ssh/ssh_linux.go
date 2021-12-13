package ssh

import (
	"aeolustec.com/capclient/cap"

	"log"
	"os/exec"
	"strconv"
)

//go:generate curl --insecure "https://the.earth.li/~sgtatham/putty/latest/w64/putty.exe" --output embeds/putty.exe

func run_ssh(conn cap.Connection) {
	cmd := exec.Command(
		"x-terminal-emulator",
		"--",
		"ssh",
		"localhost",
		"-p",
		strconv.Itoa(cap.SSH_LOCAL_PORT),
	)
	err := cmd.Run()
	if err != nil {
		log.Println("gnome-terminal FAIL: ", err)
	}
}
