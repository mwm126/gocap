package cap

import (
	"aeolustec.com/capclient/cap/connection"
	"net"
	"testing"
)

type FakeConnManager struct {
	username string
	address  net.IP
}

func (t *FakeConnManager) GetConnection() connection.Connection {
	return nil
}
func (cm *FakeConnManager) Connect(
	user, pass string,
	ext_addr,
	server net.IP,
	port uint,
	pw_expired_cb func(connection.PasswordChecker),
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
		{"none", false, false, false, 1},
		{"w", false, false, true, 2},
		{"j", false, true, false, 2},
		{"f", true, false, false, 2},
		{"fj", true, true, false, 3},
		{"fw", true, false, true, 3},
		{"jw", false, true, true, 3},
		{"fjw", true, true, true, 4},
	}
	for _, tc := range testCases {
		t.Run(tc.label, func(t *testing.T) {
			conn_man := &FakeConnManager{"the_user", net.IPv4(1, 1, 1, 1)}
			cfg := config{
				Enable_fe261: tc.fe261,
				Enable_joule: tc.joule,
				Enable_watt:  tc.watt,
			}

			client := NewClient(cfg, conn_man)

			if got := len(client.Tabs.Items); got != tc.ntabs {
				t.Errorf("Got %d; want %d", got, tc.ntabs)
			}
		})
	}
}
