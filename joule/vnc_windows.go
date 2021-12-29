package joule

import (
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

//go:embed TurboVNC-2.2.7/app
var content embed.FS

func VncCmd(vncviewer_path, otp string, localPort int) *exec.Cmd {
	return exec.Command(
		vncviewer_path,
		fmt.Sprintf("127.0.0.1:%d", localPort),
		"/password",
		otp,
	)
}

func RunVnc(otp, displayNumber string, localPort int) {
	doRunVnc("/TurboVNC-2.2.7/app/vncviewer.exe", otp, displayNumber, localPort)
}
