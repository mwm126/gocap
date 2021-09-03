package main

import (
	"aeolustec.com/capclient/cap"
	"crypto/rand"
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
	app      fyne.App
}

func newClient(knocker cap.Knocker) Client {
	a := app.New()
	w := a.NewWindow("CAP Client")

	joule := cap.NewJouleTab(knocker, a)
	watt := cap.NewWattTab(knocker, a)

	tabs := container.NewAppTabs(
		joule.Tab,
		watt.Tab,
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	w.SetContent(tabs)

	tabs.Append(container.NewTabItemWithIcon("About", theme.HomeIcon(), widget.NewLabel("The CAP client is used for connecting to Joule, Watt, and other systems using the CAP protocol.")))
	return Client{w, joule, watt, knocker, a}

}

func main() {
	var entropy [32]byte
	rand.Read(entropy[:])

	yk := cap.UsbYubikey{}
	knk := cap.NewPortKnocker(&yk, entropy)

	client := newClient(&knk)
	client.window.ShowAndRun()
}
