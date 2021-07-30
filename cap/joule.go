package cap

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type JouleTab struct {
	Tab           *container.TabItem
	NetworkSelect *widget.Select
	UsernameEntry *widget.Entry
	PasswordEntry *widget.Entry
	LoginBtn      *widget.Button
}

func NewJouleTab(knocker Knocker) JouleTab {
	username := widget.NewEntry()
	username.SetPlaceHolder("Enter username...")
	password := widget.NewEntry()
	password.SetPlaceHolder("Enter password...")
	network := widget.NewSelect(getNetworkNames(), func(s string) {})
	network.SetSelected("external")
	login := widget.NewButton("Login", func() {
		knocker.Knock(username.Text, password.Text, network.Selected)
	})
	jouleContent := container.NewVBox(username, password, network, login)
	tab := container.NewTabItem("Joule", jouleContent)
	joule := JouleTab{tab, network, username, password, login}
	return joule
}
