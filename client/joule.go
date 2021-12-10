package client

import (
	"aeolustec.com/capclient/cap"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type JouleTab struct {
	app    fyne.App
	Tabs   *container.AppTabs
	CapTab CapTab
}

func NewJouleConnected(app fyne.App, cfg config, conn_man cap.ConnectionManager) JouleTab {
	var joule_tab JouleTab
	tabs := container.NewAppTabs()
	cont := container.NewMax(tabs)

	joule_tab = JouleTab{
		app,
		tabs,
		NewCapTab("Joule", "NETL SuperComputer", cfg.Joule_Ips, conn_man,
			func(conn cap.Connection) {
				joule_tab.Connect(conn)
			}, cont),
	}
	return joule_tab
}

func (t *JouleTab) Connect(conn cap.Connection) {
	homeTab := newJouleHome(t.CapTab.closeConnection)
	sshTab := newSsh(conn)
	vncTab := newVncTab(t.app, conn)
	vncTabItem := vncTab.TabItem

	cfg := GetConfig()
	fwdTab := NewPortForwardTab(t.app, cfg.Joule_Forwards, func(fwds []string) {
		conn.UpdateForwards(fwds)
		SaveForwards(fwds)
	})

	t.Tabs.SetItems([]*container.TabItem{homeTab, sshTab, vncTabItem, fwdTab.TabItem})
}

func newJouleHome(close_cb func()) *container.TabItem {
	close := widget.NewButton("Disconnect", close_cb)
	box := container.NewVBox(widget.NewLabel("Connected!"), close)
	return container.NewTabItem("Home", box)
}
