package main

import (
	"aeolustec.com/capclient/cap"
	"crypto/rand"
	"embed"
)

func main() {
	var entropy [32]byte
	rand.Read(entropy[:])

	yk := cap.UsbYubikey{}
	knk := cap.NewPortKnocker(&yk, entropy)

	cfg := cap.GetConfig()
	cfg.Enable_joule = true
	cfg.Enable_watt = true
	cap.WriteConfig(cfg)

	client := cap.NewClient(&knk, content)
	client.Run()
}

//go:generate go run gen.go

//go:embed embeds/*
var content embed.FS
