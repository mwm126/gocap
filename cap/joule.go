package cap

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"time"
)

type JouleTab struct {
	Tab           *container.TabItem
	NetworkSelect *widget.Select
	UsernameEntry *widget.Entry
	PasswordEntry *widget.Entry
	LoginBtn      *widget.Button
	connection    *CapConnection
}

func NewJouleTab(knocker Knocker, a fyne.App) JouleTab {
	var joule *JouleTab
	username := widget.NewEntry()
	username.SetPlaceHolder("Enter username...")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Enter password...")
	network := widget.NewSelect(getNetworkNames(), func(s string) {})
	networks := getNetworks()
	network.SetSelected("external")

	var card *widget.Card

	var jouleLogin, jouleConnecting, jouleConnected *fyne.Container
	login := widget.NewButton("Login", func() {
		card.SetContent(jouleConnecting)
		go func() {
			conn, err := newCapConnection(username.Text, password.Text, networks[network.Selected].JouleAddress, knocker)
			if err != nil {
				log.Println("Unable to make CAP Connection")
				card.SetContent(jouleLogin)
				return
			}
			joule.connection = conn
			time.Sleep(1 * time.Second)
			card.SetContent(jouleConnected)
		}()
	})
	jouleLogin = container.NewVBox(username, password, network, login)

	connecting := widget.NewLabel("Connecting......")
	cancel := widget.NewButton("Cancel", func() {
		card.SetContent(jouleLogin)
	})
	jouleConnecting = container.NewVBox(connecting, cancel)

	ssh := widget.NewButton("Connect SSH", func() {
		ConnectSsh(a, "localhost", username.Text, password.Text)
	})
	close := widget.NewButton("Close", func() {
		card.SetContent(jouleLogin)
	})
	jouleConnected = container.NewVBox(widget.NewLabel("Connected!"), ssh, close)

	card = widget.NewCard("Login to Joule", "(NETL Machine Learning system)", jouleLogin)

	tab := container.NewTabItem("Joule", card)
	joule = &JouleTab{tab, network, username, password, login, nil}
	return *joule
}
