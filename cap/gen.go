//go:build ignore
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
	download(PUTTY_FILENAME, PUTTY_URL)
}

func download(path, url string) {
	if _, err := os.Stat(path); err == nil {
		log.Printf("%s already downloaded.", path)
		return
	}

	g := got.New()
	err := g.Download(url, path)
	if err != nil {
		log.Printf("Problem downloading %s: %s", path, err)
	} else {
		log.Printf("Successfully downloaded %s.", path)
	}
}
