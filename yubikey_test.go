// +build yubikey

package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestYubikeySerial(t *testing.T) {

	serial_num := find_serial()

	assert.Equal(t, 5417533, serial_num)
}

func TestYubikeyChalResp(t *testing.T) {

	resp := challenge_response("123456")

	assert.Equal(t, 32, len(resp))
}

func TestYubikeyChalRespHMAC(t *testing.T) {

	resp := challenge_response_hmac("123456")

	assert.Equal(t, 20, len(resp))
}
