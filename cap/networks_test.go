package cap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCapConnection(t *testing.T) {
	networks := getNetworks()

	assert.Equal(t, Network{
		true,
		(""),
		("204.154.140.51"),
		("204.154.139.11"),
		62201,
		("172.16.0.1"),
		("204.154.140.10"),
	}, networks["external"])
}

func TestGetExternalIp(t *testing.T) {
	ip := GetExternalIp()

	assert.NotEqual(t, ip.String(), "127.0.0.1")
}
