package cap

import (
	"log"
	"net"
	"time"

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
	connection    *CapConnection
}

func NewWattTab(knocker Knocker, a fyne.App) WattTab {
	var watt *WattTab
	username := widget.NewEntry()
	username.SetPlaceHolder("Enter username...")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Enter password...")

	cfg := GetConfig()
	networkNames := make([]string, 0, len(cfg.Watt_Ips))
	for network := range cfg.Watt_Ips {
		networkNames = append(networkNames, network)
	}

	network := widget.NewSelect(networkNames, func(s string) {})
	network.SetSelected("external")

	var card *widget.Card

	var wattLogin, wattConnecting, wattConnected *fyne.Container
	login := widget.NewButton("Login", func() {
		card.SetContent(wattConnecting)
		go func() {
			var addr net.IP
			if network.Selected == "external" {
				addr = GetExternalIp()
			} else {
				addr = net.ParseIP(cfg.Watt_Ips[network.Selected])
			}
			conn, err := newCapConnection(username.Text, password.Text, addr, knocker)
			if err != nil {
				log.Println("Unable to make CAP Connection")
				card.SetContent(wattLogin)
				return
			}
			watt.connection = conn
			time.Sleep(1 * time.Second)
			card.SetContent(wattConnected)
		}()
	})
	wattLogin = container.NewVBox(username, password, network, login)

	connecting := widget.NewLabel("Connecting......")
	cancel := widget.NewButton("Cancel", func() {
		card.SetContent(wattLogin)
	})
	wattConnecting = container.NewVBox(connecting, cancel)

	ssh := widget.NewButton("Connect SSH", func() {
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
