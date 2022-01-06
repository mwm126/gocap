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
	window fyne.Window
	Tabs   *container.AppTabs
	CapTab *login.CapTab
	vncTab *VncTab
}

func NewJouleConnected(
	app fyne.App,
	w fyne.Window,
	service login.Service,
	conn_man *cap.ConnectionManager,
	login_info login.LoginInfo) JouleTab {
	var joule_tab JouleTab
	tabs := container.NewAppTabs()
	cont := container.NewMax(tabs)

	joule_tab = JouleTab{
		app,
		w,
		tabs,
		login.NewCapTab("Joule", "NETL SuperComputer", service, conn_man,
			func(conn *cap.Connection) {
				joule_tab.Connect(conn)
			}, cont, login_info),
		nil,
	}
	return joule_tab
}

func (t *JouleTab) Connect(conn *cap.Connection) {
	homeTab := newJouleHome(
		func() {
			t.vncTab.Close()
			t.CapTab.CloseConnection()
		})
	sshTab := ssh.NewSsh(conn)
	t.vncTab = newVncTab(t.app, t.window, conn, &ExeRunner{}, cap.FreePortFinder{})
	vncTabItem := t.vncTab.TabItem

	cfg := config.GetConfig()
	fwdTab := forwards.NewPortForwardTab(t.app, cfg.Joule_Forwards, func(fwds []string) {
		conn.UpdateForwards(fwds)
		config.SaveForwards(fwds)
	})

	t.Tabs.SetItems([]*container.TabItem{homeTab, vncTabItem, sshTab, fwdTab.TabItem})
}

func newJouleHome(close_cb func()) *container.TabItem {
	close := widget.NewButton("Disconnect", close_cb)
	box := container.NewVBox(widget.NewLabel("Connected!"), close)
	return container.NewTabItem("Home", box)
}
