package joule

import (
	"errors"
	"io"
	"log"
	"net"
	"testing"

	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/cap/sshtunnel"
	"fyne.io/fyne/v2/test"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/ssh"
)

func NewFakeVncClient(server net.IP, user, pass string) (cap.Client, error) {
	client := FakeVncClient{}
	return &client, nil
}

type FakeVncClient struct{}

func (fsc FakeVncClient) CleanExec(command string) (string, error) {
	replies := map[string]string{
		"hostname": "the_hostname",
		`ping -c 1 the_hostname| grep PING|awk '{print $3}'| sed "s/(//"|sed "s/)//"`: "1.2.3.4",
		`id|sed "s/uid=//"|sed "s/(/ /"|awk '{print $1}'`:                             "the_uid",
		"ps auxnww|grep Xvnc|grep -v grep":                                            `8227 27248  0.0  0.0  92620 70572 ?        S    Aug03   0:234 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :123 -desktop TurboVNC: login03:5 (the_user) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5905 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1`,
	}
	reply, exists := replies[command]
	if !exists {
		log.Println(command, " : ", reply)
		return "", errors.New("Unexpected command")
	}
	return replies[command], nil
}

func (fsc FakeVncClient) Close() {
}

func (client FakeVncClient) CheckPasswordExpired(
	pass string,
	pw_expired_cb func(cap.Client),
	ch chan string,
) error {
	return nil
}

func (sc FakeVncClient) OpenSSHTunnel(
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

func (fsc *FakeVncClient) Start(command string) (io.ReadCloser, io.ReadCloser, error) {
	return nil, nil, nil
}

func (fsc *FakeVncClient) Wait() error {
	return nil
}

type StubYubikey struct{}

func (yk *StubYubikey) FindSerial() (int32, error) {
	return 0, nil
}

func (yk *StubYubikey) ChallengeResponse(chal [6]byte) ([16]byte, error) {
	return [16]byte{}, nil
}

func (yk *StubYubikey) ChallengeResponseHMAC(chal cap.SHADigest) ([20]byte, error) {
	return [20]byte{}, nil
}

func TestVncTab(t *testing.T) {
	a := test.NewApp()

	t.Run("Vnc Refresh Sessions", func(t *testing.T) {
		knk := cap.NewKnocker(&StubYubikey{}, 0)
		conn_man := cap.NewCapConnectionManager(NewFakeVncClient, knk)
		ext := net.IPv4(1, 1, 1, 1)
		srv := net.IPv4(1, 1, 1, 1)
		conn, err := conn_man.Connect(
			"the_user",
			"pass",
			ext,
			srv,
			123,
			func(cap.Client) {},
			make(chan string),
		)
		if err != nil {
			t.Error(err)
		}
		vncTab := newVncTab(a, conn)

		want := 0
		got := len(vncTab.sessions)
		if want != got {
			t.Error("Initially # of sessions should be 0 but was: ", len(vncTab.sessions))
		}

		test.Tap(vncTab.refresh_btn)

		want = 1
		got = len(vncTab.sessions)
		if want != got {
			t.Error("After refresh # of sessions should be 1 but was: ", len(vncTab.sessions))
		}
	})

	t.Run("Vnc New Session", func(t *testing.T) {
		var conn cap.Connection
		// conn.sessions = []cap.Session{init_session}
		vncTab := newVncTab(a, &conn)

		want := 0
		got := len(vncTab.sessions)
		if want != got {
			t.Error("Initially # of sessions should be 0 but was: ", len(vncTab.sessions))
		}

		test.Tap(vncTab.new_btn) // Should show new session dialog (without error)
	})

}

func _TestNewSessionDialog(t *testing.T) {
	a := test.NewApp()

	// init_session := cap.Session{
	//	Username:      "the_user",
	//	DisplayNumber: ":123",
	//	Geometry:      "1661x888",
	//	DateCreated:   "2021-12-02",
	//	HostAddress:   "localhost",
	//	HostPort:      "789",
	// }

	default_rezs := []string{
		"800x600",
		"1024x768",
		"1280x1024",
		"1600x1200",
	}

	t.Run("Test Preset Resolution", func(t *testing.T) {
		var conn cap.Connection
		// conn.sessions = []cap.Session{init_session}
		vncTab := newVncTab(a, &conn)
		w := test.NewWindow(nil)
		vsf := vncTab.NewVncSessionForm(w, default_rezs)
		last_index := len(vsf.preset_select.Options) - 1
		vsf.preset_select.SetSelectedIndex(last_index)
		vsf.Form.OnSubmit()

		want := []cap.Session{
			{
				Username:      "the_user",
				DisplayNumber: ":123",
				Geometry:      "1661x888",
				DateCreated:   "2021-12-02",
				HostAddress:   "localhost",
				HostPort:      "789",
			},
			{
				Username:      "test_user",
				DisplayNumber: ":77",
				Geometry:      "1600x1200",
				DateCreated:   "2222-33-44",
				HostAddress:   "localhost",
				HostPort:      "8088",
			}}
		// got := conn.sessions
		got := make([]cap.Session, 2)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Mismatch: %s", diff)
		}
	})

	t.Run("Test Custom Resolution", func(t *testing.T) {
		var conn cap.Connection
		// conn.sessions = []cap.Session{init_session}
		vncTab := newVncTab(a, &conn)
		w := test.NewWindow(nil)
		vsf := vncTab.NewVncSessionForm(w, default_rezs)
		vsf.xres_entry.SetText("999")
		vsf.yres_entry.SetText("555")
		vsf.Form.OnSubmit()

		want := []cap.Session{
			{
				Username:      "the_user",
				DisplayNumber: ":123",
				Geometry:      "1661x888",
				DateCreated:   "2021-12-02",
				HostAddress:   "localhost",
				HostPort:      "789",
			},
			{
				Username:      "test_user",
				DisplayNumber: ":77",
				Geometry:      "999x555",
				DateCreated:   "2222-33-44",
				HostAddress:   "localhost",
				HostPort:      "8088",
			}}
		// got := conn.sessions
		got := make([]cap.Session, 2)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Mismatch: %s", diff)
		}
	})

}

func TestVncCmd(t *testing.T) {
	cmd := VncCmd("abc", "xyz")

	want := "echo abc | env -u LD_LIBRARY_PATH vncviewer_HPCEE -highqual -autopass 127.0.0.1::10055 &"
	got := cmd
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Mismatch: %s", diff)
	}
}
