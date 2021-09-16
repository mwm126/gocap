package cap

import (
	"embed"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

var PUTTY_FILENAME = "embeds/putty.exe"

//go:generate go run gen.go

//go:embed embeds/*
var content embed.FS

func run_ssh(conn_man *CapConnectionManager) {
	conn := conn_man.GetConnection()
	if conn == nil {
		log.Println("Warning: no connection; unable to run Putty", conn_man)
		return
	}
	username := conn.connectionInfo.username
	password := conn.connectionInfo.password
	file, err := ioutil.TempFile("", "putty.*.exe")
	defer os.Remove(file.Name())
	if err != nil {
		log.Fatal("could not open tempfile", err)
	}

	the_putty, err := content.ReadFile(PUTTY_FILENAME)
	if err != nil {
		log.Fatal("Could not get embed", err)
	}
	file.Write(the_putty)
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
