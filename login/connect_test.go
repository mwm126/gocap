package login

import (
	"testing"

	"aeolustec.com/capclient/cap"
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

func TestNewCapTab(t *testing.T) {
	knk := cap.NewKnocker(&StubYubikey{}, 0)
	conn_man := cap.NewCapConnectionManager(NewFakeClient, knk)

	tab := NewCapTab(
		"some_tab",
		"some new tab for some service",
		cap.Service{},
		conn_man,
		func(c *cap.Connection) {},
		nil,
		LoginInfo{Network: "alb_admin", Username: "the_user", Password: "the_pass"},
	)

	tab.CloseConnection()
}
