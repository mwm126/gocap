package cap

import (
	"aeolustec.com/capclient/cap/connection"

	"bufio"
	"errors"
	"log"
	"net"
	"os"
	"time"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type CapTab struct {
	Tab                *container.TabItem
	connection_manager connection.ConnectionManager
	networkSelect      *widget.Select
	usernameEntry      *widget.Entry
	passwordEntry      *widget.Entry
	loginBtn           *widget.Button
	card               *widget.Card
	login              *fyne.Container
	connecting         *fyne.Container
	change_password    *fyne.Container
	pw_expired_cb      func(connection.PasswordChecker)
}

func NewCapTab(tabname,
	desc string,
	ips map[string]string,
	conn_man connection.ConnectionManager,
	connected_cb func(connection connection.Connection),
	connected *fyne.Container) CapTab {
	tab := &CapTab{}
	connect_cancelled := false
	ch := make(chan string)

	tab.connection_manager = conn_man
	tab.login = tab.NewLogin(ips, func(user, pass string, ext_ip, srv_ip net.IP) {
		tab.card.SetContent(tab.connecting)
		cfg := GetConfig()
		port := cfg.CapPort
		err := tab.connection_manager.Connect(user, pass, ext_ip, srv_ip,
			port,
			tab.pw_expired_cb, ch)

		if err != nil {
			log.Println("Unable to make CAP Connection")
			tab.card.SetContent(tab.login)
			connect_cancelled = false
			return
		}

		if tab.connection_manager.GetPasswordExpired() {
			tab.card.SetContent(tab.change_password)
			return
		}

		if connect_cancelled {
			log.Println("CAP Connection cancelled.")
			tab.connection_manager.Close()
			connect_cancelled = false
			return
		}

		time.Sleep(1 * time.Second)
		tab.card.SetContent(connected)
		connected_cb(conn_man.GetConnection())
	})

	tab.pw_expired_cb = func(pw_checker connection.PasswordChecker) {
		// Detected expired password callback
		tab.connection_manager.SetPasswordExpired()
		tab.card.SetContent(tab.change_password)
	}
	tab.connecting = func() *fyne.Container {
		connecting := widget.NewLabel("Connecting......")
		cancel := widget.NewButton("Cancel", func() {
			connect_cancelled = true
			tab.card.SetContent(tab.login)
		})
		return container.NewVBox(connecting, cancel)
	}()

	tab.change_password = NewChangePassword(func(new_password string) {
		ch <- new_password
		connect_cancelled = true
		tab.card.SetContent(tab.login)
	})
	tab.card = widget.NewCard(tabname, desc, tab.login)

	tab.Tab = container.NewTabItem(tabname, tab.card)
	tab.connection_manager = conn_man
	return *tab
}

func (t *CapTab) NewLogin(network_ips map[string]string,
	connect_cb func(user, pass string, ext_ip, srv_ip net.IP)) *fyne.Container {
	username := widget.NewEntry()
	username.SetPlaceHolder("Enter username...")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Enter password...")

	uname, pword, err := GetSavedLogin()
	if err == nil {
		username.SetText(uname)
		password.SetText(pword)
	}

	networkNames := make([]string, 0, len(network_ips))
	for network := range network_ips {
		networkNames = append(networkNames, network)
	}
	network := widget.NewSelect(networkNames, func(s string) {})

	network.SetSelected("external")
	login := widget.NewButton("Login", func() {
		var ext_addr, server_addr net.IP
		if network.Selected == "external" {
			ext_addr = GetExternalIp()
		} else {
			ext_addr = net.ParseIP(GetConfig().External_Ips[network.Selected])
		}
		server_addr = net.ParseIP(network_ips[network.Selected])
		go connect_cb(username.Text, password.Text, ext_addr, server_addr)
	})
	t.networkSelect = network
	t.usernameEntry = username
	t.passwordEntry = password
	t.loginBtn = login
	return container.NewVBox(username, password, network, login)
}

func (t *CapTab) closeConnection() {
	t.connection_manager.Close()
	t.card.SetContent(t.login)
}

func GetSavedLogin() (string, string, error) {
	file, err := os.Open(".cap-credentials")
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", "", err
	}
	if len(lines) < 2 {
		return "", "", errors.New("Could not read .cap-credentials")
	}
	return lines[0], lines[1], nil
}
