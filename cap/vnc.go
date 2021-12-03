package cap

import (
	"aeolustec.com/capclient/cap/connection"

	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type VncTab struct {
	// app                fyne.App
	connection  connection.Connection
	refresh_btn *widget.Button
	new_vnc     *widget.Button
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
	// vt.session_labels.Set(labels)
	return err
}

func newVncTab(conn connection.Connection) *VncTab {
	t := VncTab{
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
	t.new_vnc = widget.NewButton("New VNC Session", func() { run_ssh(conn) })
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

	vcard := widget.NewCard("GUI", "List of VNC Sessions", t.new_vnc)
	box := container.NewBorder(vcard, t.refresh_btn, nil, nil, sessions)

	return container.NewTabItem("VNC", box)
}
