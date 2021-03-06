package ssh

import (
	"aeolustec.com/capclient/cap"

	_ "embed"
	"log"
	"os"
	"os/exec"
	"strconv"
)

//go:embed embeds/putty.exe
var putty []byte

func run_ssh(conn *cap.Connection) {
	file, err := os.CreateTemp("", "putty.*.exe")
	defer os.Remove(file.Name())
	if err != nil {
		log.Fatal("could not open tempfile", err)
	}

	_, err = file.Write(putty)
	if err != nil {
		log.Fatal("could not write ", putty, " because: ", err)
	}
	file.Close()
	cmd := exec.Command(file.Name(),
		"127.0.0.1",
		"-l",
		conn.GetUsername(),
		"-pw",
		conn.GetPassword(),
		"-P",
		strconv.Itoa(cap.SSH_LOCAL_PORT),
	)
	err = cmd.Run()
	if err != nil {
		log.Println("problem running tempdir putty:  ", err)
	}
}
