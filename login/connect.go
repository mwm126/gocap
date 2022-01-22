package login

import (
	"log"

	"aeolustec.com/capclient/cap"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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
	ConnectForm        *ConnectForm
	LoginInfo          LoginInfo
	LoginCard          *LoginCard
}

func NewCapTab(tabname,
	desc string,
	service cap.Service,
	conn_man *cap.ConnectionManager,
	connected_cb func(cap *cap.Connection),
	connected *fyne.Container, login_info LoginInfo) *CapTab {

	var tab CapTab

	tab.connection_manager = conn_man
	tab.LoginInfo = login_info
	tab.ConnectForm = NewConnectForm(service, login_info, func(linfo LoginInfo) {
		conn, err := tab.LoginCard.handle_login(service, linfo, connected)
		if err != nil {
			log.Println("Could not login to service ", service, " because ", err)
			return
		}
		tab.connection = conn
		connected_cb(tab.connection)
	})
	tab.connection_manager.AddYubikeyCallback(tab.ConnectForm.setEnabled)
	tab.LoginCard = NewLoginCard(conn_man, tabname, desc, service, tab.ConnectForm.Container)
	tab.Tab = container.NewTabItem(tabname, tab.LoginCard.Card)
	return &tab
}

func (t *CapTab) CloseConnection() {
	if t.connection == nil {
		log.Println("No connection connection; cannot close connection")
		return
	}
	defer t.connection.Close()
	t.connection = nil
	t.LoginCard.Card.SetContent(t.ConnectForm.Container)
}
