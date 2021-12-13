//go:build integration
// +build integration

package watt

import (
	"net"
	"testing"
	"time"

	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/login"
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

type WattSpyKnocker struct {
	username string
	address  net.IP
	knocked  bool
}

func (sk *WattSpyKnocker) Knock(username string, address net.IP, port uint) error {
	sk.knocked = true
	sk.username = username
	sk.address = address
	return nil
}

func TestWattLoginButton(t *testing.T) {
	a := app.New()

	var conn_man FakeConnectionManager
	login.InitServices(nil)
	var watt_service login.Service
	services, _ := login.FindServices()
	for _, service := range services {
		if service.Name == "watt" {
			watt_service = service
		}
	}
	wattTab := NewWattConnected(a, watt_service, &conn_man)

	test.Type(wattTab.CapTab.UsernameEntry, "the_user")
	wattTab.CapTab.NetworkSelect.SetSelected("vpn")

	test.Tap(wattTab.CapTab.LoginBtn)

	time.Sleep(100 * time.Millisecond)

	t.Run("Test username entry", func(t *testing.T) {
		want := "the_user"
		got := conn_man.username
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Mismatch: %s", diff)
		}
	})

	t.Run("Test address selection", func(t *testing.T) {
		want := net.IPv4(199, 249, 243, 253)
		got := conn_man.address
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Mismatch: %s", diff)
		}
	})
}
