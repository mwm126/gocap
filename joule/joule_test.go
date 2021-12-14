package joule

import (
	"net"
	"testing"
	"time"

	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/login"
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
	a := test.NewApp()

	var conn_man FakeConnectionManager

	login.InitServices(nil)
	var joule_service login.Service
	services, _ := login.FindServices()
	for _, service := range services {
		if service.Name == "joule" {
			joule_service = service
		}
	}

	jouleTab := NewJouleConnected(
		a,
		joule_service,
		&conn_man,
		login.LoginInfo{"alb_admin", "the_user", "the_pass"},
	)

	test.Tap(jouleTab.CapTab.ConnectBtn)

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

	t.Run("Test Login", func(t *testing.T) {
		fake_conn := &FakeConnection{}
		jouleTab.Connect(fake_conn)
	})
}
