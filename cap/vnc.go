package cap

import (
	"aeolustec.com/capclient/cap/connection"

	"fmt"
	"log"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type VncTab struct {
	app         fyne.App
	connection  connection.Connection
	refresh_btn *widget.Button
	new_btn     *widget.Button
	// save     SaveCallback
	session_labels binding.StringList
	sessions       []connection.Session
}

func (vt *VncTab) refresh() error {
	// refreshes sessions attribute
	sessions, err := vt.connection.FindSessions()
	if err != nil {
		log.Println("Warning - unable to refresh: ", err)
	}
	vt.sessions = sessions
	labels := make([]string, 0)
	for _, session := range vt.sessions {
		label := fmt.Sprintf(
			"Session %s - %s - %s",
			session.DisplayNumber,
			session.Geometry,
			session.DateCreated,
		)
		labels = append(labels, label)

	}
	vt.session_labels.Set(labels)
	return err
}

func newVncTab(a fyne.App, conn connection.Connection) *VncTab {
	if conn == nil {
		panic("Invalid")
	}
	t := VncTab{
		app:        a,
		connection: conn,
		// save: cb,
		session_labels: binding.BindStringList(
			&[]string{
				"1920x1080 :4   2021-12-21",
				"800x600   :5   2021-12-21",
			},
		),
	}

	t.refresh_btn = widget.NewButton("Refresh", func() { t.refresh() })
	t.new_btn = widget.NewButton("New VNC Session", t.showNewVncSessionDialog)
	return &t
}

func newVncTabItem(t *VncTab) *container.TabItem {
	sessions := widget.NewListWithData(t.session_labels,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(fwd binding.DataItem, obj fyne.CanvasObject) {
			obj.(*widget.Label).Bind(fwd.(binding.String))
		})

	vcard := widget.NewCard("GUI", "List of VNC Sessions", t.new_btn)
	box := container.NewBorder(vcard, t.refresh_btn, nil, nil, sessions)

	return container.NewTabItem("VNC", box)
}

func (t *VncTab) showNewVncSessionDialog() {
	win := t.app.NewWindow("Add Vnc Session")

	preset_select := widget.NewSelect(
		[]string{
			"800x600",
			"1024x768",
			"1280x1024",
			"1600x1200",
		}, func(string) {})
	xres_entry := widget.NewEntry()
	yres_entry := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Presets", Widget: preset_select},
			{Text: "X-resolution", Widget: xres_entry},
			{Text: "Y-resolution", Widget: yres_entry},
		},
		OnSubmit: func() {
			// new_fwd := fmt.Sprintf("%s:%s:%s", local_p.Text, remote_h.Text, remote_p.Text)
			// t.addPortForward(new_fwd)
			// fwds, _ := t.forwards.Get()
			// t.save(fwds)
			win.Close()
		},
		OnCancel:   func() { win.Close() },
		SubmitText: "Create Session",
		CancelText: "Cancel",
	}
	win.SetContent(form)
	win.Show()
}
