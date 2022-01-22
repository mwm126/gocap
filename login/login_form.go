package login

import (
	"aeolustec.com/capclient/cap"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// The LoginForm shows the widgets for logging in (MGN/PIT),(ADMIN/SCI),username,password
type LoginForm struct {
	Container     *fyne.Container
	NetworkSelect *widget.Select
	UsernameEntry *widget.Entry
	PasswordEntry *widget.Entry
	LoginButton   *widget.Button
}

// Returns pointer to a new LoginForm. Service defines what networks to list in
// the dropdown. When submitting the form, function connect_cb() is called with
// the entries for network, username, and password.
func NewLoginForm(
	service cap.Service,
	connect_cb func(LoginInfo),
	init_uname, init_pword string) *LoginForm {
	username := widget.NewEntry()
	username.SetText(init_uname)
	username.SetPlaceHolder("Enter username...")
	password := widget.NewPasswordEntry()
	password.SetText(init_pword)
	password.SetPlaceHolder("Enter password...")

	network_ips := make(map[string]string)
	external_ips := make(map[string]string)
	networkNames := make([]string, 0, len(service.Networks))
	for name, val := range service.Networks {
		network_ips[name] = val.CapServerAddress
		external_ips[name] = val.ClientExternalAddress
		networkNames = append(networkNames, name)
	}
	network := widget.NewSelect(networkNames, func(s string) {})

	network.SetSelected("external")
	login := widget.NewButton("Login", func() {
		login_info := LoginInfo{
			network.Selected, username.Text, password.Text,
		}
		go connect_cb(login_info)
	})
	return &LoginForm{
		Container:     container.NewVBox(username, password, network, login),
		NetworkSelect: network,
		UsernameEntry: username,
		PasswordEntry: password,
		LoginButton:   login,
	}
}

func (f *LoginForm) setEnabled(enabled bool) {
	if enabled {
		f.NetworkSelect.Enable()
		f.UsernameEntry.Enable()
		f.PasswordEntry.Enable()
		f.LoginButton.Enable()
	} else {
		f.NetworkSelect.Disable()
		f.UsernameEntry.Disable()
		f.PasswordEntry.Disable()
		f.LoginButton.Disable()
	}
}
