package cap

import (
	"aeolustec.com/capclient/cap/connection"
	_ "embed"
	"log"
)

//go:generate curl --insecure "https://the.earth.li/~sgtatham/putty/latest/w64/putty.exe" --output embeds/putty.exe

func run_ssh(conn connection.Connection) {
	log.Println("TODO")
}
