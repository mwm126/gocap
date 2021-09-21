package main

import (
	"crypto/rand"
	"log"

	"aeolustec.com/capclient/cap"
)

func main() {
	var entropy [32]byte
	_, err := rand.Read(entropy[:])
	if err != nil {
		log.Fatal("Unable to get entropy to send CAP packet")
	}

	yk := new(cap.UsbYubikey)
	knk := cap.NewPortKnocker(yk, entropy)

	cfg := cap.GetConfig()
	cfg.Enable_joule = true
	cfg.Enable_watt = true
	cap.WriteConfig(cfg)

	client := cap.NewClient(knk)
	client.Run()
}
