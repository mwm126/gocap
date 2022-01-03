package watt

import (
	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/config"
	"aeolustec.com/capclient/forwards"
	"aeolustec.com/capclient/login"
	"aeolustec.com/capclient/ssh"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type WattTab struct {
	app    fyne.App
	Tabs   *container.AppTabs
	CapTab *login.CapTab
}

func NewWattConnected(
	app fyne.App,
	service login.Service,
	conn_man *cap.ConnectionManager,
	login_info login.LoginInfo,
) WattTab {
	var watt_tab WattTab
	tabs := container.NewAppTabs()
	cont := container.NewMax(tabs)

	watt_tab = WattTab{
		app,
		tabs,
		login.NewCapTab("Watt", "NETL SuperComputer", service, conn_man,
			func(conn *cap.Connection) {
				watt_tab.Connect(conn)
			}, cont, login_info),
	}
	return watt_tab
}

func (t *WattTab) Connect(conn *cap.Connection) {
	homeTab := newWattHome(t.CapTab.CloseConnection)
	sshTab := ssh.NewSsh(conn)

	instanceTab := NewInstanceTab(conn.GetClient())

	cfg := config.GetConfig()
	fwdTab := forwards.NewPortForwardTab(t.app, cfg.Watt_Forwards, func(fwds []string) {
		conn.UpdateForwards(fwds)
		config.SaveForwards(fwds)
	})

	t.Tabs.SetItems([]*container.TabItem{homeTab, instanceTab.TabItem, sshTab, fwdTab.TabItem})
}

func newWattHome(close_cb func()) *container.TabItem {
	close := widget.NewButton("Disconnect", close_cb)
	box := container.NewVBox(widget.NewLabel("Connected!"), close)
	return container.NewTabItem("Home", box)
}
