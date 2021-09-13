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
	password string
	address  net.IP
	knocked  bool
}

func (sk *JouleSpyKnocker) Knock(username, password string, address net.IP) {
	sk.knocked = true
	sk.username = username
	sk.password = password
	sk.address = address
}

func TestJouleLoginButton(t *testing.T) {
	spy := &JouleSpyKnocker{}
	a := app.New()
	jouleTab := NewJouleTab(spy, a)
	test.Type(jouleTab.usernameEntry, "the_user")
	test.Type(jouleTab.passwordEntry, "the_pass")
	jouleTab.networkSelect.SetSelected("alb_admin")

	test.Tap(jouleTab.loginBtn)

	time.Sleep(100 * time.Millisecond)
	assert.True(t, spy.knocked)
	assert.Equal(t, "the_user", spy.username)
	assert.Equal(t, "the_pass", spy.password)
	assert.Equal(t, net.IPv4(204, 154, 139, 11), spy.address)
}
