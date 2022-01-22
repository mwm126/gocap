package joule

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
)

//go:embed embeds/turbovnc/bin/*
//go:embed embeds/turbovnc/share/*
var vnc_content embed.FS

func VncCmd(vncviewer_path, otp string, localPort uint) *exec.Cmd {
	return exec.Command(
		vncviewer_path,
		fmt.Sprintf("127.0.0.1::%d", localPort),
		fmt.Sprintf("-Password=%s", otp),
	)
}

func RunVnc(otp, displayNumber string, localPort uint) {
	vnchome := extractVncToTempDir(otp, displayNumber, localPort)
	defer os.RemoveAll(vnchome)

	vnc_cmd := vnchome + "/embeds/turbovnc/bin/vncviewer"
	err := os.Chmod(vnc_cmd, 0755)
	if err != nil {
		log.Println("Could not run ", vnc_cmd, " because ", err)
		return
	}

	cmd := VncCmd(vnc_cmd, otp, localPort)
	log.Println("\n\n\nRunVnc: ", cmd)

	if os.Getenv("GOCAP_DEMO") == "" {
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Println("vncviewer output: ", string(output))
			log.Println("vncviewer error: ", err)
		}
	}

}
