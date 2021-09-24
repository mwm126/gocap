package cap

import (
	_ "embed"
	"log"
)

//go:generate curl --insecure "https://the.earth.li/~sgtatham/putty/latest/w64/putty.exe" --output embeds/putty.exe

func run_ssh(conn_man *CapConnectionManager) {
	log.Println("TODO")
}
