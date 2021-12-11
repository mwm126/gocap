package client

import (
	"aeolustec.com/capclient/cap"

	_ "embed"
	"log"
	"os"
	"os/exec"
	"strconv"
)

//go:generate curl --insecure "https://the.earth.li/~sgtatham/putty/latest/w64/putty.exe" --output embeds/putty.exe

//go:embed embeds/putty.exe
var putty []byte

func run_ssh(conn cap.Connection) {
	username := conn.connectionInfo.username
	password := conn.connectionInfo.password
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
		username,
		"-pw",
		password,
		"-P",
		strconv.Itoa(SSH_LOCAL_PORT),
	)
	err = cmd.Run()
	if err != nil {
		log.Println("problem running tempdir putty:  ", err)
	}
}