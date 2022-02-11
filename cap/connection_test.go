package cap

import (
	"net"
	"testing"

	"aeolustec.com/capclient/cap/sshtunnel"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/ssh"
)

func NewFakeClient(server net.IP, user, pass string, port uint) (Client, error) {
	client := FakeClient{}
	return &client, nil
}

type CmdResult struct {
	Out string
	Err error
}
type FakeClient struct {
	ActivatedShell []string
	Outputs        map[string]CmdResult
}

func (fsc FakeClient) CleanExec(command string) (string, error) {
	return "", nil
}

func (fsc FakeClient) Close() {
}

func (client FakeClient) CheckPasswordExpired(
	pass string,
	pw_expired_cb func(Client),
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

type StubYubikey struct{}

func (yk *StubYubikey) FindSerial() (int32, error) {
	return 0, nil
}

func (yk *StubYubikey) ChallengeResponse(chal [6]byte) ([16]byte, error) {
	return [16]byte{}, nil
}

func (yk *StubYubikey) ChallengeResponseHMAC(chal SHADigest) ([20]byte, error) {
	return [20]byte{}, nil
}

func TestCapConnection(t *testing.T) {
	username := "testusername"
	password := "testpassword"
	ext_ip := net.IPv4(11, 22, 33, 44)
	server := net.IPv4(55, 66, 77, 88)

	knk := NewKnocker(&StubYubikey{}, 0)
	conn_man := NewCapConnectionManager(NewFakeClient, knk)
	ch := make(chan string)
	conn, err := conn_man.Connect(
		username,
		password,
		ext_ip,
		server,
		123,
		22,
		func(pwc Client) {},
		ch,
	)
	if err != nil {
		t.Error("Error making cap connection:", err)
	}

	t.Run("Test connection username", func(t *testing.T) {
		want := "testusername"
		got := conn.GetUsername()
		if want != got {
			t.Errorf("Did not set connection username: want %s but got %s", want, got)
		}
	})
	t.Run("Test connection password", func(t *testing.T) {
		want := "testpassword"
		got := conn.GetPassword()
		if want != got {
			t.Errorf("Did not set connection password: want %s but got %s", want, got)
		}
	})

	t.Run("Test UpdateForwards", func(t *testing.T) {
		conn.UpdateForwards([]string{"123,4.5.6.7,890"})
	})
}

func TestParseVncProcesses(t *testing.T) {
	sessions := parseSessions(
		"mmeredith",
		`8227 27248  0.0  0.0  92620 70572 ?        S    Aug03   0:20 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :5 -desktop TurboVNC: login03:5 (not_mark) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5905 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
8227 27248  0.0  0.0  92620 70572 ?        S    Aug03   0:20 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :5 -desktop TurboVNC: login03:5 (mmeredith) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5905 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1
8227 27248  0.0  0.0  92620 70572 ?        S    Aug03   0:20 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :5 -desktop TurboVNC: login03:5 (meredithm) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5905 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1`,
	)

	want := Session{
		"mmeredith",
		":5",
		"3840x2160",
		"Aug03",
		5905,
	}
	got := sessions[0]
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Mismatch: %s", diff)
	}
}
