package joule

import (
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"aeolustec.com/capclient/cap"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

const VNC_LOCAL_PORT = 10055

type VncTab struct {
	TabItem        *container.TabItem
	app            fyne.App
	connection     *cap.Connection
	refresh_btn    *widget.Button
	new_btn        *widget.Button
	session_labels binding.StringList
	sessions       []cap.Session
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

func newVncTab(a fyne.App, conn *cap.Connection) *VncTab {
	t := VncTab{
		app:        a,
		connection: conn,
		// save: cb,
		session_labels: binding.BindStringList(&[]string{}),
	}

	sessions := widget.NewListWithData(t.session_labels,
		func() fyne.CanvasObject {
			hidden := widget.NewLabel("")
			hidden.Hide()
			connect_btn := widget.NewButton("Connect", func() {
				strs := strings.Split(hidden.Text, ",")
				owner_uid := strs[0]
				display := strs[1]

				otp := get_otp(conn, owner_uid, display)
				RunVnc(conn, otp, display)
			})
			label := widget.NewLabel("template")
			delete_btn := widget.NewButton("Kill", func() {
				KillSession(conn, hidden.Text, hidden.Text)
			})
			return container.NewHBox(hidden, connect_btn, label, delete_btn)
		},
		func(session binding.DataItem, obj fyne.CanvasObject) {
			box, ok := obj.(*fyne.Container)
			if ok {
				box.Objects[0].(*widget.Label).Bind(session.(binding.String))
				box.Objects[2].(*widget.Label).Bind(session.(binding.String))
			} else {
				log.Println("Warning: could not update VNC session list: ", box, session)
			}
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
			otp, displayNumber, err := t.connection.CreateVncSession(
				xres_entry.Text,
				yres_entry.Text,
			)
			if err == nil {
				t.refresh()
				go RunVnc(t.connection, otp, displayNumber)
			}
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

func get_otp(conn *cap.Connection, owner_uid, display string) string {
	// Use vncpasswd to generate OTP for sessions we own
	// log.Println("Checking %s ? %s", owner_uid, ssh.uid)
	if owner_uid == conn.GetUid() {
		return *get_owner_otp(conn, display)
	}

	return get_shared_otp(display)
}

func get_owner_otp(conn *cap.Connection, display string) *string {
	client := conn.GetClient()
	// loginName := conn.GetUsername()
	// if !display.startswith(loginName) {
	// panic("BAD DISPLAY")
	// }
	command := fmt.Sprintf("vncpasswd -o -display %s", display)
	log.Println("VNCPASSWD: ", command)
	output, err := client.CleanExec(command)
	if err != nil {
		log.Println("Error executing: ", command)
		return nil
	}
	for _, line := range strings.Split(output, "\n") {
		log.Println(line)
		prefix := "Full control one-time password:"
		if strings.HasPrefix(line, prefix) {
			otp := strings.TrimSpace(line[len(prefix):])
			return &otp
		}
	}
	log.Println("OTP not found from vncpasswd")
	return nil
}

func get_shared_otp(display string) string {
	// Use session manager to make OTP for shared sessions
	var response []string
	// response := ssh.sessionMessage("OTP", display).split()

	log.Println(response)

	nonce := strings.TrimSpace(response[0])
	encOTP := strings.TrimSpace(response[1])
	return decryptOTP([]byte(nonce), []byte(encOTP))
}

func MAC(message []byte) string {
	digest := cap.MakeSHADigest(
		// self._session_manager_secret[0:32],
		message,
	// self._session_manager_secret[33:64],
	)
	return hex.EncodeToString(digest[:])
}

func decryptOTP(nonce, encOTP []byte) string {

	key := MAC(nonce)
	decOTP := ""
	for ii, key_char := range key {
		encOTP_char := int64(encOTP[ii])
		k, _ := strconv.ParseInt(string(key_char), 16, 8)
		e := encOTP_char - 65
		o := (e - k) % 16
		decOTP += fmt.Sprint(o)
	}
	return decOTP
}

func KillSession(conn *cap.Connection, otp, displayNumber string) {
}
