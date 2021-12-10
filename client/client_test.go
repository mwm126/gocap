package client

import (
	"net"
	"testing"

	"aeolustec.com/capclient/cap"
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
			var services []Service
			var cfg config
			if tc.fe261 {
				services = append(services, Service{Name: "fe261"})
			}
			if tc.joule {
				services = append(services, Service{Name: "joule"})
			}
			if tc.watt {
				services = append(services, Service{Name: "watt"})
			}
			InitServices(&services)

			a := test.NewApp()
			w := test.NewWindow(nil)
			client := NewClient(a, w, cfg, conn_man)

			// test.Tap(client.LoginTab.loginBtn)
			client.LoginTab.ConnectedCallback(&FakeConnection{})

			if got := len(client.Tabs.Items); got != tc.ntabs {
				t.Errorf("Got %d; want %d", got, tc.ntabs)
			}
		})
	}
}

type FakeConnection struct {
	sessions []cap.Session
}

func (c *FakeConnection) FindSessions() ([]cap.Session, error) {
	return c.sessions, nil
}

func (c *FakeConnection) GetUsername() string {
	return "test_user"
}

func (conn *FakeConnection) UpdateForwards(fwds []string) {}

func (conn *FakeConnection) CreateVncSession(xres string, yres string) (string, string, error) {
	conn.sessions = append(conn.sessions, cap.Session{
		Username:      "test_user",
		DisplayNumber: ":77",
		Geometry:      "geo",
		DateCreated:   "2222-33-44",
		HostAddress:   "localhost",
		HostPort:      "8088",
	})
	return "", "", nil
}
