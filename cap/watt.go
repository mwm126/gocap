package cap

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type WattTab struct {
	Tab           *container.TabItem
	NetworkSelect *widget.Select
	UsernameEntry *widget.Entry
	PasswordEntry *widget.Entry
	LoginBtn      *widget.Button
}

func NewWattConnected(app fyne.App,
	conn_man *CapConnectionManager,
	close_cb func()) *fyne.Container {

	ssh := widget.NewButton("Connect SSH", func() {
	})
	close := widget.NewButton("Close", close_cb)
	wattConnected := container.NewVBox(widget.NewLabel("Connected!"), ssh, close)
	return wattConnected
}
