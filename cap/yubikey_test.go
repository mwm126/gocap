//go:build yubikey
// +build yubikey

package cap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYubikeySerial(t *testing.T) {
	yk := new(UsbYubikey)
	serial_num, _ := yk.FindSerial()

	assert.Equal(t, int32(5417533), serial_num)
}

func TestYubikeyChalResp(t *testing.T) {
	yk := new(UsbYubikey)
	resp, _ := yk.challengeResponse([6]byte{0, 1, 2, 3, 4, 5})

	assert.Equal(t, 16, len(resp))
}

func TestYubikeyChalRespHMAC(t *testing.T) {
	yk := new(UsbYubikey)
	resp, _ := yk.challengeResponseHMAC([32]byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5,
	})

	assert.Equal(t, 20, len(resp))
}
