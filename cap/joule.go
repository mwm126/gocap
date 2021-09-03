package cap

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type JouleTab struct {
	Tab        *container.TabItem
	connection *CapConnection
}

func NewJouleTab(knocker Knocker, a fyne.App) JouleTab {
	joule := &JouleTab{}
	var card *widget.Card
	var jouleLogin, jouleConnecting, jouleConnected *fyne.Container
	connect_cancelled := false

	jouleLogin = NewJouleLogin(func(user, pass, host string) {
		card.SetContent(jouleConnecting)
		conn, err := newCapConnection(user, pass, host, knocker)

		if err != nil {
			log.Println("Unable to make CAP Connection")
			card.SetContent(jouleLogin)
			connect_cancelled = false
			return
		}

		if connect_cancelled {
			log.Println("CAP Connection cancelled.")
			conn.close()
			connect_cancelled = false
			return
		}

		joule.connection = conn
		time.Sleep(1 * time.Second)
		card.SetContent(jouleConnected)
	})

	jouleConnecting = NewJouleConnecting(func() {
		connect_cancelled = true
		card.SetContent(jouleLogin)
	})

	jouleConnected = NewJouleConnected(a, joule, func() {
		joule.connection.close()
		joule.connection = nil
		card.SetContent(jouleLogin)
	})

	card = widget.NewCard("Joule 2.0", "NETL Supercomputer", jouleLogin)

	joule.Tab = container.NewTabItem("Joule", card)
	return *joule
}

func NewJouleLogin(connect_cb func(user, pass, host string)) *fyne.Container {
	username := widget.NewEntry()
	username.SetPlaceHolder("Enter username...")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Enter password...")
	network := widget.NewSelect(getNetworkNames(), func(s string) {})
	networks := getNetworks()
	network.SetSelected("external")
	login := widget.NewButton("Login", func() {
		go connect_cb(username.Text, password.Text, networks[network.Selected].JouleAddress)
	})
	return container.NewVBox(username, password, network, login)
}

func NewJouleConnecting(cancel_cb func()) *fyne.Container {
	connecting := widget.NewLabel("Connecting......")
	cancel := widget.NewButton("Cancel", func() {
		cancel_cb()
	})
	return container.NewVBox(connecting, cancel)
}

func NewJouleConnected(app fyne.App, joule *JouleTab, close_cb func()) *fyne.Container {

	homeTab := newJouleHome(close_cb)
	sshTab := newJouleSsh()

	vcard := widget.NewCard("GUI", "TODO", nil)
	vncTab := container.NewTabItem("VNC", vcard)

	fwdTab := newJouleFwds(app)

	tabs := container.NewAppTabs(
		homeTab,
		sshTab,
		vncTab,
		fwdTab,
	)
	return container.NewMax(tabs)
}

func newJouleHome(close_cb func()) *container.TabItem {
	close := widget.NewButton("Disconnect", func() {
		close_cb()
	})
	box := container.NewVBox(widget.NewLabel("Connected!"), close)
	return container.NewTabItem("Home", box)
}

func newJouleSsh() *container.TabItem {
	ssh := widget.NewButton("New SSH Session", func() {
		cmd := exec.Command("gnome-terminal", "--", "ssh", "localhost", "-p", strconv.Itoa(SSH_LOCAL_PORT))
		err := cmd.Run()
		if err != nil {
			log.Println("gnome-terminal FAIL: ", err)
		}
	})
	label := widget.NewLabel(fmt.Sprintf("or run in a terminal:  ssh localhost -p %d", SSH_LOCAL_PORT))
	box := container.NewVBox(widget.NewLabel("To create new Terminal (SSH) Session in gnome-terminal:"), ssh, label)
	return container.NewTabItem("SSH", box)
}

func newJouleFwds(app fyne.App) *container.TabItem {
	forwards := binding.BindStringList(
		&[]string{
			"10022:localhost:22",
			"20022:localhost:33",
		},
	)
	var to_be_removed widget.ListItemID

	add := widget.NewButton("Add", func() { addJouleFwd(app, forwards) })
	var remove *widget.Button
	remove = widget.NewButton("Remove", func() {
		removeForward(forwards, to_be_removed)
		remove.Disable()
	})
	remove.Disable()

	list := widget.NewListWithData(forwards,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(fwd binding.DataItem, obj fyne.CanvasObject) {
			obj.(*widget.Label).Bind(fwd.(binding.String))
		})
	list.OnUnselected = func(id widget.ListItemID) { remove.Disable() }
	list.OnSelected = func(id widget.ListItemID) {
		if id < 2 {
			remove.Disable()
			log.Println("Cannot remove fixed forward #", id)
		} else {
			remove.Enable()
			to_be_removed = id
		}
	}

	box := container.NewBorder(add, remove, nil, nil, list)
	return container.NewTabItem("Port Fowards", box)
}

func addJouleFwd(app fyne.App, forwards binding.StringList) {
	win := app.NewWindow("Add Port Forward")

	local_p := widget.NewEntry()
	local_p.SetPlaceHolder("Local Port")
	remote_h := widget.NewEntry()
	remote_h.SetPlaceHolder("Remote Host")
	remote_p := widget.NewEntry()
	remote_p.SetPlaceHolder("Remote Port")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Local Port", Widget: local_p},
			{Text: "Remote Host", Widget: remote_h},
			{Text: "Remote Port", Widget: remote_p},
		},
		OnSubmit: func() {
			new_fwd := fmt.Sprintf("%s:%s:%s", local_p.Text, remote_h.Text, remote_p.Text)
			log.Println("Adding forward:", new_fwd)
			forwards.Append(new_fwd)
			win.Close()
		},
		OnCancel:   func() { win.Close() },
		SubmitText: "Ok",
		CancelText: "Cancel",
	}
	win.SetContent(form)
	win.Show()
}

func removeForward(forwards binding.StringList, to_be_removed int) {
	fwds, _ := forwards.Get()
	for i := range fwds {
		if i == to_be_removed {
			fwds = append(fwds[:i], fwds[i+1:]...)
			break
		}
	}
	forwards.Set(fwds)
}
