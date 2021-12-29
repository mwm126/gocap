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

//go:embed joule/TurboVNC-2.2.7/app/
var content embed.FS

func VncCmd(vncviewer_path, otp string, localPort int) string {
	return fmt.Sprintf(
		"env -u LD_LIBRARY_PATH %s 127.0.0.1::%d -Password='%s'",
		vncviewer_path,
		localPort,
		otp,
	)
}

func RunVnc(otp, displayNumber string, localPort int) {
	vnchome, err := ioutil.TempDir("", "capclient")
	if err != nil {
		log.Fatal("could not open tempfile", err)
	}
	defer os.RemoveAll(vnchome)

	fs.WalkDir(content, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			dirname := vnchome + "/" + path
			os.Mkdir(dirname, 0755)
			return nil
		}
		src := path
		input, err := content.ReadFile(src)
		if err != nil {
			fmt.Println(err)
			return err
		}
		dest := vnchome + "/" + path
		err = ioutil.WriteFile(dest, input, 0644)
		if err != nil {
			fmt.Println("Error creating", dest)
			fmt.Println(err)
			return err
		}
		return nil
	})

	vnc_cmd := vnchome + "/vncviewer.exe"
	err = os.Chmod(vnc_cmd, 0755)
	if err != nil {
		log.Fatal("could not make ", vnc_cmd, " executable because ", err)
	}

	cmd_string := VncCmd(vnc_cmd, otp, localPort)
	cmd := exec.Command("sh", "-c", cmd_string)
	log.Println("\n\n\nRunVnc: ", cmd)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Println("vncviewer output: ", string(output))
		log.Println("vncviewer error: ", err)
	}
}
