package client

import (
	"testing"

	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/config"
	"aeolustec.com/capclient/login"
	"fyne.io/fyne/v2/test"
)

func TestWattClient(t *testing.T) {
	yk := new(cap.UsbYubikey)
	knk := cap.NewKnocker(yk, 0)
	conn_man := cap.NewCapConnectionManager(NewFakeClient, knk)
	var services []login.Service
	services = append(services, login.Service{Name: "watt"})
	login.SetDemoServices(services)
	a := test.NewApp()
	w := test.NewWindow(nil)
	var cfg config.Config
	client := NewClient(a, w, cfg, conn_man)

	test.Tap(client.LoginTab.LoginForm.LoginButton)
	client.setupServices(
		&login.LoginInfo{Network: "", Username: "", Password: ""},
		services,
	)

	if got := len(client.Tabs.Items); got != 2 {
		t.Errorf("Got %d; want %d", got, 2)
	}
}
