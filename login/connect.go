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

type LoginInfo struct {
	Network  string
	Username string
	Password string
}

type CapTab struct {
	Tab                *container.TabItem
	connection_manager *cap.ConnectionManager
	connection         *cap.Connection
	LoginInfo          LoginInfo
	ConnectBtn         *widget.Button
	card               *widget.Card
	login              *fyne.Container
	connecting         *fyne.Container
	change_password    *fyne.Container
	pw_expired_cb      func(cap.Client)
	ConnectedCallback  func(cap.Connection)
}

func NewCapTab(tabname,
	desc string,
	service Service,
	conn_man *cap.ConnectionManager,
	connected_cb func(cap *cap.Connection),
	connected *fyne.Container, login_info LoginInfo) CapTab {
	tab := &CapTab{}
	connect_cancelled := false
	ch := make(chan string)

	port := service.CapPort
	tab.connection_manager = conn_man
	tab.LoginInfo = login_info
	tab.connection_manager.AddYubikeyCallback(func(enable bool) {
		if enable {
			tab.Enable()
		} else {
			tab.Disable()
		}
	})

	tab.login = tab.NewLogin(service, func(user, pass string, ext_ip, srv_ip net.IP) {
		tab.card.SetContent(tab.connecting)

		conn, err := tab.connection_manager.Connect(user, pass, ext_ip, srv_ip,
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
			conn.Close()
			connect_cancelled = false
			return
		}

		time.Sleep(1 * time.Second)
		tab.card.SetContent(connected)

		if tab.connection == nil {
			return
		}
		connected_cb(tab.connection)
	})

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

func (t *CapTab) NewLogin(
	service Service,
	connect_cb func(user, pass string, ext_ip, srv_ip net.IP)) *fyne.Container {

	network_ips := make(map[string]string)
	external_ips := make(map[string]string)
	for name, val := range service.Networks {
		network_ips[name] = val.CapServerAddress
		external_ips[name] = val.ClientExternalAddress
	}
	t.ConnectBtn = widget.NewButton("Connect", func() {
		var ext_addr, server_addr net.IP
		if t.LoginInfo.Network == "external" {
			ext_addr = GetExternalIp()
		} else {
			ext_addr = net.ParseIP(external_ips[t.LoginInfo.Network])
		}
		server_addr = net.ParseIP(network_ips[t.LoginInfo.Network])
		go connect_cb(t.LoginInfo.Username, t.LoginInfo.Password, ext_addr, server_addr)
	})
	return container.NewVBox(t.ConnectBtn)
}

func (t *CapTab) CloseConnection() {
	if t.connection == nil {
		log.Println("No connection connection; cannot close connection")
		return
	}
	defer t.connection.Close()
	t.connection = nil
	t.card.SetContent(t.login)
}

func (t *CapTab) Disable() {
	t.ConnectBtn.Disable()
}

func (t *CapTab) Enable() {
	t.ConnectBtn.Enable()
}
