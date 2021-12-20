package joule

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"aeolustec.com/capclient/cap"
)

func VncCmd(otp, displayNumber string) string {
	return fmt.Sprintf(
		"echo %s | env -u LD_LIBRARY_PATH vncviewer_HPCEE -highqual -autopass 127.0.0.1::%s &",
		otp,
		strconv.Itoa(VNC_LOCAL_PORT),
	)
}

func RunVnc(conn *cap.Connection, otp, displayNumber string) {
	cmd_string := VncCmd(otp, displayNumber)
	cmd := exec.Command(cmd_string)
	if err := cmd.Run(); err != nil {
		log.Println("vncviewer FAIL: ", err)
	}
}
