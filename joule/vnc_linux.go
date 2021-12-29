package joule

import (
	"embed"
	"fmt"
	"os/exec"
)

//go:embed turbovnc/bin/*
//go:embed turbovnc/share/*
var content embed.FS

func VncCmd(vncviewer_path, otp string, localPort int) *exec.Cmd {
	return exec.Command(
		vncviewer_path,
		fmt.Sprintf("127.0.0.1::%d", localPort),
		fmt.Sprintf("-Password='%s'", otp),
	)
}

func RunVnc(vncviewer_p, otp, displayNumber string, localPort int) {
	doRunVnc("/turbovnc/bin/vncviewer", otp, displayNumber, localPort)
}
