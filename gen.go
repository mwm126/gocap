// +build ignore

package main

// This program generates embeds/ directory and downloads external dependencies.
// It can be invoked by running:  go generate

import (
	"github.com/melbahja/got"
	"log"
	"os"
)

var PUTTY_URL = "https://the.earth.li/~sgtatham/putty/latest/w64/putty.exe"
var PUTTY_FILENAME = "embeds/putty.exe"

func main() {
	err := os.Mkdir("embeds", 0755)
	if err != nil {
		log.Println(err)
	}

	if _, err := os.Stat(PUTTY_FILENAME); err == nil {
		log.Println("Putty already downloaded.")
		return
	}

	g := got.New()
	err = g.Download(PUTTY_URL, PUTTY_FILENAME)
	if err != nil {
		log.Println("Problem downloading Putty:", err)
	} else {
		log.Println("Successfully downloaded Putty.")
	}
}
