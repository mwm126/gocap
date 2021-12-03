//go:build integration
// +build integration

package cap

import (
	"net"
	"testing"
	"time"

	"aeolustec.com/capclient/cap/connection"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

type FakeConnectionManager struct {
	username string
	address  net.IP
}

func (t *FakeConnectionManager) GetConnection() connection.Connection {
	return nil
}

func (t *FakeConnectionManager) Close() {
}

func (cm *FakeConnectionManager) Connect(
	user, pass string,
	ext_addr,
	server net.IP,
	port uint,
	pw_expired_cb func(connection.PasswordChecker),
	ch chan string) error {
	cm.username = user
	cm.address = server
	return nil
}

func (c *FakeConnectionManager) GetPasswordExpired() bool {
	return false
}
func (c *FakeConnectionManager) SetPasswordExpired() {}

func TestJouleLoginButton(t *testing.T) {
	a := app.New()

	var conn_man FakeConnectionManager
	cfg := GetConfig()

	jouleTab := NewJouleConnected(a, cfg, &conn_man)

	test.Type(jouleTab.CapTab.usernameEntry, "the_user")
	test.Type(jouleTab.CapTab.passwordEntry, "the_pass")
	jouleTab.CapTab.networkSelect.SetSelected("alb_admin")

	test.Tap(jouleTab.CapTab.loginBtn)

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, "the_user", conn_man.username)
	assert.Equal(t, net.IPv4(204, 154, 139, 11), conn_man.address)
}
