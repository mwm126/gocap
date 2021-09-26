//go:build integration
// +build integration

package cap

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

type StubYubikey struct{}

func (yk *StubYubikey) FindSerial() (int32, error) {
	return 0, nil
}

func (yk *StubYubikey) challengeResponse(chal [6]byte) ([16]byte, error) {
	return [16]byte{}, nil
}

func (yk *StubYubikey) challengeResponseHMAC(chal SHADigest) ([20]byte, error) {
	return [20]byte{}, nil
}

func TestCapConnection(t *testing.T) {
	username := "testusername"
	password := "testpassword"
	ext_ip := net.IPv4(11, 22, 33, 44)
	server := net.IPv4(55, 66, 77, 88)

	var fake_yk StubYubikey
	var entropy [32]byte
	var fake_kckr Knocker
	fake_kckr = NewPortKnocker(&fake_yk, entropy)
	conn_man := NewCapConnectionManager(fake_kckr)
	err := conn_man.Connect(username, password, ext_ip, server)
	if err != nil {
		assert.FailNow(t, "failed to make cap connection", err)
	}

	assert.NotNil(t, conn_man.connection.connectionInfo)
	assert.Equal(t, conn_man.connection.connectionInfo.username, "testusername")
	assert.Equal(t, conn_man.connection.connectionInfo.password, "testpassword")
}
