package main

import (
	"log"

	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/client"
	"aeolustec.com/capclient/config"
	"aeolustec.com/capclient/login"
	"fyne.io/fyne/v2/app"
)

func main() {
	cfg := config.GetConfig()
	cfg.Enable_joule = true
	cfg.Enable_watt = true
	config.WriteConfig(cfg)

	yk := new(cap.UsbYubikey)
	knk := cap.NewKnocker(yk, cfg.YubikeyTimeout)
	knk.StartMonitor()
	conn_man := cap.NewCapConnectionManager(cap.NewSshClient, knk)
	a := app.New()
	w := a.NewWindow("CAP Client")

	err := login.InitServices(nil)
	if err != nil {
		log.Println("Could not contact Service List server:", err)
		return
	}
	client := client.NewClient(a, w, cfg, conn_man)
	client.Run()
}
