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
	server := net.IPv4(11, 22, 33, 44)

	var fake_yk FakeYubikey
	var entropy [32]byte
	var fake_kckr Knocker
	fake_kckr = NewPortKnocker(&fake_yk, entropy)
	conn_man := NewCapConnectionManager(fake_kckr)
	conn, _ := conn_man.newCapConnection(username, password, server)
	defer conn.close()

	assert.Equal(t, conn.connectionInfo.username, "testusername")
	assert.Equal(t, conn.connectionInfo.password, "testpassword")
}
