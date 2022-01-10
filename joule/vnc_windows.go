package joule

import (
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
)

//go:embed embeds/TurboVNC-2.2.7/app
var vnc_content embed.FS

func VncCmd(vncviewer_path, otp string, localPort uint) *exec.Cmd {
	return exec.Command(
		vncviewer_path,
		fmt.Sprintf("127.0.0.1:%d", localPort),
		"/password",
		otp,
	)
}

func RunVnc(otp, displayNumber string, localPort int) {
	vnchome := extractVncToTempDir(otp, displayNumber, localPort)
	defer os.RemoveAll(vnchome)

	vnc_cmd := vnchome + "/embeds/TurboVNC-2.2.7/app/vncviewer.exe"
	err := os.Chmod(vnc_cmd, 0755)
	if err != nil {
		log.Fatal("could not make ", vnc_cmd, " executable because ", err)
	}

	cmd := VncCmd(vnc_cmd, otp, localPort)
	log.Println("\n\n\nRunVnc: ", cmd)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Println("vncviewer output: ", string(output))
		log.Println("vncviewer error: ", err)
	}

}
