package joule

import (
	"embed"
	"fmt"
	"os/exec"
)

// Must install TurboVNC under /Applications
//go:embed TurboVNC-Mac
var content embed.FS

func VncCmd(vncviewer_path, otp string, localPort int) *exec.Cmd {
	return exec.Command(
		vncviewer_path,
		fmt.Sprintf("127.0.0.1::%d", localPort),
		fmt.Sprintf("-Password=%s", otp),
	)
}

func RunVnc(otp, displayNumber string, localPort int) {
	doRunVnc("/TurboVNC-Mac/Contents/MacOS/TurboVNC Viewer", otp, displayNumber, localPort)
}
