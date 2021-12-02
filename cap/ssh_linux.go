package cap

import (
	"aeolustec.com/capclient/cap/connection"

	"log"
	"os/exec"
	"strconv"
)

//go:generate curl --insecure "https://the.earth.li/~sgtatham/putty/latest/w64/putty.exe" --output embeds/putty.exe

func run_ssh(conn_man *connection.CapConnectionManager) {
	cmd := exec.Command(
		"x-terminal-emulator",
		"--",
		"ssh",
		"localhost",
		"-p",
		strconv.Itoa(connection.SSH_LOCAL_PORT),
	)
	err := cmd.Run()
	if err != nil {
		log.Println("gnome-terminal FAIL: ", err)
	}
}
