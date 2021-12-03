package cap

import (
	"aeolustec.com/capclient/cap/connection"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewJouleConnected(app fyne.App,
	conn_man connection.ConnectionManager,
	close_cb func()) *fyne.Container {

	conn := conn_man.GetConnection()

	homeTab := newJouleHome(close_cb)
	sshTab := newSsh(conn)
	vncTab := newVncTab(conn)
	vncTabItem := newVncTabItem(vncTab)

	cfg := GetConfig()
	fwdTab := newPortForwardTab(app, cfg.Joule_Forwards, func(fwds []string) {
		conn.UpdateForwards(fwds)
		SaveForwards(fwds)
	})

	tabs := container.NewAppTabs(
		homeTab,
		sshTab,
		vncTabItem,
		fwdTab,
	)
	return container.NewMax(tabs)
}

func newJouleHome(close_cb func()) *container.TabItem {
	close := widget.NewButton("Disconnect", close_cb)
	box := container.NewVBox(widget.NewLabel("Connected!"), close)
	return container.NewTabItem("Home", box)
}
