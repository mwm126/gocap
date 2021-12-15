package main

import (
	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/config"
	"aeolustec.com/capclient/fe261"
	"aeolustec.com/capclient/joule"
	"aeolustec.com/capclient/login"
	"aeolustec.com/capclient/watt"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	cfg := config.GetConfig()
	cfg.Enable_joule = true
	cfg.Enable_watt = true
	config.WriteConfig(cfg)

	yk := new(cap.UsbYubikey)
	knk := cap.NewPortKnocker(yk, cfg.YubikeyTimeout)
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
	LoginTab login.LoginTab
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
		Name:    "ServiceList",
		CapPort: 62201,
		Networks: map[string]login.Network{
			"external": {
				ClientExternalAddress: "0.0.0.0",
				CapServerAddress:      "204.154.139.11",
			},
		},
	}

	connctd := container.NewVBox(widget.NewLabel("Connected!"))
	tabs := container.NewAppTabs(about_tab)

	uname, pword, _ := login.GetSavedLogin()
	login_tab := login.NewLoginTab("Login", "NETL SuperComputer", service, conn_man,
		func(login_info login.LoginInfo) {
			services, err := login.FindServices()
			if err != nil {
				return
			}
			for _, service := range services {
				if service.Name == "joule" {
					joule := joule.NewJouleConnected(a, service, conn_man, login_info)
					tabs.Append(joule.CapTab.Tab)
				}
				if service.Name == "watt" {
					watt := watt.NewWattConnected(a, service, conn_man, login_info)
					tabs.Append(watt.CapTab.Tab)
				}
				if service.Name == "fe261" {
					fe261 := fe261.NewFe261Connected(a, service, conn_man, login_info)
					tabs.Append(fe261.CapTab.Tab)
				}
			}
			w.SetContent(tabs)

			// fe261_tab.Connect(conn)
		}, connctd, uname, pword)

	tabs.Append(login_tab.Tab)

	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)

	return Client{tabs, w, a, login_tab}
}

func (client *Client) Run() {
	client.window.ShowAndRun()
}
