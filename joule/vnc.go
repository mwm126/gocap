package joule

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"aeolustec.com/capclient/cap"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type VncRunner interface {
	RunVnc(otp, display string, port int)
}

type ExeRunner struct{}

func (r *ExeRunner) RunVnc(otp, display string, port int) {
	RunVnc(otp, display, port)
}

type ItemButton struct {
	widget.Button
	session cap.Session
}

type ItemList struct {
	items []cap.Session
}

func (i *ItemList) Set(items []cap.Session) {
	i.items = items
}

func (i *ItemList) AddListener(listener binding.DataListener) {
}

func (i *ItemList) RemoveListener(listener binding.DataListener) {
}

func (i *ItemList) GetItem(index int) (binding.DataItem, error) {
	var s cap.Session
	if index > i.Length()-1 {
		return &s, errors.New(fmt.Sprintf("Invalid index %d > max index %d", index, i.Length()-1))
	}
	return i.items[index], nil
}

func (i *ItemList) Length() int {
	return len(i.items)
}

type PortFinder interface {
	FindPort() (int, error)
}

type FreePortFinder struct{}

func (fpf FreePortFinder) FindPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

type VncTab struct {
	TabItem     *container.TabItem
	List        *widget.List
	VncRunner   VncRunner
	list_items  map[cap.Session]*fyne.Container
	app         fyne.App
	connection  *cap.Connection
	refresh_btn *widget.Button
	new_btn     *widget.Button
	sessions    *ItemList
	PortFinder  PortFinder
}

func (vt *VncTab) refresh() error {
	sessions, err := vt.connection.FindSessions()
	if err != nil {
		log.Println("Warning - unable to refresh: ", err)
	}
	items := make([]cap.Session, 0)
	for _, session := range sessions {
		items = append(items, session)
	}
	vt.sessions.Set(items)
	vt.List.Refresh()
	return err
}

func newVncTab(a fyne.App, conn *cap.Connection, vnc_runner VncRunner, pf PortFinder) *VncTab {
	t := VncTab{
		list_items: make(map[cap.Session]*fyne.Container, 0),
		VncRunner:  vnc_runner,
		app:        a,
		connection: conn,
		sessions:   &ItemList{},
		PortFinder: pf,
	}

	t.List = widget.NewListWithData(t.sessions,
		func() fyne.CanvasObject {
			label := widget.NewLabel("placeholder")

			connect_btn := &ItemButton{}
			connect_btn.ExtendBaseWidget(connect_btn)
			connect_btn.Text = "Connect"
			connect_btn.OnTapped = func() {
				local_p, err := t.PortFinder.FindPort()
				if err != nil {
					log.Println("Could not find free port for VNC session: ", err)
					return
				}
				remote_h := conn.GetAddress()
				remote_p := connect_btn.session.HostPort
				tunnel, err := conn.NewTunnel(local_p, remote_h, remote_p)
				if err != nil {
					log.Println("Could not forward VNC port: ", err)
					return
				}
				display_num := connect_btn.session.DisplayNumber
				display := fmt.Sprintf("%s%s", remote_h, display_num)
				otp := get_otp(conn, conn.GetUid(), display)
				if otp == nil {
					log.Println("WARNING: OTP is missing, cannot connect.")
					return
				}
				connect_btn.Disable()
				t.VncRunner.RunVnc(*otp, display, local_p)
				tunnel.Close()
				connect_btn.Enable()
			}
			delete_btn := widget.NewButton("Kill", func() {
				KillSession(conn, connect_btn.session.Username, connect_btn.session.DisplayNumber)
			})
			return container.NewHBox(connect_btn, label, delete_btn)
		},
		func(item binding.DataItem, obj fyne.CanvasObject) {
			box, ok := obj.(*fyne.Container)
			if !ok {
				log.Println("Warning: could not update VNC session list: ", box, item)
				return
			}
			session := item.(cap.Session)
			box.Objects[0].(*ItemButton).session = session
			box.Objects[1].(*widget.Label).SetText(session.Label())
			t.list_items[session] = box
		})

	t.refresh_btn = widget.NewButton("Refresh", func() { t.refresh() })
	t.new_btn = widget.NewButton("New VNC Session", t.showNewVncSessionDialog)
	vcard := widget.NewCard("GUI", "List of VNC Sessions", t.new_btn)
	box := container.NewBorder(vcard, t.refresh_btn, nil, nil, t.List)

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
			defer win.Close()
			local_p, err := t.PortFinder.FindPort()
			if err != nil {
				log.Println("Could not find free port for vnc client: ", err)
				return
			}
			otp, displayNumber, err := t.connection.CreateVncSession(
				xres_entry.Text,
				yres_entry.Text,
			)
			if err != nil {
				log.Println("Could not create VNC session: ", err)
				return
			}
			t.refresh()
			go RunVnc(otp, displayNumber, local_p)
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

func get_otp(conn *cap.Connection, owner_uid, display string) *string {
	// Use vncpasswd to generate OTP for sessions we own
	// log.Println("Checking %s ? %s", owner_uid, ssh.uid)
	if owner_uid == conn.GetUid() {
		log.Println("owner is uid", owner_uid)
		return get_owner_otp(conn, display)
	}

	log.Println("owner NOT uid", owner_uid, conn.GetUid())
	return get_shared_otp(display)
}

func get_owner_otp(conn *cap.Connection, display string) *string {
	client := conn.GetClient()
	address := conn.GetAddress()
	if !strings.HasPrefix(display, address) {
		log.Println("Display does not match loginName:", display, address)
		return nil
	}
	command := fmt.Sprintf("vncpasswd -o -display %s", display)
	log.Println("VNCPASSWD: ", command)
	output, err := client.CleanExec(command)
	if err != nil {
		log.Println("Error: ", err)
		log.Println("Output: ", output)
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

func get_shared_otp(display string) *string {
	// Use session manager to make OTP for shared sessions
	response := []string{"abc", "123"}
	// response := ssh.sessionMessage("OTP", display).split()

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

func decryptOTP(nonce, encOTP []byte) *string {
	key := MAC(nonce)
	decOTP := ""
	for ii, key_char := range key {
		encOTP_char := int64(encOTP[ii])
		k, _ := strconv.ParseInt(string(key_char), 16, 8)
		e := encOTP_char - 65
		o := (e - k) % 16
		decOTP += fmt.Sprint(o)
	}
	return &decOTP
}

func KillSession(conn *cap.Connection, otp, displayNumber string) {
}
