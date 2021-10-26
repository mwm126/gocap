package cap

import (
	"encoding/hex"
	"log"
)

// Yubikey interface used by other code (can be real or faked)
type Yubikey interface {
	FindSerial() (int32, error)
	challengeResponse(chal [6]byte) ([16]byte, error)
	challengeResponseHMAC(chal SHADigest) ([20]byte, error)
}

// UsbYubikey implementation (for actual yubikey)
type UsbYubikey struct{}

func (yk *UsbYubikey) FindSerial() (int32, error) {
	out, err := run_yk_info()
	if err != nil {
		return 0, error(err)
	}
	return out, err
}

func (yk *UsbYubikey) challengeResponse(chal [6]byte) ([16]byte, error) {
	challengeArgument := hex.EncodeToString(chal[:])
	out, err := run_yk_chalresp(challengeArgument)
	if err != nil {
		log.Println("Could not get challenge response", err)
		return [16]byte{}, err
	}
	return out, nil
}

func (yk *UsbYubikey) challengeResponseHMAC(chal SHADigest) ([20]byte, error) {
	var hmac [20]byte
	hex_chal := hex.EncodeToString(chal[:])
	hmac, err := run_yk_hmac(hex_chal)
	if err != nil {
		log.Println("Unable to get HMAC challenge response:", err)
	}
	return hmac, err
}
