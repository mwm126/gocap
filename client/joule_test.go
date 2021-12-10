//go:build integration
// +build integration

package client

import (
	"net"
	"testing"
	"time"

	"aeolustec.com/capclient/cap"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"github.com/google/go-cmp/cmp"
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

	t.Run("Test username entry", func(t *testing.T) {
		want := "the_user"
		got := conn_man.username
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Mismatch: %s", diff)
		}
	})

	t.Run("Test address selection", func(t *testing.T) {
		want := net.IPv4(204, 154, 139, 11)
		got := conn_man.address
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Mismatch: %s", diff)
		}
	})
}
