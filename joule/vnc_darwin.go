package joule

import (
	"aeolustec.com/capclient/cap"
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

func RunVnc(conn cap.Connection, otp, displayNumber string) {
	cmd := exec.Command(
		fmt.Sprintf(
			"echo %s | env -u LD_LIBRARY_PATH vncviewer_HPCEE -highqual -autopass 127.0.0.1::%s &",
			otp,
			strconv.Itoa(VNC_LOCAL_PORT),
		),
	)
	err := cmd.Run()
	if err != nil {
		log.Println("gnome-terminal FAIL: ", err)
	}
}
