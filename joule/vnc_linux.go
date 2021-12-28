package joule

import (
	"fmt"
	"log"
	"os/exec"

	"aeolustec.com/capclient/cap"
)

func VncCmd(otp, displayNumber string, localPort int) string {
	return fmt.Sprintf(
		"echo %s | env -u LD_LIBRARY_PATH vncviewer_HPCEE -highqual -autopass 127.0.0.1::%d",
		otp,
		localPort,
	)
}

func RunVnc(conn *cap.Connection, otp, displayNumber string, localPort int) {
	cmd_string := VncCmd(otp, displayNumber, localPort)
	cmd := exec.Command("sh", "-c", cmd_string)
	log.Println("\n\n\nRunVnc: ", cmd)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Println("vncviewer output: ", string(output))
		log.Println("vncviewer error: ", err)
	}
}
