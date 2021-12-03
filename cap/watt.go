package cap

import (
	"aeolustec.com/capclient/cap/connection"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type WattTab struct {
	app    fyne.App
	Tabs   *container.AppTabs
	CapTab CapTab
}

func NewWattConnected(app fyne.App, cfg config, conn_man connection.ConnectionManager) WattTab {
	var watt_tab WattTab
	tabs := container.NewAppTabs()
	cont := container.NewMax(tabs)

	watt_tab = WattTab{
		app,
		tabs,
		NewCapTab("Watt", "NETL SuperComputer", cfg.Watt_Ips, conn_man,
			func(conn connection.Connection) {
				watt_tab.Connect(conn)
			}, cont),
	}
	return watt_tab
}

func (t *WattTab) Connect(conn connection.Connection) {
	homeTab := newWattHome(t.CapTab.closeConnection)
	sshTab := newSsh(conn)

	cfg := GetConfig()
	fwdTab := newPortForwardTab(t.app, cfg.Watt_Forwards, func(fwds []string) {
		conn.UpdateForwards(fwds)
		SaveForwards(fwds)
	})

	t.Tabs.SetItems([]*container.TabItem{homeTab, sshTab, fwdTab})
}

func newWattHome(close_cb func()) *container.TabItem {
	close := widget.NewButton("Disconnect", close_cb)
	box := container.NewVBox(widget.NewLabel("Connected!"), close)
	return container.NewTabItem("Home", box)
}
