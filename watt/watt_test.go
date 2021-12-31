package watt

import (
	"net"
	"testing"
	"time"

	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/cap/sshtunnel"
	"aeolustec.com/capclient/login"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/test"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/ssh"
)

func NewFakeClient(server net.IP, user, pass, port string) (cap.Client, error) {
	client := FakeClient{}
	return &client, nil
}

type FakeClient struct {
	ActivatedShell []string
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

type WattSpyKnocker struct {
	username string
	address  net.IP
	knocked  bool
}

func (sk *WattSpyKnocker) Knock(username string, address net.IP, port uint) error {
	sk.knocked = true
	sk.username = username
	sk.address = address
	return nil
}

type FakeYubikey struct{}

func (yk *FakeYubikey) FindSerial() (int32, error) {
	return 0, nil
}

func (yk *FakeYubikey) ChallengeResponse(chal [6]byte) ([16]byte, error) {
	return [16]byte{}, nil
}

func (yk *FakeYubikey) ChallengeResponseHMAC(chal cap.SHADigest) ([20]byte, error) {
	return [20]byte{}, nil
}

func TestWattLoginButton(t *testing.T) {
	a := app.New()

	knk := cap.NewKnocker(&FakeYubikey{}, 0)
	conn_man := cap.NewCapConnectionManager(NewFakeClient, knk)
	err := login.InitServices(nil)
	if err != nil {
		t.Fatal(err)
	}
	var watt_service login.Service
	services, _ := login.FindServices()
	for _, service := range services {
		if service.Name == "watt" {
			watt_service = service
		}
	}
	wattTab := NewWattConnected(
		a,
		watt_service,
		conn_man,
		login.LoginInfo{Network: "vpn", Username: "the_user", Password: ""},
	)

	test.Tap(wattTab.CapTab.ConnectForm.ConnectButton)

	time.Sleep(100 * time.Millisecond)

	t.Run("Test username entry", func(t *testing.T) {
		var client FakeClient
		conn, err := cap.NewCapConnection(client, "the_user", "the_pass")
		if err != nil {
			t.Error(err)
		}

		want := "the_user"
		got := conn.GetUsername()
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Mismatch: %s", diff)
		}
	})

	// t.Run("Test address selection", func(t *testing.T) {
	//	var client FakeClient
	//	conn, err := cap.NewCapConnection(client, "the_user", "the_pass")
	//	if err != nil {
	//		t.Error(err)
	//	}

	//	want := net.IPv4(199, 249, 243, 253)
	//	got := conn.GetAddress()
	//	if diff := cmp.Diff(want, got); diff != "" {
	//		t.Errorf("Mismatch: %s", diff)
	//	}
	// })

	t.Run("Test Login", func(t *testing.T) {
		var fake_conn cap.Connection
		wattTab.Connect(&fake_conn)
	})
}
