//go:build yubikey
// +build yubikey

package connection

import (
	"testing"
)

func TestYubikeySerial(t *testing.T) {
	yk := new(UsbYubikey)
	serial_num, err := yk.FindSerial()

	if err != nil {
		t.Error("Error with HMAC Challenge Response:", err)
	}
	want := int32(5417533)
	got := serial_num
	if want != got {
		t.Errorf("response should be %d bytes long, but was %d", want, got)
	}
}

func TestYubikeyChalResp(t *testing.T) {
	yk := new(UsbYubikey)
	resp, err := yk.ChallengeResponse([6]byte{0, 1, 2, 3, 4, 5})

	if err != nil {
		t.Error("Error with HMAC Challenge Response:", err)
	}
	want := 16
	got := len(resp)
	if want != got {
		t.Errorf("response should be %d bytes long, but was %d", want, got)
	}
}

func TestYubikeyChalRespHMAC(t *testing.T) {
	yk := new(UsbYubikey)

	resp, err := yk.ChallengeResponseHMAC([32]byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5,
	})

	if err != nil {
		t.Error("Error with HMAC Challenge Response:", err)
	}
	want := 20
	got := len(resp)
	if want != got {
		t.Errorf("response should be %d bytes long, but was %d", want, got)
	}
}
