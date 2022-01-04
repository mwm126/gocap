package joule

import (
	"embed"
	"fmt"
	"os/exec"
)

//go:embed embeds/turbovnc/bin/*
//go:embed embeds/turbovnc/share/*
var content embed.FS

func VncCmd(vncviewer_path, otp string, localPort int) *exec.Cmd {
	return exec.Command(
		vncviewer_path,
		fmt.Sprintf("127.0.0.1::%d", localPort),
		fmt.Sprintf("-Password=%s", otp),
	)
}

func RunVnc(otp, displayNumber string, localPort int) {
	doRunVnc("/turbovnc/bin/vncviewer", otp, displayNumber, localPort)
}
