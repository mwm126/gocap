package joule

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

type JouleTab struct {
	app    fyne.App
	Tabs   *container.AppTabs
	CapTab login.CapTab
}

func NewJouleConnected(
	app fyne.App,
	service login.Service,
	conn_man cap.ConnectionManager,
) JouleTab {
	var joule_tab JouleTab
	tabs := container.NewAppTabs()
	cont := container.NewMax(tabs)

	joule_tab = JouleTab{
		app,
		tabs,
		login.NewCapTab("Joule", "NETL SuperComputer", service, conn_man,
			func(conn cap.Connection) {
				joule_tab.Connect(conn)
			}, cont),
	}
	return joule_tab
}

func (t *JouleTab) Connect(conn cap.Connection) {
	homeTab := newJouleHome(t.CapTab.CloseConnection)
	sshTab := ssh.NewSsh(conn)
	vncTab := newVncTab(t.app, conn)
	vncTabItem := vncTab.TabItem

	cfg := config.GetConfig()
	fwdTab := forwards.NewPortForwardTab(t.app, cfg.Joule_Forwards, func(fwds []string) {
		conn.UpdateForwards(fwds)
		config.SaveForwards(fwds)
	})

	t.Tabs.SetItems([]*container.TabItem{homeTab, sshTab, vncTabItem, fwdTab.TabItem})
}

func newJouleHome(close_cb func()) *container.TabItem {
	close := widget.NewButton("Disconnect", close_cb)
	box := container.NewVBox(widget.NewLabel("Connected!"), close)
	return container.NewTabItem("Home", box)
}
