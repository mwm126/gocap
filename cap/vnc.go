package cap

import (
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
	app                fyne.App
	connection_manager *CapConnectionManager
	// save     SaveCallback
	session_labels binding.StringList
	sessions       []Session
}

func (vt *VncTab) refresh() error {
	// refreshes sessions attribute
	session, err := vt.connection_manager.connection.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	b, err := session.CombinedOutput("ps auxnww|grep Xvnc|grep -v grep")
	vt.sessions = findSessions(vt.connection_manager.connection.connectionInfo.username, string(b))
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

func newVncTab(app fyne.App, conn_man *CapConnectionManager) *container.TabItem {
	t := VncTab{
		app:                app,
		connection_manager: conn_man,
		// save: cb,
		session_labels: binding.BindStringList(
			&[]string{
				"1920x1080 :4   2021-12-21",
				"800x600   :5   2021-12-21",
			},
		),
	}

	new_vnc := widget.NewButton("New VNC Session", func() { run_ssh(conn_man) })
	refresh_btn := widget.NewButton("Refresh", func() { t.refresh() })

	sessions := widget.NewListWithData(t.session_labels,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(fwd binding.DataItem, obj fyne.CanvasObject) {
			obj.(*widget.Label).Bind(fwd.(binding.String))
		})

	vcard := widget.NewCard("GUI", "List of VNC Sessions", new_vnc)
	box := container.NewBorder(vcard, refresh_btn, nil, nil, sessions)

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
