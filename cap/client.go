package cap

import (
	"aeolustec.com/capclient/cap/connection"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Client represents the Main window of CAP client
type Client struct {
	window fyne.Window
	app    fyne.App
}

func NewClient(conn_man connection.ConnectionManager) Client {
	a := app.New()
	w := a.NewWindow("CAP Client")

	cfg := GetConfig()
	about_tab := container.NewTabItemWithIcon(
		"About",
		theme.HomeIcon(),
		widget.NewLabel(
			"The CAP client is used for connecting to Joule, Watt, and other systems using the CAP protocol.",
		),
	)
	tabs := container.NewAppTabs(about_tab)

	if cfg.Enable_joule {
		joule := NewJouleConnected(a, cfg, conn_man)
		tabs.Append(joule.CapTab.Tab)
	}
	if cfg.Enable_watt {
		watt := NewWattConnected(a, cfg, conn_man)
		tabs.Append(watt.CapTab.Tab)
	}
	if cfg.Enable_fe261 {
		fe261 := NewFe261Connected(a, cfg, conn_man)
		tabs.Append(fe261.CapTab.Tab)
	}

	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)

	return Client{w, a}
}

func (client *Client) Run() {
	client.window.ShowAndRun()
}
