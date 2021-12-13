//go:build ignore
// +build ignore

package main

// This program generates embeds/ directory and downloads external dependencies.
// It can be invoked by running:  go generate

import (
	"log"
	"os"

	"github.com/melbahja/got"
)

var PUTTY_URL = "https://the.earth.li/~sgtatham/putty/latest/w64/putty.exe"
var PUTTY_FILENAME = "embeds/putty.exe"

var YK_URL = "https://developers.yubico.com/yubikey-personalization/Releases/ykpers-1.20.0-win64.zip"
var YK_FILENAME = "embeds/ykpers-1.20.0-win64.zip"

func main() {
	download(PUTTY_FILENAME, PUTTY_URL)
	download(YK_FILENAME, YK_URL)
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
