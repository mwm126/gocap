package cap

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewFe261Connected(app fyne.App,
	conn_man *CapConnectionManager,
	close_cb func()) *fyne.Container {

	homeTab := newHome(close_cb)
	sshTab := newSsh(conn_man)

	cfg := GetConfig()
	fwdTab := newPortForwardTab(app, cfg.Fe261_Forwards, func(fwds []string) {
		cfg := GetConfig()
		cfg.Fe261_Forwards = fwds[2:]
		WriteConfig(cfg)
	})

	tabs := container.NewAppTabs(
		homeTab,
		sshTab,
		fwdTab,
	)
	return container.NewMax(tabs)
}

func newHome(close_cb func()) *container.TabItem {
	close := widget.NewButton("Disconnect", func() {
		close_cb()
	})
	box := container.NewVBox(widget.NewLabel("Connected!"), close)
	return container.NewTabItem("Home", box)
}
