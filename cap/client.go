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
	about_tab := container.NewTabItemWithIcon("About", theme.HomeIcon(), widget.NewLabel("The CAP client is used for connecting to Joule, Watt, and other systems using the CAP protocol."))
	tabs := container.NewAppTabs(about_tab)

	if cfg.Enable_joule {
		joule := NewJouleTab(knocker, a)
		tabs.Append(joule.Tab)
	}
	if cfg.Enable_watt {
		watt := NewWattTab(knocker, a)
		tabs.Append(watt.Tab)
	}
	if cfg.Enable_fe261 {
		fe261 := NewFe261Tab(knocker, a)
		tabs.Append(fe261.Tab)
	}

	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)

	return Client{w, knocker, a}
}

func (client *Client) Run() {
	client.window.ShowAndRun()
}
