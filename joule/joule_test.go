package joule

import (
	"io"
	"net"
	"testing"
	"time"

	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/cap/sshtunnel"
	"aeolustec.com/capclient/login"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/ssh"
)

func NewFakeClient(server net.IP, user, pass string, port uint) (cap.Client, error) {
	client := FakeClient{}
	return &client, nil
}

type FakeClient struct{}

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

func (fsc *FakeClient) Start(command string) (io.ReadCloser, io.ReadCloser, error) {
	return nil, nil, nil
}

func (fsc *FakeClient) Wait() error {
	return nil
}

func TestJouleLoginButton(t *testing.T) {
	a := test.NewApp()
	var co fyne.CanvasObject
	w := test.NewWindow(co)

	knk := cap.NewKnocker(&StubYubikey{}, 0)
	conn_man := cap.NewCapConnectionManager(NewFakeClient, knk)

	var joule_service login.Service
	services, _ := login.FindServices()
	for _, service := range services {
		if service.Name == "joule" {
			joule_service = service
		}
	}

	jouleTab := NewJouleConnected(
		a,
		w,
		joule_service,
		conn_man,
		login.LoginInfo{Network: "alb_admin", Username: "the_user", Password: "the_pass"},
	)

	test.Tap(jouleTab.CapTab.ConnectForm.ConnectButton)

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
	//	want := net.IPv4(204, 154, 139, 11)
	//	got := conn.GetAddress()
	//	if diff := cmp.Diff(want, got); diff != "" {
	//		t.Errorf("Mismatch: %s", diff)
	//	}
	// })

	t.Run("Test Login", func(t *testing.T) {
		var client FakeClient
		fake_conn, err := cap.NewCapConnection(client, "", "")
		if err != nil {
			t.Error(err)
		}

		jouleTab.Connect(fake_conn)
	})
}
