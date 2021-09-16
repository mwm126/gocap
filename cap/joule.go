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

	vcard := widget.NewCard("GUI", "TODO", nil)
	vncTab := container.NewTabItem("VNC", vcard)

	cfg := GetConfig()
	fwdTab := newPortForwardTab(app, cfg.Joule_Forwards, func(fwds []string) {
		cfg := GetConfig()
		cfg.Joule_Forwards = fwds[2:]
		WriteConfig(cfg)
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

func saveForwards(fwds []string) {
	cfg := GetConfig()
	cfg.Joule_Forwards = fwds[2:]
	WriteConfig(cfg)
}
