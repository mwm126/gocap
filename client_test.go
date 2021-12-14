package main

import (
	"net"
	"testing"

	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/config"
	"aeolustec.com/capclient/login"
	"fyne.io/fyne/v2/test"
)

type FakeConnManager struct {
	username string
	address  net.IP
}

func (t *FakeConnManager) GetConnection() cap.Connection {
	return nil
}
func (cm *FakeConnManager) Connect(
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
func (t *FakeConnManager) Close()              {}
func (c *FakeConnManager) SetPasswordExpired() {}
func (c *FakeConnManager) GetPasswordExpired() bool {
	return false
}

func TestClient(t *testing.T) {
	testCases := []struct {
		label string
		fe261 bool
		joule bool
		watt  bool
		ntabs int
	}{
		{"none", false, false, false, 2},
		{"w", false, false, true, 3},
		{"j", false, true, false, 3},
		{"f", true, false, false, 3},
		{"fj", true, true, false, 4},
		{"fw", true, false, true, 4},
		{"jw", false, true, true, 4},
		{"fjw", true, true, true, 5},
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			conn_man := &FakeConnManager{"the_user", net.IPv4(1, 1, 1, 1)}
			var services []login.Service
			var cfg config.Config
			if tc.fe261 {
				services = append(services, login.Service{Name: "fe261"})
			}
			if tc.joule {
				services = append(services, login.Service{Name: "joule"})
			}
			if tc.watt {
				services = append(services, login.Service{Name: "watt"})
			}
			login.InitServices(&services)

			a := test.NewApp()
			w := test.NewWindow(nil)
			client := NewClient(a, w, cfg, conn_man)

			// test.Tap(client.LoginTab.LoginBtn)
			client.LoginTab.ConnectedCallback(login.LoginInfo{"", "", ""})

			if got := len(client.Tabs.Items); got != tc.ntabs {
				t.Errorf("Got %d; want %d", got, tc.ntabs)
			}
		})
	}
}
