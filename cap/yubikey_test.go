// +build yubikey

package cap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYubikeySerial(t *testing.T) {

	serial_num := findSerial()

	assert.Equal(t, 5417533, serial_num)
}

func TestYubikeyChalResp(t *testing.T) {

	resp := challengeResponse("123456")

	assert.Equal(t, 32, len(resp))
}

func TestYubikeyChalRespHMAC(t *testing.T) {

	resp := challengeResponseHMAC("123456")

	assert.Equal(t, 20, len(resp))
}
