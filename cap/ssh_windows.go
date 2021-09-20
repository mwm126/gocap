package cap

import (
	_ "embed"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
)

//go:generate go run gen.go

//go:embed embeds/putty.exe
var putty []byte

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
