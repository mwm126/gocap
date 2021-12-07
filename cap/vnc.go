package cap

import (
	"aeolustec.com/capclient/cap/connection"
	"fmt"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"log"
	"strings"
)

type VncTab struct {
	TabItem        *container.TabItem
	app            fyne.App
	connection     connection.Connection
	refresh_btn    *widget.Button
	new_btn        *widget.Button
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
		session_labels: binding.BindStringList(&[]string{}),
	}

	sessions := widget.NewListWithData(t.session_labels,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(fwd binding.DataItem, obj fyne.CanvasObject) {
			obj.(*widget.Label).Bind(fwd.(binding.String))
		})

	t.refresh_btn = widget.NewButton("Refresh", func() { t.refresh() })
	t.new_btn = widget.NewButton("New VNC Session", t.showNewVncSessionDialog)
	vcard := widget.NewCard("GUI", "List of VNC Sessions", t.new_btn)
	box := container.NewBorder(vcard, t.refresh_btn, nil, nil, sessions)

	t.TabItem = container.NewTabItem("VNC", box)
	return &t
}

func (t *VncTab) showNewVncSessionDialog() {
	win := t.app.NewWindow("Add Vnc Session")
	DEFAULT_RESOLUTIONS := []string{"800x600", "1024x768", "1280x1024", "1600x1200"}
	f := t.NewVncSessionForm(win, DEFAULT_RESOLUTIONS)
	win.SetContent(f.Form)
	win.Show()
}

type VncSessionForm struct {
	Form          *widget.Form
	preset_select *widget.Select
	xres_entry    *widget.Entry
	yres_entry    *widget.Entry
}

func (t *VncTab) NewVncSessionForm(win fyne.Window, rezs []string) *VncSessionForm {
	xres_entry := widget.NewEntry()
	yres_entry := widget.NewEntry()
	preset_select := widget.NewSelect(
		rezs, func(text string) {
			res := strings.Split(text, "x")
			xres_entry.SetText(res[0])
			yres_entry.SetText(res[1])
		})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Presets", Widget: preset_select},
			{Text: "X-resolution", Widget: xres_entry},
			{Text: "Y-resolution", Widget: yres_entry},
		},
		OnSubmit: func() {
			t.connection.CreateVncSession(xres_entry.Text, yres_entry.Text)
			win.Close()
		},
		OnCancel:   func() { win.Close() },
		SubmitText: "Create Session",
		CancelText: "Cancel",
	}
	vsf := &VncSessionForm{
		form,
		preset_select,
		xres_entry,
		yres_entry,
	}
	return vsf
}
