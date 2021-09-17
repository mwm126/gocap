package main

import (
	"aeolustec.com/capclient/cap"
	"crypto/rand"
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

	client := cap.NewClient(knk)
	client.Run()
}
