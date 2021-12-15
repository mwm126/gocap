package joule

import (
	"aeolustec.com/capclient/cap"
	"fmt"
	"fyne.io/fyne/v2/test"
	"github.com/google/go-cmp/cmp"
	"testing"
)

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

func NewFakeKnocker() *cap.PortKnocker {
	var fake_yk StubYubikey
	return cap.NewPortKnocker(&fake_yk, 0)
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

func (c *FakeConnection) GetPassword() string {
	return "test_pwd"
}

func (conn *FakeConnection) UpdateForwards(fwds []string) {}

func (conn *FakeConnection) CreateVncSession(xres string, yres string) (string, string, error) {
	conn.sessions = append(conn.sessions, cap.Session{
		Username:      "test_user",
		DisplayNumber: ":77",
		Geometry:      fmt.Sprintf("%sx%s", xres, yres),
		DateCreated:   "2222-33-44",
		HostAddress:   "localhost",
		HostPort:      "8088",
	})
	return "", "", nil
}

func TestVncTab(t *testing.T) {
	a := test.NewApp()

	init_session := cap.Session{
		Username:      "the_user",
		DisplayNumber: ":123",
		Geometry:      "1661x888",
		DateCreated:   "2021-12-02",
		HostAddress:   "localhost",
		HostPort:      "789",
	}

	t.Run("Vnc Refresh Sessions", func(t *testing.T) {
		var conn FakeConnection
		conn.sessions = []cap.Session{init_session}
		vncTab := newVncTab(a, &conn)

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
		var conn FakeConnection
		conn.sessions = []cap.Session{init_session}
		vncTab := newVncTab(a, &conn)

		want := 0
		got := len(vncTab.sessions)
		if want != got {
			t.Error("Initially # of sessions should be 0 but was: ", len(vncTab.sessions))
		}

		test.Tap(vncTab.new_btn) // Should show new session dialog (without error)
	})

}

func TestNewSessionDialog(t *testing.T) {
	a := test.NewApp()
	init_session := cap.Session{
		Username:      "the_user",
		DisplayNumber: ":123",
		Geometry:      "1661x888",
		DateCreated:   "2021-12-02",
		HostAddress:   "localhost",
		HostPort:      "789",
	}
	default_rezs := []string{
		"800x600",
		"1024x768",
		"1280x1024",
		"1600x1200",
	}

	t.Run("Test Preset Resolution", func(t *testing.T) {
		var conn FakeConnection
		conn.sessions = []cap.Session{init_session}
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
		got := conn.sessions
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Mismatch: %s", diff)
		}
	})

	t.Run("Test Custom Resolution", func(t *testing.T) {
		var conn FakeConnection
		conn.sessions = []cap.Session{init_session}
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
		got := conn.sessions
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Mismatch: %s", diff)
		}
	})

}
