package joule

import (
	"embed"
	"fmt"
	"os/exec"
)

//go:embed embeds/TurboVNC-2.2.7/app
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
