package cap

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
)

type WattSpyKnocker struct {
	username string
	password string
	network  string
	knocked  bool
}

func (sk *WattSpyKnocker) Knock(username, password, network string) {
	sk.knocked = true
	sk.username = username
	sk.password = password
	sk.network = network
}

func TestWattLoginButton(t *testing.T) {
	spy := &WattSpyKnocker{}
	a := app.New()
	wattTab := NewWattTab(spy, a)
	test.Type(wattTab.UsernameEntry, "the_user")
	test.Type(wattTab.PasswordEntry, "the_pass")
	wattTab.NetworkSelect.SetSelected("external")

	test.Tap(wattTab.LoginBtn)

	time.Sleep(100 * time.Millisecond)
	assert.True(t, spy.knocked)
	assert.Equal(t, "the_user", spy.username)
	assert.Equal(t, "the_pass", spy.password)
	assert.Equal(t, "204.154.140.10", spy.network)
}
