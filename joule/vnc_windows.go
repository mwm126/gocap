package joule

import (
	"aeolustec.com/capclient/cap"
	"log"
	"os/exec"
	"strconv"
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
