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
	password string
	address  net.IP
	knocked  bool
}

func (sk *WattSpyKnocker) Knock(username, password string, address net.IP) {
	sk.knocked = true
	sk.username = username
	sk.password = password
	sk.address = address
}

func TestWattLoginButton(t *testing.T) {
	spy := &WattSpyKnocker{}
	a := app.New()
	wattTab := NewWattTab(spy, a)
	test.Type(wattTab.UsernameEntry, "the_user")
	test.Type(wattTab.PasswordEntry, "the_pass")
	wattTab.NetworkSelect.SetSelected("vpn")

	test.Tap(wattTab.LoginBtn)

	time.Sleep(100 * time.Millisecond)
	assert.True(t, spy.knocked)
	assert.Equal(t, "the_user", spy.username)
	assert.Equal(t, "the_pass", spy.password)
	assert.Equal(t, net.IPv4(199, 249, 243, 253), spy.address)
}
