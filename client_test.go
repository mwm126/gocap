package main

import (
	"net"
	"testing"

	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/cap/sshtunnel"
	"aeolustec.com/capclient/config"
	"aeolustec.com/capclient/login"
	"fyne.io/fyne/v2/test"
	"golang.org/x/crypto/ssh"
)

type FakeClient struct{}

func NewFakeClient(server net.IP, user, pass string) (cap.Client, error) {
	client := FakeClient{}
	return &client, nil
}

func (fsc FakeClient) CleanExec(command string) (string, error) {
	return "", nil
}

func (fsc FakeClient) Close() {
}

func (client FakeClient) CheckPasswordExpired(
	pass string,
	pw_expired_cb func(cap.Client),
	ch chan string,
) error {
	return nil
}

func (sc FakeClient) OpenSSHTunnel(
	user, pass string,
	local_port int,
	remote_addr string,
	remote_port int,
) sshtunnel.SSHTunnel {
	return *sshtunnel.NewSSHTunnel(
		nil,
		"testuser@localhost",
		ssh.Password(pass),
		"rem_addr:123",
		"123",
	)
}

func (sc FakeClient) Dial(protocol, endpoint string) (net.Conn, error) {
	return nil, nil
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
			knk := cap.NewKnocker(&FakeYubikey{}, 0)
			conn_man := cap.NewCapConnectionManager(NewFakeClient, knk)
			var services []login.Service
			if tc.fe261 {
				services = append(services, login.Service{Name: "fe261"})
			}
			if tc.joule {
				services = append(services, login.Service{Name: "joule"})
			}
			if tc.watt {
				services = append(services, login.Service{Name: "watt"})
			}
			if err := login.InitServices(&services); err != nil {
				t.Fatal(err)
			}
			a := test.NewApp()
			w := test.NewWindow(nil)
			var cfg config.Config
			client := NewClient(a, w, cfg, conn_man)

			test.Tap(client.LoginTab.LoginForm.LoginButton)
			client.setupServices(login.LoginInfo{Network: "", Username: "", Password: ""}, services)

			if got := len(client.Tabs.Items); got != tc.ntabs {
				t.Errorf("Got %d; want %d", got, tc.ntabs)
			}
		})
	}
}

type FakeYubikey struct {
	Available bool
}

func (yk *FakeYubikey) YubikeyAvailable() bool {
	return true
}

func (yk *FakeYubikey) FindSerial() (int32, error) {
	return 5417533, nil
}

func (yk *FakeYubikey) ChallengeResponse(chal [6]byte) ([16]byte, error) {
	var resp [16]byte
	return resp, nil
}

func (yk *FakeYubikey) ChallengeResponseHMAC(chal cap.SHADigest) ([20]byte, error) {
	var hmac [20]byte
	return hmac, nil
}
