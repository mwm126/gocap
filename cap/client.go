package cap

import (
	"embed"
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

func NewClient(knocker Knocker, content embed.FS) Client {
	a := app.New()
	w := a.NewWindow("CAP Client")

	cfg := GetConfig()
	about_tab := container.NewTabItemWithIcon("About", theme.HomeIcon(), widget.NewLabel("The CAP client is used for connecting to Joule, Watt, and other systems using the CAP protocol."))
	tabs := container.NewAppTabs(about_tab)

	conn_man := &CapConnectionManager{}

	if cfg.Enable_joule {
		var joule CapTab
		joule = NewCapTab("Joule", "NETL SuperComputer", cfg.Joule_Ips, knocker, conn_man, a, content,
			NewJouleConnected(a, conn_man, content, joule.closeConnection))
		tabs.Append(joule.Tab)
	}
	if cfg.Enable_watt {
		var watt CapTab
		watt = NewCapTab("Watt", "NETL SuperComputer", cfg.Watt_Ips, knocker, conn_man, a, content,
			NewWattConnected(a, conn_man, content, watt.closeConnection))
		tabs.Append(watt.Tab)
	}
	if cfg.Enable_fe261 {
		var fe261 CapTab
		fe261 = NewCapTab("FE261", "NETL system", cfg.Fe261_Ips, knocker, conn_man, a, content,
			NewFe261Connected(a, conn_man, content, fe261.closeConnection))
		tabs.Append(fe261.Tab)
	}

	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)

	return Client{w, knocker, a}
}

func (client *Client) Run() {
	client.window.ShowAndRun()
}
