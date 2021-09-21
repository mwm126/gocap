//go:build integration
// +build integration

package cap

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestCapConnection(t *testing.T) {
	username := "testusername"
	password := "testpassword"
	ext_ip := net.IPv4(11, 22, 33, 44)
	server := net.IPv4(55, 66, 77, 88)

	var fake_yk FakeYubikey
	var entropy [32]byte
	var fake_kckr Knocker
	fake_kckr = NewPortKnocker(&fake_yk, entropy)
	conn_man := NewCapConnectionManager(fake_kckr)
	conn, _ := conn_man.newCapConnection(username, password, ext_ip, server)
	assert.NotNil(t, conn)
	defer conn.close()

	assert.NotNil(t, conn.connectionInfo)
	assert.Equal(t, conn.connectionInfo.username, "testusername")
	assert.Equal(t, conn.connectionInfo.password, "testpassword")
}
