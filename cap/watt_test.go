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

func (sk *WattSpyKnocker) Knock(username string, address net.IP) error {
	sk.knocked = true
	sk.username = username
	sk.address = address
	return nil
}

func TestWattLoginButton(t *testing.T) {
	spy := &WattSpyKnocker{}
	a := app.New()

	var fake_yk FakeYubikey
	var entropy [32]byte
	fake_kckr := NewPortKnocker(&fake_yk, entropy)
	conn_man := NewCapConnectionManager(fake_kckr)
	cfg := GetConfig()
	wattTab := NewCapTab("Watt", "NETL SuperComputer", cfg.Watt_Ips, conn_man,
		NewWattConnected(a, conn_man, func() {}))

	test.Type(wattTab.usernameEntry, "the_user")
	test.Type(wattTab.passwordEntry, "the_pass")
	wattTab.networkSelect.SetSelected("vpn")

	test.Tap(wattTab.loginBtn)

	time.Sleep(100 * time.Millisecond)
	assert.True(t, spy.knocked)
	assert.Equal(t, "the_user", spy.username)
	assert.Equal(t, net.IPv4(199, 249, 243, 253), spy.address)
}
