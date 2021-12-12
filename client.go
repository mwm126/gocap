package main

import (
	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/config"
	"aeolustec.com/capclient/fe261"
	"aeolustec.com/capclient/joule"
	"aeolustec.com/capclient/login"
	"aeolustec.com/capclient/watt"
	"crypto/rand"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
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

// Client represents the Main window of CAP client
type Client struct {
	Tabs     *container.AppTabs
	window   fyne.Window
	app      fyne.App
	LoginTab login.CapTab
}

func NewClient(
	a fyne.App,
	w fyne.Window,
	cfg config.Config,
	conn_man cap.ConnectionManager,
) Client {

	about_tab := container.NewTabItemWithIcon(
		"About",
		theme.HomeIcon(),
		widget.NewLabel(
			"The CAP client is used for connecting to Joule, Watt, and other systems using the CAP protocol.",
		),
	)
	service := login.Service{ // TODO: placeholder for real ServiceList service
		Name: "ServiceList",
		Networks: map[string]login.Network{
			"external": {
				ClientExternalAddress: "0.0.0.0",
				CapServerAddress:      "204.154.139.11",
			},
		},
	}

	connctd := container.NewVBox(widget.NewLabel("Connected!"))
	tabs := container.NewAppTabs(about_tab)
	login_tab := login.NewLoginTab("Login", "NETL SuperComputer", service, conn_man,
		func(conn cap.Connection) {
			services, err := login.FindServices()
			if err != nil {
				return
			}
			for _, service := range services {
				if service.Name == "joule" {
					joule := joule.NewJouleConnected(a, service, conn_man)
					tabs.Append(joule.CapTab.Tab)
				}
				if service.Name == "watt" {
					watt := watt.NewWattConnected(a, service, conn_man)
					tabs.Append(watt.CapTab.Tab)
				}
				if service.Name == "fe261" {
					fe261 := fe261.NewFe261Connected(a, service, conn_man)
					tabs.Append(fe261.CapTab.Tab)
				}
			}
			w.SetContent(tabs)

			// fe261_tab.Connect(conn)
		}, connctd)

	tabs.Append(login_tab.Tab)

	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)

	return Client{tabs, w, a, login_tab}
}

func (client *Client) Run() {
	client.window.ShowAndRun()
}
