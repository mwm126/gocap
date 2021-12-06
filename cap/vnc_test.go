//go:build integration
// +build integration

package cap

import (
	"aeolustec.com/capclient/cap/connection"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"testing"
)

type StubYubikey struct{}

func (yk *StubYubikey) FindSerial() (int32, error) {
	return 0, nil
}

func (yk *StubYubikey) ChallengeResponse(chal [6]byte) ([16]byte, error) {
	return [16]byte{}, nil
}

func (yk *StubYubikey) ChallengeResponseHMAC(chal connection.SHADigest) ([20]byte, error) {
	return [20]byte{}, nil
}

func NewFakeKnocker() *connection.PortKnocker {
	var fake_yk StubYubikey
	var entropy [32]byte
	return connection.NewPortKnocker(&fake_yk, entropy)
}

type FakeConnection struct{}

func (c *FakeConnection) FindSessions() ([]connection.Session, error) {
	sessions := make([]connection.Session, 0, 10)
	return sessions, nil
}

func (c *FakeConnection) GetUsername() string {
	return ""
}

func (conn *FakeConnection) UpdateForwards(fwds []string) {}

func DiableTestVncTabRefresh(t *testing.T) {
	a := app.New()
	w := a.NewWindow("Hello")
	var conn FakeConnection
	vncTab := newVncTab(a, &conn)
	tabItem := newVncTabItem(vncTab)
	tabs := container.NewAppTabs(tabItem)
	w.SetContent(container.NewVBox(
		tabs,
	))

	test.Tap(vncTab.refresh_btn)

	want := 0
	got := len(vncTab.sessions)
	if want != got {
		t.Error("Could not refresh sessions")
	}
}

func TestVncTabNewSession(t *testing.T) {
	a := app.New()
	w := a.NewWindow("Hello")
	var conn FakeConnection
	vncTab := newVncTab(a, &conn)
	tabItem := newVncTabItem(vncTab)
	tabs := container.NewAppTabs(tabItem)
	w.SetContent(container.NewVBox(
		tabs,
	))

	test.Tap(vncTab.new_btn)

	// if 0 != len(vncTab.sessions) {
	// 	t.Error("Number of sessions should be 1 but was: ", len(vncTab.sessions))
	// }
	// want := connection.Session{}
	// got := vncTab.sessions[0]
	// if want != got {
	// 	t.Error("BAD THING u DID")
	// }

}
