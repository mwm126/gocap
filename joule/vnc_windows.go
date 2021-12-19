package joule

import (
	"log"
	"os/exec"
	"strconv"

	"aeolustec.com/capclient/cap"
)

func RunVnc(conn *cap.Connection, otp, displayNumber string) {
	cmd := exec.Command(
		"x-terminal-emulator",
		"--",
		"vnc",
		"localhost",
		"-p",
		strconv.Itoa(VNC_LOCAL_PORT),
	)
	err := cmd.Run()
	if err != nil {
		log.Println("gnome-terminal FAIL: ", err)
	}
}
