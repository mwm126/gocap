package main

import (
	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/config"
	"aeolustec.com/capclient/login"
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

	cfg := config.GetConfig()
	cfg.Enable_joule = true
	cfg.Enable_watt = true
	config.WriteConfig(cfg)

	conn_man := cap.NewCapConnectionManager(knk)
	a := app.New()
	w := a.NewWindow("CAP Client")

	login.InitServices(nil)
	client := NewClient(a, w, cfg, conn_man)
	client.Run()
}
