package client

import (
	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/config"
	"aeolustec.com/capclient/fe261"
	"aeolustec.com/capclient/joule"
	"aeolustec.com/capclient/login"
	"aeolustec.com/capclient/watt"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Client represents the Main window of CAP client
type Client struct {
	connectionManager *cap.ConnectionManager
	Tabs              *container.AppTabs
	window            fyne.Window
	app               fyne.App
	LoginTab          *login.LoginTab
}

func NewClient(
	a fyne.App,
	w fyne.Window,
	cfg config.Config,
	conn_man *cap.ConnectionManager,
) *Client {
	var client Client

	service := cap.Service{ // TODO: placeholder for real ServiceList service
		Name:    "ServiceList",
		CapPort: 62201,
		SshPort: 22,
		Networks: map[string]cap.Network{
			"external": {
				ClientExternalAddress: "0.0.0.0",
				CapServerAddress:      "204.154.139.11",
			},
		},
	}

	connctd := container.NewVBox(widget.NewLabel("Connected!"))

	uname, pword, _ := login.GetSavedLogin()
	login_tab := login.NewLoginTab(
		"Login",
		"NETL SuperComputer",
		service,
		conn_man,
		client.setupServices,
		connctd,
		uname,
		pword,
	)

	client = Client{conn_man, nil, w, a, login_tab}
	client.setupServices(nil, make([]cap.Service, 0))
	return &client
}

func (client *Client) setupServices(login_info *login.LoginInfo, services []cap.Service) {
	about_tab := container.NewTabItemWithIcon(
		"About",
		theme.HomeIcon(),
		widget.NewLabel(
			"The CAP client is used for connecting to Joule, Watt, and other systems using the CAP protocol.",
		),
	)

	tabs := container.NewAppTabs(about_tab)
	tabs.SetTabLocation(container.TabLocationLeading)

	if login_info == nil {
		tabs.Append(client.LoginTab.Tab)
	} else {
		for _, service := range services {
			if service.Name == "joule" {
				joule := joule.NewJouleConnected(
					client.app,
					client.window,
					service,
					client.connectionManager,
					*login_info,
				)
				tabs.Append(joule.CapTab.Tab)
			}
			if service.Name == "watt" {
				watt := watt.NewWattConnected(
					client.app,
					service,
					client.connectionManager,
					*login_info,
				)
				tabs.Append(watt.CapTab.Tab)
			}
			if service.Name == "fe261" {
				fe261 := fe261.NewFe261Connected(
					client.app,
					service,
					client.connectionManager,
					*login_info,
				)
				tabs.Append(fe261.CapTab.Tab)
			}
		}
	}
	client.Tabs = tabs
	client.window.SetContent(client.Tabs)
}

func (client *Client) Run() {
	client.window.ShowAndRun()
}
