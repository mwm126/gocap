package login

import (
	"aeolustec.com/capclient/cap"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"log"
	"net"
	"time"
)

type LoginTab struct {
	Tab                *container.TabItem
	connection_manager cap.ConnectionManager
	NetworkSelect      *widget.Select
	UsernameEntry      *widget.Entry
	PasswordEntry      *widget.Entry
	LoginBtn           *widget.Button
	card               *widget.Card
	login              *fyne.Container
	connecting         *fyne.Container
	change_password    *fyne.Container
	pw_expired_cb      func(cap.Client)
	ConnectedCallback  func(LoginInfo)
}

func NewLoginTab(tabname,
	desc string,
	service Service,
	conn_man cap.ConnectionManager,
	connected_cb func(login_info LoginInfo),
	connected *fyne.Container,
	username, password string) LoginTab {

	tab := &LoginTab{}
	connect_cancelled := false
	ch := make(chan string)

	tab.ConnectedCallback = connected_cb
	tab.connection_manager = conn_man
	tab.connection_manager.SetYubikeyCallback(func(serial int32) {
		if serial == 0 {
			tab.Disable()
		} else {
			tab.Enable()
		}
	})
	tab.login = tab.NewLogin(service, func(network, user, pass string, ext_ip, srv_ip net.IP) {
		tab.card.SetContent(tab.connecting)
		err := tab.connection_manager.Connect(user, pass, ext_ip, srv_ip,
			service.CapPort,
			tab.pw_expired_cb, ch)

		if err != nil {
			log.Println("Unable to make CAP Connection: ", err)
			tab.card.SetContent(tab.login)
			connect_cancelled = false
			return
		}
		log.Println("Made CAP Connection: ", tab.connection_manager.GetConnection())

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
		login_info := LoginInfo{
			Network:  network,
			Username: username,
			Password: password,
		}
		connected_cb(login_info)
	}, username, password)

	tab.pw_expired_cb = func(pw_checker cap.Client) {
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
	return *tab
}

func (t *LoginTab) NewLogin(
	service Service,
	connect_cb func(network, user, pass string, ext_ip, srv_ip net.IP),
	uname, pword string) *fyne.Container {
	username := widget.NewEntry()
	username.SetText(uname)
	username.SetPlaceHolder("Enter username...")
	password := widget.NewPasswordEntry()
	password.SetText(pword)
	password.SetPlaceHolder("Enter password...")

	network_ips := make(map[string]string)
	external_ips := make(map[string]string)
	networkNames := make([]string, 0, len(service.Networks))
	for name, val := range service.Networks {
		network_ips[name] = val.CapServerAddress
		external_ips[name] = val.ClientExternalAddress
		networkNames = append(networkNames, name)
	}
	network := widget.NewSelect(networkNames, func(s string) {})

	network.SetSelected("external")
	login := widget.NewButton("Login", func() {
		var ext_addr, server_addr net.IP
		if network.Selected == "external" {
			ext_addr = GetExternalIp()
		} else {
			ext_addr = net.ParseIP(external_ips[network.Selected])
		}
		server_addr = net.ParseIP(network_ips[network.Selected])
		go connect_cb(network.Selected, username.Text, password.Text, ext_addr, server_addr)
	})
	t.NetworkSelect = network
	t.UsernameEntry = username
	t.PasswordEntry = password
	t.LoginBtn = login
	return container.NewVBox(username, password, network, login)
}

func (t *LoginTab) CloseConnection() {
	t.connection_manager.Close()
	t.card.SetContent(t.login)
}

func (t *LoginTab) Disable() {
	t.NetworkSelect.Disable()
	t.UsernameEntry.Disable()
	t.PasswordEntry.Disable()
	t.LoginBtn.Disable()
}

func (t *LoginTab) Enable() {
	t.NetworkSelect.Enable()
	t.UsernameEntry.Enable()
	t.PasswordEntry.Enable()
	t.LoginBtn.Enable()
}
