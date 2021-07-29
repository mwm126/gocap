package main

import (
	"aeolustec.com/capclient/cap"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Client represents the Main window of CAP client
type Client struct {
	window   fyne.Window
	jouleTab cap.JouleTab
	wattTab  cap.WattTab
	knocker  cap.Knocker
}

func newClient(knocker cap.Knocker) Client {
	a := app.New()
	w := a.NewWindow("CAP Client")
	w.Resize(fyne.NewSize(300, 200))

	joule := cap.NewJouleTab(knocker)
	watt := cap.NewWattTab(knocker)

	tabs := container.NewAppTabs(
		joule.Tab,
		watt.Tab,
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)

	tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("The CAP client is used for connecting to Joule, Watt, and other systems using the CAP protocol.")))
	return Client{w, joule, watt, knocker}

}

func main() {
	knk := cap.PortKnocker{}
	client := newClient(&knk)
	client.window.ShowAndRun()
}
