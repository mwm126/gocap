package login

import (
	"testing"

	"aeolustec.com/capclient/cap"
)

func TestFormEnableDisable(t *testing.T) {
	form := NewConnectForm(
		cap.Service{},
		LoginInfo{Network: "alb_admin", Username: "the_user", Password: "the_pass"},
		func(i LoginInfo) {},
	)

	form.setEnabled(true)
	if form.ConnectButton.Disabled() {
		t.Errorf("Enable Login Form failed")
	}

	form.setEnabled(false)
	if !form.ConnectButton.Disabled() {
		t.Errorf("Disable Login Form failed")
	}
}
