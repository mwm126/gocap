package login

import (
	"log"

	"aeolustec.com/capclient/cap"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

type LoginTab struct {
	Tab                *container.TabItem
	connection_manager *cap.ConnectionManager
	connection         *cap.Connection
	LoginForm          *LoginForm
	LoginCard          *LoginCard
}

func NewLoginTab(tabname,
	desc string,
	service Service,
	conn_man *cap.ConnectionManager,
	login_cb func(*LoginInfo, []Service),
	connected *fyne.Container,
	username, password string) *LoginTab {

	var tab LoginTab

	tab.connection_manager = conn_man

	tab.LoginForm = NewLoginForm(service, func(linfo LoginInfo) {
		conn, err := tab.LoginCard.handle_login(service, linfo, connected)
		if err != nil {
			log.Println("Could not login to lookup services: ", err)
		}
		tab.connection = conn
		services, err := FindServices()
		if err != nil {
			log.Println("Could not find services: ", err)
		}
		login_cb(&linfo, services)

	}, username, password)
	tab.connection_manager.AddYubikeyCallback(tab.LoginForm.setEnabled)
	tab.LoginCard = NewLoginCard(conn_man, tabname, desc, service, tab.LoginForm.Container)

	tab.Tab = container.NewTabItem(tabname, tab.LoginCard.Card)
	return &tab
}

func (t *LoginTab) CloseConnection() {
	if t.connection == nil {
		log.Println("No connection connection; cannot close connection")
		return
	}
	defer t.connection.Close()
	t.connection = nil
	t.LoginCard.Card.SetContent(t.LoginForm.Container)
}
