//go:build integration
// +build integration

package client

import (
	"net"
	"testing"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"github.com/google/go-cmp/cmp"
)

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

func _TestWattLoginButton(t *testing.T) {
	a := app.New()

	var conn_man FakeConnectionManager
	cfg := GetConfig()
	wattTab := NewJouleConnected(a, cfg, &conn_man)

	test.Type(wattTab.CapTab.usernameEntry, "the_user")
	wattTab.CapTab.networkSelect.SetSelected("vpn")

	test.Tap(wattTab.CapTab.loginBtn)

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
