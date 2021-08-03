package cap

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"time"
)

type WattTab struct {
	Tab           *container.TabItem
	NetworkSelect *widget.Select
	UsernameEntry *widget.Entry
	PasswordEntry *widget.Entry
	LoginBtn      *widget.Button
	connection    *CapConnection
}

func NewWattTab(knocker Knocker, a fyne.App) WattTab {
	var watt *WattTab
	username := widget.NewEntry()
	username.SetPlaceHolder("Enter username...")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Enter password...")
	network := widget.NewSelect(getNetworkNames(), func(s string) {})
	networks := getNetworks()
	network.SetSelected("external")

	var card *widget.Card

	var wattLogin, wattConnecting, wattConnected *fyne.Container
	login := widget.NewButton("Login", func() {
		card.SetContent(wattConnecting)
		conn, err := newCapConnection(username.Text, password.Text, networks[network.Selected].WattAddress, knocker)
		if err != nil {
			log.Println("Unable to make CAP Connection")
			card.SetContent(wattLogin)
			return
		}
		watt.connection = conn
		time.Sleep(1 * time.Second)
		card.SetContent(wattConnected)
	})
	wattLogin = container.NewVBox(username, password, network, login)

	connecting := widget.NewLabel("Connecting......")
	cancel := widget.NewButton("Cancel", func() {
		card.SetContent(wattLogin)
	})
	wattConnecting = container.NewVBox(connecting, cancel)

	ssh := widget.NewButton("Connect SSH", func() {
		ConnectSsh(a, "localhost", username.Text, password.Text)
	})
	close := widget.NewButton("Close", func() {
		card.SetContent(wattLogin)
	})
	wattConnected = container.NewVBox(widget.NewLabel("Connected!"), ssh, close)

	card = widget.NewCard("Login to Watt", "(NETL Machine Learning system)", wattLogin)

	tab := container.NewTabItem("Watt", card)
	watt = &WattTab{tab, network, username, password, login, nil}
	return *watt
}
