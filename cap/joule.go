package cap

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type JouleTab struct {
	Tab           *container.TabItem
	connection    *CapConnection
	networkSelect *widget.Select
	usernameEntry *widget.Entry
	passwordEntry *widget.Entry
	loginBtn      *widget.Button
}

func NewJouleTab(knocker Knocker, a fyne.App) JouleTab {
	joule := &JouleTab{}
	var card *widget.Card
	var jouleLogin, jouleConnecting, jouleConnected *fyne.Container
	connect_cancelled := false

	jouleLogin = joule.NewJouleLogin(func(user, pass string, host net.IP) {
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
		joule.connection.close()
		joule.connection = nil
		card.SetContent(jouleLogin)
	})

	card = widget.NewCard("Joule 2.0", "NETL Supercomputer", jouleLogin)

	joule.Tab = container.NewTabItem("Joule", card)
	return *joule
}

func (t *JouleTab) NewJouleLogin(connect_cb func(user, pass string, host net.IP)) *fyne.Container {
	username := widget.NewEntry()
	username.SetPlaceHolder("Enter username...")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Enter password...")

	cfg := GetConfig()
	networkNames := make([]string, 0, len(cfg.Joule_Ips))
	for network := range cfg.Joule_Ips {
		networkNames = append(networkNames, network)
	}
	network := widget.NewSelect(networkNames, func(s string) {})

	network.SetSelected("external")
	login := widget.NewButton("Login", func() {
		var addr net.IP
		if network.Selected == "external" {
			addr = GetExternalIp()
		} else {
			addr = net.ParseIP(cfg.Joule_Ips[network.Selected])
		}
		go connect_cb(username.Text, password.Text, addr)
	})
	t.networkSelect = network
	t.usernameEntry = username
	t.passwordEntry = password
	t.loginBtn = login
	return container.NewVBox(username, password, network, login)
}

func NewJouleConnecting(cancel_cb func()) *fyne.Container {
	connecting := widget.NewLabel("Connecting......")
	cancel := widget.NewButton("Cancel", func() {
		cancel_cb()
	})
	return container.NewVBox(connecting, cancel)
}

func NewJouleConnected(app fyne.App, joule *JouleTab, close_cb func()) *fyne.Container {

	homeTab := newJouleHome(close_cb)
	sshTab := newJouleSsh()

	vcard := widget.NewCard("GUI", "TODO", nil)
	vncTab := container.NewTabItem("VNC", vcard)

	cfg := GetConfig()
	fwdTab := newPortForwardTab(app, cfg.Joule_Forwards, func(fwds []string) {
		cfg := GetConfig()
		cfg.Joule_Forwards = fwds[2:]
		WriteConfig(cfg)
	})

	tabs := container.NewAppTabs(
		homeTab,
		sshTab,
		vncTab,
		fwdTab,
	)
	return container.NewMax(tabs)
}

func newJouleHome(close_cb func()) *container.TabItem {
	close := widget.NewButton("Disconnect", func() {
		close_cb()
	})
	box := container.NewVBox(widget.NewLabel("Connected!"), close)
	return container.NewTabItem("Home", box)
}

func newJouleSsh() *container.TabItem {
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

func saveForwards(fwds []string) {
	cfg := GetConfig()
	cfg.Joule_Forwards = fwds[2:]
	WriteConfig(cfg)
}
