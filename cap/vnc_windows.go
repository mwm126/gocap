package cap

import (
	"log"
	"os/exec"
	"strconv"
)

func RunVnc(conn connection.Connection, otp, displayNumber string) {
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
