package cap

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"time"
)

type JouleTab struct {
	Tab        *container.TabItem
	connection *CapConnection
}

func NewJouleTab(knocker Knocker, a fyne.App) JouleTab {
	joule := &JouleTab{}
	var card *widget.Card
	var jouleLogin, jouleConnecting, jouleConnected *fyne.Container
	connect_cancelled := false

	jouleLogin = NewJouleLogin(func(user, pass, host string) {
		card.SetContent(jouleConnecting)
		conn, err := newCapConnection(user, pass, host, knocker)

		if err != nil {
			log.Println("Unable to make CAP Connection")
			card.SetContent(jouleLogin)
			connect_cancelled = false
			return
		}

		if connect_cancelled {
			log.Println("CAP Connection cancelled.")
			conn.close()
			connect_cancelled = false
			return
		}

		joule.connection = conn
		time.Sleep(1 * time.Second)
		card.SetContent(jouleConnected)
	})

	jouleConnecting = NewJouleConnecting(func() {
		connect_cancelled = true
		card.SetContent(jouleLogin)
	})

	jouleConnected = NewJouleConnected(a, joule, func() {
		card.SetContent(jouleLogin)
	})

	card = widget.NewCard("Login to Joule", "(NETL Machine Learning system)", jouleLogin)

	joule.Tab = container.NewTabItem("Joule", card)
	return *joule
}

func NewJouleLogin(connect_cb func(user, pass, host string)) *fyne.Container {
	username := widget.NewEntry()
	username.SetPlaceHolder("Enter username...")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Enter password...")
	network := widget.NewSelect(getNetworkNames(), func(s string) {})
	networks := getNetworks()
	network.SetSelected("external")
	login := widget.NewButton("Login", func() {
		go connect_cb(username.Text, password.Text, networks[network.Selected].JouleAddress)
	})
	return container.NewVBox(username, password, network, login)
}

func NewJouleConnecting(cancel_cb func()) *fyne.Container {
	connecting := widget.NewLabel("Connecting......")
	cancel := widget.NewButton("Cancel", func() {
		cancel_cb()
	})
	return container.NewVBox(connecting, cancel)
}

func NewJouleConnected(a fyne.App, joule *JouleTab, close_cb func()) *fyne.Container {
	var card *widget.Card
	ssh := widget.NewButton("Connect SSH", func() {
		log.Println(joule)
		log.Println(joule.connection)
		log.Println(joule.connection.connectionInfo)
	})
	card = widget.NewCard("Connect SSH", "(NETL Machine Learning system)", ssh)
	close := widget.NewButton("Close", func() {
		close_cb()
	})
	jouleConnected := container.NewVBox(widget.NewLabel("Connected!"), card, close)
	return jouleConnected
}
