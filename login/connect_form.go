package login

import (
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// The LoginForm shows the widgets for logging in (MGN/PIT),(ADMIN/SCI),username,password
type ConnectForm struct {
	Container     *fyne.Container
	ConnectButton *widget.Button
}

// Returns pointer to a new LoginForm. Service defines what networks to list in
// the dropdown. When submitting the form, function connect_cb() is called with
// the entries for network, username, and password.
func NewConnectForm(
	service Service,
	login_info LoginInfo,
	connect_cb func(LoginInfo)) *ConnectForm {

	connect := widget.NewButton("Connect", func() {
		go connect_cb(login_info)
	})
	return &ConnectForm{
		Container:     container.NewVBox(connect),
		ConnectButton: connect,
	}
}

func (f *ConnectForm) setEnabled(enabled bool) {
	if enabled {
		f.ConnectButton.Enable()
	} else {
		f.ConnectButton.Disable()
	}
}
