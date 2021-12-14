package login

import (
	"aeolustec.com/capclient/cap"
	"fmt"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
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

type FakeConnection struct {
	sessions []cap.Session
}

func (c *FakeConnection) FindSessions() ([]cap.Session, error) {
	return c.sessions, nil
}

func (c *FakeConnection) GetUsername() string {
	return "test_user"
}

func (c *FakeConnection) GetPassword() string {
	return "test_pwd"
}

func (conn *FakeConnection) UpdateForwards(fwds []string) {}

func (conn *FakeConnection) CreateVncSession(xres string, yres string) (string, string, error) {
	conn.sessions = append(conn.sessions, cap.Session{
		Username:      "test_user",
		DisplayNumber: ":77",
		Geometry:      fmt.Sprintf("%sx%s", xres, yres),
		DateCreated:   "2222-33-44",
		HostAddress:   "localhost",
		HostPort:      "8088",
	})
	return "", "", nil
}

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
		func(conn cap.Connection, login_info LoginInfo) {
			ct := NewCapTab("test tab", "for testing", Service{},
				conn_man, func(cap cap.Connection) {},
				connctd, login_info)
			tabs.Append(ct.Tab)
			w.SetContent(tabs)
		}, connctd, "", "")
	w = test.NewWindow(tabs)

	conn := &FakeConnection{}
	login_tab.ConnectedCallback(conn, login_info)

	login_tab.CloseConnection()

}
