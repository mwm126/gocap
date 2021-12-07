package main

import (
	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/cap/connection"
	"crypto/rand"
	"fyne.io/fyne/v2/app"
	"log"
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
	a := app.New()
	w := a.NewWindow("CAP Client")
	client := cap.NewClient(a, w, cfg, conn_man)
	client.Run()
}
