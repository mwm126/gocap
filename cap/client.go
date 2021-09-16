package cap

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Client represents the Main window of CAP client
type Client struct {
	window fyne.Window
	// jouleTab JouleTab
	// wattTab  WattTab
	knocker Knocker
	app     fyne.App
}

func NewClient(knocker Knocker) Client {
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
		conn_man := NewCapConnectionManager(knocker)
		var joule CapTab
		joule = NewCapTab("Joule", "NETL SuperComputer", cfg.Joule_Ips, conn_man,
			NewJouleConnected(a, conn_man, joule.closeConnection))
		tabs.Append(joule.Tab)
	}
	if cfg.Enable_watt {
		conn_man := NewCapConnectionManager(knocker)
		var watt CapTab
		watt = NewCapTab("Watt", "NETL SuperComputer", cfg.Watt_Ips, conn_man,
			NewWattConnected(a, conn_man, watt.closeConnection))
		tabs.Append(watt.Tab)
	}
	if cfg.Enable_fe261 {
		conn_man := NewCapConnectionManager(knocker)
		var fe261 CapTab
		fe261 = NewCapTab("FE261", "NETL system", cfg.Fe261_Ips, conn_man,
			NewFe261Connected(a, conn_man, fe261.closeConnection))
		tabs.Append(fe261.Tab)
	}

	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)

	return Client{w, knocker, a}
}

func (client *Client) Run() {
	client.window.ShowAndRun()
}
