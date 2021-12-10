package client

import (
	"aeolustec.com/capclient/cap"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Client represents the Main window of CAP client
type Client struct {
	Tabs     *container.AppTabs
	window   fyne.Window
	app      fyne.App
	LoginTab CapTab
}

func NewClient(a fyne.App, w fyne.Window, cfg config, conn_man cap.ConnectionManager) Client {

	about_tab := container.NewTabItemWithIcon(
		"About",
		theme.HomeIcon(),
		widget.NewLabel(
			"The CAP client is used for connecting to Joule, Watt, and other systems using the CAP protocol.",
		),
	)

	connctd := container.NewVBox(widget.NewLabel("Connected!"))
	tabs := container.NewAppTabs(about_tab)
	login_tab := NewLoginTab("Login", "NETL SuperComputer", cfg.Joule_Ips, conn_man,
		func(conn cap.Connection) {
			services, err := FindServices()
			if err != nil {
				return
			}
			for _, service := range services {
				if service.Name == "joule" {
					joule := NewJouleConnected(a, cfg, conn_man)
					tabs.Append(joule.CapTab.Tab)
				}
				if service.Name == "watt" {
					watt := NewWattConnected(a, cfg, conn_man)
					tabs.Append(watt.CapTab.Tab)
				}
				if service.Name == "fe261" {
					fe261 := NewFe261Connected(a, cfg, conn_man)
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
