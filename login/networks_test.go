package login

import (
	"testing"
)

func TestGetExternalIp(t *testing.T) {
	ip := GetExternalIp()

	if ip.String() == "127.0.0.1" {
		t.Error("Should get external IP address, not", ip)
	}
}
