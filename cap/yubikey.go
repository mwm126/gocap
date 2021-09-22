package cap

import (
	"encoding/hex"
	"log"
	"strings"
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
	responseStr := strings.TrimSpace(string(out))
	response := modhexDecode(responseStr[:16])
	return response, nil
}

func modhexDecode(m string) [16]byte {
	mod2hex := map[rune]byte{
		'c': 0x0,
		'b': 0x1,
		'd': 0x2,
		'e': 0x3,
		'f': 0x4,
		'g': 0x5,
		'h': 0x6,
		'i': 0x7,
		'j': 0x8,
		'k': 0x9,
		'l': 0xa,
		'n': 0xb,
		'r': 0xc,
		't': 0xd,
		'u': 0xe,
		'v': 0xf,
	}
	var h [16]byte
	for i, c := range m {
		h[i] = mod2hex[c]
	}
	return h

}

func (yk *UsbYubikey) challengeResponseHMAC(chal SHADigest) ([20]byte, error) {
	var hmac [20]byte
	hex_chal := hex.EncodeToString(chal[:])
	out, err := run_yk_hmac(hex_chal)
	if err != nil {
		log.Println("Unable to run ykchalresp:", err)
		return hmac, err
	}
	responseHex := strings.TrimSpace(out)
	response, err := hex.DecodeString(responseHex)
	if err != nil {
		log.Fatal(response, err)
	}

	copy(hmac[:], response[:20])
	return hmac, nil
}
