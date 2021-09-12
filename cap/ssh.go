package cap

import (
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func newSsh(conn_man *CapConnectionManager, content embed.FS) *container.TabItem {
	ssh := widget.NewButton("New SSH Session", func() {
		os := runtime.GOOS
		switch os {
		case "windows":
			log.Println("Windows")
			run_putty(conn_man, content)
		case "darwin":
			log.Println("TODO")
		case "linux":
			log.Println("Linux")
			run_ssh()
		default:
			log.Printf("%s.\n", os)
		}

	})
	label := widget.NewLabel(fmt.Sprintf("or run in a terminal:  ssh localhost -p %d", SSH_LOCAL_PORT))
	box := container.NewVBox(widget.NewLabel("To create new Terminal (SSH) Session in gnome-terminal:"), ssh, label)
	return container.NewTabItem("SSH", box)
}

func run_ssh() {
	cmd := exec.Command("x-terminal-emulator", "--", "ssh", "localhost", "-p", strconv.Itoa(SSH_LOCAL_PORT))
	err := cmd.Run()
	if err != nil {
		log.Println("gnome-terminal FAIL: ", err)
	}
}

func run_putty(conn_man *CapConnectionManager, content embed.FS) {
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

	the_putty, err := content.ReadFile("embeds/putty.exe")
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
