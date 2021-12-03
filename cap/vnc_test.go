//go:build integration
// +build integration

package cap

import (
	"aeolustec.com/capclient/cap/connection"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
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

func TestVncTab(t *testing.T) {
	a := app.New()
	w := a.NewWindow("Hello")

	var conn FakeConnection

	vncTab := newVncTab(&conn)
	tabItem := newVncTabItem(vncTab)
	tabs := container.NewAppTabs(tabItem)

	w.SetContent(container.NewVBox(
		tabs,
	))

	test.Tap(vncTab.refresh_btn)

	expected := make([]connection.Session, 0)
	assert.Equal(t, vncTab.sessions, expected)
}
