package cap

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewJouleConnected(app fyne.App,
	conn_man *CapConnectionManager,
	close_cb func()) *fyne.Container {

	homeTab := newJouleHome(close_cb)
	sshTab := newSsh(conn_man)
	vncTab := newVncTab(app, conn_man)

	cfg := GetConfig()
	conn := conn_man.connection
	fwdTab := newPortForwardTab(app, cfg.Joule_Forwards, func(fwds []string) {
		conn.UpdateForwards(fwds)
		SaveForwards(fwds)
	})

	tabs := container.NewAppTabs(
		homeTab,
		sshTab,
		vncTab,
		fwdTab,
	)
	return container.NewMax(tabs)
}

func newJouleHome(close_cb func()) *container.TabItem {
	close := widget.NewButton("Disconnect", close_cb)
	box := container.NewVBox(widget.NewLabel("Connected!"), close)
	return container.NewTabItem("Home", box)
}
