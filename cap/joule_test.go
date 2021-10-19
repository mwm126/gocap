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

type JouleSpyKnocker struct {
	username string
	address  net.IP
	knocked  bool
}

func (sk *JouleSpyKnocker) Knock(username string, address net.IP) error {
	sk.knocked = true
	sk.username = username
	sk.address = address
	return nil
}

func TestJouleLoginButton(t *testing.T) {
	spy := &JouleSpyKnocker{}
	a := app.New()

	var fake_yk FakeYubikey
	var entropy [32]byte
	fake_kckr := NewPortKnocker(&fake_yk, entropy)
	conn_man := NewCapConnectionManager(fake_kckr)
	cfg := GetConfig()
	jouleTab := NewCapTab("Joule", "NETL SuperComputer", cfg.Joule_Ips, conn_man,
		NewJouleConnected(a, conn_man, func() {}))

	test.Type(jouleTab.usernameEntry, "the_user")
	test.Type(jouleTab.passwordEntry, "the_pass")
	jouleTab.networkSelect.SetSelected("alb_admin")

	test.Tap(jouleTab.loginBtn)

	time.Sleep(100 * time.Millisecond)
	assert.True(t, spy.knocked)
	assert.Equal(t, "the_user", spy.username)
	assert.Equal(t, net.IPv4(204, 154, 139, 11), spy.address)
}
