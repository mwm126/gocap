package cap

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type VncTab struct {
	app fyne.App
	// save     SaveCallback
	sessions binding.StringList
}

func newVncTab(app fyne.App, conn_man *CapConnectionManager) *container.TabItem {

	t := VncTab{
		app: app,
		// save: cb,
		sessions: binding.BindStringList(
			&[]string{
				"1920x1080 :4   2021-12-21",
				"800x600   :5   2021-12-21",
			},
		),
	}

	new_vnc := widget.NewButton("New VNC Session", func() { run_ssh(conn_man) })

	sessions := widget.NewListWithData(t.sessions,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(fwd binding.DataItem, obj fyne.CanvasObject) {
			obj.(*widget.Label).Bind(fwd.(binding.String))
		})

	box := container.NewVBox(
		widget.NewLabel("To create new Terminal (SSH) Session in gnome-terminal:"),
		new_vnc,
		sessions,
	)

	vcard := widget.NewCard("GUI", "List of VNC Sessions", box)

	return container.NewTabItem("VNC", vcard)
}
