package cap

import (
	"log"
)

// Yubikey interface used by other code (can be real or faked)
type Yubikey interface {
	FindSerial() (int32, error)
	ChallengeResponse(chal [6]byte) ([16]byte, error)
	ChallengeResponseHMAC(chal SHADigest) ([20]byte, error)
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

func (yk *UsbYubikey) ChallengeResponse(chal [6]byte) ([16]byte, error) {
	out, err := run_yk_chalresp(chal)
	if err != nil {
		log.Println("Could not get challenge response", err)
		return [16]byte{}, err
	}
	return out, nil
}

func (yk *UsbYubikey) ChallengeResponseHMAC(chal SHADigest) ([20]byte, error) {
	var hmac [20]byte
	hmac, err := run_yk_hmac(chal)
	if err != nil {
		log.Println("Unable to get HMAC challenge response:", err)
	}
	return hmac, err
}
