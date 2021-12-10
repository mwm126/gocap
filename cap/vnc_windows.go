package cap

import (
	"aeolustec.com/capclient/cap/connection"

	"log"
	"os/exec"
	"strconv"
)

func run_vnc(conn connection.Connection, otp, displayNumber string) {
	cmd := exec.Command(
		"x-terminal-emulator",
		"--",
		"vnc",
		"localhost",
		"-p",
		strconv.Itoa(connection.VNC_LOCAL_PORT),
	)
	err := cmd.Run()
	if err != nil {
		log.Println("gnome-terminal FAIL: ", err)
	}
}
