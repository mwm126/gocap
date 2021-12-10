package cap

import (
	"aeolustec.com/capclient/cap/connection"
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

func run_vnc(conn connection.Connection, otp, displayNumber string) {
	cmd := exec.Command(fmt.Sprintf("echo %s | env -u LD_LIBRARY_PATH vncviewer_HPCEE -highqual -autopass 127.0.0.1::%s &",
		otp, strconv.Itoa(connection.VNC_LOCAL_PORT)))
	err := cmd.Run()
	if err != nil {
		log.Println("gnome-terminal FAIL: ", err)
	}
}
