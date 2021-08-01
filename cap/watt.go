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

func NewWattTab(knocker Knocker) WattTab {
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
		time.Sleep(2 * time.Second)
		card.SetContent(wattConnected)
	})

	wattLogin = container.NewVBox(username, password, network, login)
	close := widget.NewButton("Close", func() {
		card.SetContent(wattLogin)
	})

	connecting := widget.NewLabel("Connecting......")
	wattConnecting = container.NewVBox(connecting)

	wattConnected = container.NewVBox(widget.NewLabel("Connected!"), close)

	card = widget.NewCard("Login to Watt", "(NETL Machine Learning system)", wattLogin)

	tab := container.NewTabItem("Watt", card)
	watt = &WattTab{tab, network, username, password, login, nil}
	return *watt
}
