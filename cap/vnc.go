package cap

import (
	"aeolustec.com/capclient/cap/connection"

	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Session struct {
	Username      string
	DisplayNumber string
	Geometry      string
	DateCreated   string
	HostAddress   string
	HostPort      string
}

type VncTab struct {
	// app                fyne.App
	connection_manager *connection.CapConnectionManager
	refresh_btn        *widget.Button
	new_vnc            *widget.Button
	// save     SaveCallback
	session_labels binding.StringList
	sessions       []Session
}

func (vt *VncTab) refresh() error {
	// refreshes sessions attribute

	b, err := connection.CleanExec(vt.connection_manager.GetConnection().GetClient(), "ps auxnww|grep Xvnc|grep -v grep")

	vt.sessions = findSessions(vt.connection_manager.GetUsername(), string(b))
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

func newVncTab(conn_man *connection.CapConnectionManager) *VncTab {
	t := VncTab{
		connection_manager: conn_man,
		// save: cb,
		session_labels: binding.BindStringList(
			&[]string{
				"1920x1080 :4   2021-12-21",
				"800x600   :5   2021-12-21",
			},
		),
	}

	t.refresh_btn = widget.NewButton("Refresh", func() { t.refresh() })
	t.new_vnc = widget.NewButton("New VNC Session", func() { run_ssh(conn_man) })
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

func get_field(fields []string, fieldname string) string {
	for ii, field := range fields {
		if field == fieldname {
			return fields[ii+1]
		}
	}
	return ""
}

func findSessions(username string, text string) []Session {
	sessions := make([]Session, 0, 10)
	for _, line := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		session, err := parseVncLine(line)
		if err != nil {
			continue
		}
		if session.Username == username {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

func parseVncLine(line string) (Session, error) {
	fields := strings.Fields(line)
	username := fields[15][1 : len(fields[15])-1]
	session := Session{
		Username:      username,
		DisplayNumber: fields[11],
		Geometry:      get_field(fields, "-geometry"),
		DateCreated:   fields[8],
		HostAddress:   "localhost",
		HostPort:      get_field(fields, "-rfbport"),
	}
	return session, nil
}
