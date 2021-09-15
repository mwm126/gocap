//go:build integration
// +build integration

package cap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCapConnection(t *testing.T) {
	username := "testusername"
	password := "testpassword"
	server := "localhost:22"

	knckr := &FakeYubikey{}
	conn := newCapConnection(username, password, server, knckr)
	defer conn.close()

	guis := conn.listGUIs()

	assert.Equal(t, "/etc/passwd\n", guis)
}
