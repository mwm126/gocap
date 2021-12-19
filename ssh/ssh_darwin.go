package ssh

import (
	_ "embed"
	"log"

	"aeolustec.com/capclient/cap"
)

//go:generate curl --insecure "https://the.earth.li/~sgtatham/putty/latest/w64/putty.exe" --output embeds/putty.exe

func run_ssh(conn *cap.Connection) {
	log.Println("TODO")
}
