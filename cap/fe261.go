package cap

import (
	"aeolustec.com/capclient/cap/connection"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Fe261Tab struct {
	app    fyne.App
	Tabs   *container.AppTabs
	CapTab CapTab
}

func NewFe261Connected(app fyne.App, cfg config, conn_man connection.ConnectionManager) Fe261Tab {
	var fe261_tab Fe261Tab
	tabs := container.NewAppTabs()
	cont := container.NewMax(tabs)

	fe261_tab = Fe261Tab{
		app,
		tabs,
		NewCapTab("FE261", "NETL SuperComputer", cfg.Fe261_Ips, conn_man,
			func(conn connection.Connection) {
				fe261_tab.Connect(conn)
			}, cont),
	}
	return fe261_tab
}

func (t *Fe261Tab) Connect(conn connection.Connection) {
	homeTab := newFe261Home(t.CapTab.closeConnection)
	sshTab := newSsh(conn)
	t.Tabs.SetItems([]*container.TabItem{homeTab, sshTab})
}

func newFe261Home(close_cb func()) *container.TabItem {
	close := widget.NewButton("Disconnect", close_cb)
	box := container.NewVBox(widget.NewLabel("Connected!"), close)
	return container.NewTabItem("Home", box)
}
