package cap

import (
	"log"
	"net"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type CapTab struct {
	Tab                *container.TabItem
	connection_manager *CapConnectionManager
	networkSelect      *widget.Select
	usernameEntry      *widget.Entry
	passwordEntry      *widget.Entry
	loginBtn           *widget.Button
	card               *widget.Card
	login              *fyne.Container
	connecting         *fyne.Container
	connected          *fyne.Container
}

func NewCapTab(tabname,
	desc string,
	ips map[string]string,
	conn_man *CapConnectionManager,
	connected *fyne.Container) CapTab {
	tab := &CapTab{}
	// var login, connecting *fyne.Container
	connect_cancelled := false

	tab.login = tab.NewLogin(ips, func(user, pass string, host net.IP) {
		tab.card.SetContent(tab.connecting)
		conn, err := newCapConnection(user, pass, host, conn_man.knocker)

		if err != nil {
			log.Println("Unable to make CAP Connection")
			tab.card.SetContent(tab.login)
			connect_cancelled = false
			return
		}

		if connect_cancelled {
			log.Println("CAP Connection cancelled.")
			conn.close()
			connect_cancelled = false
			return
		}

		tab.connection_manager = &CapConnectionManager{}
		tab.connection_manager.connection = conn
		time.Sleep(1 * time.Second)
		tab.card.SetContent(connected)
	})

	tab.connecting = NewConnecting(func() {
		connect_cancelled = true
		tab.card.SetContent(tab.login)
	})

	tab.card = widget.NewCard(tabname, desc, tab.login)

	tab.Tab = container.NewTabItem(tabname, tab.card)
	tab.connection_manager = conn_man
	return *tab
}

func (t *CapTab) NewLogin(network_ips map[string]string,
	connect_cb func(user, pass string, host net.IP)) *fyne.Container {
	username := widget.NewEntry()
	username.SetPlaceHolder("Enter username...")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Enter password...")

	networkNames := make([]string, 0, len(network_ips))
	for network := range network_ips {
		networkNames = append(networkNames, network)
	}
	network := widget.NewSelect(networkNames, func(s string) {})

	network.SetSelected("external")
	login := widget.NewButton("Login", func() {
		var addr net.IP
		if network.Selected == "external" {
			addr = GetExternalIp()
		} else {
			addr = net.ParseIP(network_ips[network.Selected])
		}

		go connect_cb(username.Text, password.Text, addr)
	})
	t.networkSelect = network
	t.usernameEntry = username
	t.passwordEntry = password
	t.loginBtn = login
	return container.NewVBox(username, password, network, login)
}

func NewConnecting(cancel_cb func()) *fyne.Container {
	connecting := widget.NewLabel("Connecting......")
	cancel := widget.NewButton("Cancel", func() {
		cancel_cb()
	})
	return container.NewVBox(connecting, cancel)
}

func (t *CapTab) closeConnection() {
	t.connection_manager.CloseConnection()
	t.card.SetContent(t.login)
}
