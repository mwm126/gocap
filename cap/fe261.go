package cap

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Fe261Tab struct {
	Tab        *container.TabItem
	connection *CapConnection
}

func NewFe261Tab(knocker Knocker, a fyne.App) Fe261Tab {
	fe261 := &Fe261Tab{}
	var card *widget.Card
	var login, connecting, connected *fyne.Container
	connect_cancelled := false

	login = NewLogin(func(user, pass, host string) {
		card.SetContent(connecting)
		conn, err := newCapConnection(user, pass, host, knocker)

		if err != nil {
			log.Println("Unable to make CAP Connection")
			card.SetContent(login)
			connect_cancelled = false
			return
		}

		if connect_cancelled {
			log.Println("CAP Connection cancelled.")
			conn.close()
			connect_cancelled = false
			return
		}

		fe261.connection = conn
		time.Sleep(1 * time.Second)
		card.SetContent(connected)
	})

	connecting = NewConnecting(func() {
		connect_cancelled = true
		card.SetContent(login)
	})

	connected = NewConnected(a, fe261, func() {
		fe261.connection.close()
		fe261.connection = nil
		card.SetContent(login)
	})

	card = widget.NewCard("FE 261 system", "fe 261", login)

	fe261.Tab = container.NewTabItem("FE261", card)
	return *fe261
}

func NewLogin(connect_cb func(user, pass, host string)) *fyne.Container {
	username := widget.NewEntry()
	username.SetPlaceHolder("Enter username...")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Enter password...")
	network := widget.NewSelect(getNetworkNames(), func(s string) {})
	networks := getNetworks()
	network.SetSelected("external")
	login := widget.NewButton("Login", func() {
		go connect_cb(username.Text, password.Text, networks[network.Selected].Fe261Address)
	})
	return container.NewVBox(username, password, network, login)
}

func NewConnecting(cancel_cb func()) *fyne.Container {
	connecting := widget.NewLabel("Connecting......")
	cancel := widget.NewButton("Cancel", func() {
		cancel_cb()
	})
	return container.NewVBox(connecting, cancel)
}

func NewConnected(app fyne.App, tab *Fe261Tab, close_cb func()) *fyne.Container {

	homeTab := newHome(close_cb)
	sshTab := newSsh()

	cfg := GetConfig()
	fwdTab := newPortForwardTab(app, cfg.Fe261_Forwards, func(fwds []string) {
		cfg := GetConfig()
		cfg.Fe261_Forwards = fwds[2:]
		WriteConfig(cfg)
	})

	tabs := container.NewAppTabs(
		homeTab,
		sshTab,
		fwdTab,
	)
	return container.NewMax(tabs)
}

func newHome(close_cb func()) *container.TabItem {
	close := widget.NewButton("Disconnect", func() {
		close_cb()
	})
	box := container.NewVBox(widget.NewLabel("Connected!"), close)
	return container.NewTabItem("Home", box)
}

func newSsh() *container.TabItem {
	ssh := widget.NewButton("New SSH Session", func() {
		cmd := exec.Command("x-terminal-emulator", "--", "ssh", "localhost", "-p", strconv.Itoa(SSH_LOCAL_PORT))
		err := cmd.Run()
		if err != nil {
			log.Println("gnome-terminal FAIL: ", err)
		}
	})
	label := widget.NewLabel(fmt.Sprintf("or run in a terminal:  ssh localhost -p %d", SSH_LOCAL_PORT))
	box := container.NewVBox(widget.NewLabel("To create new Terminal (SSH) Session in gnome-terminal:"), ssh, label)
	return container.NewTabItem("SSH", box)
}
