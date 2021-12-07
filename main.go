package main

import (
	"crypto/rand"
	"log"

	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/cap/connection"
)

func main() {
	var entropy [32]byte
	_, err := rand.Read(entropy[:])
	if err != nil {
		log.Fatal("Unable to get entropy to send CAP packet")
	}

	yk := new(connection.UsbYubikey)
	knk := connection.NewPortKnocker(yk, entropy)

	cfg := cap.GetConfig()
	cfg.Enable_joule = true
	cfg.Enable_watt = true
	cap.WriteConfig(cfg)

	conn_man := connection.NewCapConnectionManager(knk)
	client := cap.NewClient(cfg, conn_man)
	client.Run()
}
