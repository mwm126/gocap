package main

import (
	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/client"
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

	yk := new(cap.UsbYubikey)
	knk := cap.NewPortKnocker(yk, entropy)

	cfg := client.GetConfig()
	cfg.Enable_joule = true
	cfg.Enable_watt = true
	client.WriteConfig(cfg)

	conn_man := cap.NewCapConnectionManager(knk)
	a := app.New()
	w := a.NewWindow("CAP Client")

	client.InitServices(nil)
	client := client.NewClient(a, w, cfg, conn_man)
	client.Run()
}
