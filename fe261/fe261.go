package fe261

import (
	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/login"
	"aeolustec.com/capclient/ssh"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Fe261Tab struct {
	app    fyne.App
	Tabs   *container.AppTabs
	CapTab *login.CapTab
}

func NewFe261Connected(
	app fyne.App,
	service login.Service,
	conn_man *cap.ConnectionManager,
	login_info login.LoginInfo) Fe261Tab {
	var fe261_tab Fe261Tab
	tabs := container.NewAppTabs()
	cont := container.NewMax(tabs)

	fe261_tab = Fe261Tab{
		app,
		tabs,
		login.NewCapTab("FE261", "NETL SuperComputer", service, conn_man,
			func(conn *cap.Connection) {
				fe261_tab.Connect(conn)
			}, cont, login_info),
	}
	return fe261_tab
}

func (t *Fe261Tab) Connect(conn *cap.Connection) {
	homeTab := newFe261Home(t.CapTab.CloseConnection)
	sshTab := ssh.NewSsh(conn)
	t.Tabs.SetItems([]*container.TabItem{homeTab, sshTab})
}

func newFe261Home(close_cb func()) *container.TabItem {
	close := widget.NewButton("Disconnect", close_cb)
	box := container.NewVBox(widget.NewLabel("Connected!"), close)
	return container.NewTabItem("Home", box)
}
