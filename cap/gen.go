//go:build ignore
// +build ignore

package main

// This program generates embeds/ directory and downloads external dependencies.
// It can be invoked by running:  go generate

import (
	"github.com/melbahja/got"
	"log"
	"os"
	"runtime"
)

var PUTTY_URL = "https://the.earth.li/~sgtatham/putty/latest/w64/putty.exe"
var PUTTY_FILENAME = "embeds/putty.exe"

func main() {
	err := os.Mkdir("embeds", 0755)

	if err != nil {
		log.Println(err)
	}

	os := runtime.GOOS
	switch os {
	case "windows":
		download(PUTTY_FILENAME, PUTTY_URL)
	case "darwin":
	case "linux":
	default:
		log.Printf("%s.\n", os)
	}
}

func download(path, url string) {
	if _, err := os.Stat(path); err == nil {
		log.Printf("{} already downloaded.", path)
		return
	}

	g := got.New()
	err := g.Download(url, path)
	if err != nil {
		log.Printf("Problem downloading {}: {}", path, err)
	} else {
		log.Printf("Successfully downloaded {}.", path)
	}
}
