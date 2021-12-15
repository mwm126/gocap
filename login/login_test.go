package login

import (
	"aeolustec.com/capclient/cap"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

type FakeConnectionManager struct {
	username string
	address  net.IP
}

func (t *FakeConnectionManager) GetConnection() cap.Connection {
	return nil
}

func (t *FakeConnectionManager) Close() {
}

func (cm *FakeConnectionManager) Connect(
	user, pass string,
	ext_addr,
	server net.IP,
	port uint,
	pw_expired_cb func(cap.PasswordChecker),
	ch chan string) error {
	cm.username = user
	cm.address = server
	return nil
}

func (c *FakeConnectionManager) GetPasswordExpired() bool {
	return false
}
func (c *FakeConnectionManager) SetPasswordExpired() {}

func TestLoginTab(t *testing.T) {
	a := test.NewApp()
	a.Run()
	var w fyne.Window
	var login_info LoginInfo
	connctd := container.NewVBox(widget.NewLabel("Connected!"))
	service := Service{}
	conn_man := &FakeConnectionManager{}
	tabs := container.NewAppTabs()
	login_tab := NewLoginTab("Login", "NETL SuperComputer", service, conn_man,
		func(login_info LoginInfo) {
			ct := NewCapTab("test tab", "for testing", Service{},
				conn_man, func(cap cap.Connection) {},
				connctd, login_info)
			tabs.Append(ct.Tab)
			w.SetContent(tabs)
		}, connctd, "", "")
	w = test.NewWindow(tabs)

	t.Run("Enabled", func(t *testing.T) {
		assert.False(t, login_tab.NetworkSelect.Disabled())
		assert.False(t, login_tab.UsernameEntry.Disabled())
		assert.False(t, login_tab.PasswordEntry.Disabled())
		assert.False(t, login_tab.LoginBtn.Disabled())
	})

	t.Run("Disabled", func(t *testing.T) {
		login_tab.Disable()
		assert.True(t, login_tab.NetworkSelect.Disabled())
		assert.True(t, login_tab.UsernameEntry.Disabled())
		assert.True(t, login_tab.PasswordEntry.Disabled())
		assert.True(t, login_tab.LoginBtn.Disabled())

		login_tab.Enable()
		assert.False(t, login_tab.NetworkSelect.Disabled())
		assert.False(t, login_tab.UsernameEntry.Disabled())
		assert.False(t, login_tab.PasswordEntry.Disabled())
		assert.False(t, login_tab.LoginBtn.Disabled())
	})

	login_tab.ConnectedCallback(login_info)

	login_tab.CloseConnection()

}
