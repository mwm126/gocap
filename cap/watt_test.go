//go:build integration
// +build integration

package cap

import (
	"net"
	"testing"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
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

func TestWattLoginButton(t *testing.T) {
	a := app.New()

	var conn_man FakeConnectionManager
	cfg := GetConfig()
	wattTab := NewJouleConnected(a, cfg, &conn_man)

	test.Type(wattTab.CapTab.usernameEntry, "the_user")
	test.Type(wattTab.CapTab.passwordEntry, "the_pass")
	wattTab.CapTab.networkSelect.SetSelected("vpn")

	test.Tap(wattTab.CapTab.loginBtn)

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, "the_user", conn_man.username)
	assert.Equal(t, net.IPv4(199, 249, 243, 253), conn_man.address)
}
