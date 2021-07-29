package cap

import (
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

func NewWattTab(knocker Knocker) WattTab {
	username := widget.NewEntry()
	username.SetPlaceHolder("Enter username...")
	password := widget.NewEntry()
	password.SetPlaceHolder("Enter password...")
	network := widget.NewSelect([]string{"external"}, func(s string) {})
	login := widget.NewButton("Login", func() {
		knocker.Knock(username.Text, password.Text, network.Selected)
	})
	wattContent := container.NewVBox(username, password, network, login)
	tab := container.NewTabItem("Watt", wattContent)
	watt := WattTab{tab, network, username, password, login}
	return watt
}
