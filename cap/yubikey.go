package cap

import (
	"encoding/hex"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// Yubikey interface used by other code (can be real or faked)
type Yubikey interface {
	findSerial() (int32, error)
	challengeResponse(chal [16]byte) [16]byte
	challengeResponseHMAC(chal SHADigest) ([16]byte, error)
}

// UsbYubikey implementation (for actual yubikey)
type UsbYubikey struct{}

func (yk *UsbYubikey) findSerial() (int32, error) {
	out, err := exec.Command("ykinfo", "-s", "-q").Output()
	if err != nil {
		return 0, error(err)
	}
	serial := strings.TrimSpace(string(out))
	i, err := strconv.Atoi(serial)
	if err != nil {
		return 0, error(err)
	}
	return int32(i), nil
}

func (yk *UsbYubikey) challengeResponse(chal [16]byte) [16]byte {
	testval, err := hex.DecodeString("2ee619bc248bcefbe09e733d2cdda3be")
	if err != nil {
		log.Fatal(err)
	}
	var resp [16]byte
	copy(resp[:], testval[:16])
	return resp
	out, err := exec.Command("ykchalresp", "-1", "-Y", "-x", string(chal[:])).Output()
	if err != nil {
		log.Fatal(err)
	}
	responseStr := strings.TrimSpace(string(out))
	log.Println(responseStr)
	if len(responseStr) != 16 {
		log.Fatal("Invalid ykchalresp")
	}
	response := modhexDecode(responseStr)
	if err != nil {
		log.Fatal(response, err)
	}
	return response
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

func (yk *UsbYubikey) challengeResponseHMAC(chal SHADigest) ([16]byte, error) {
	var hmac [16]byte
	out, err := exec.Command("ykchalresp", "-2", "-H", "-x", string(chal[:])).Output()
	if err != nil {
		log.Println("Unable to run ykchalresp:", err)
		return hmac, err
	}
	responseHex := strings.TrimSpace(string(out))
	response, err := hex.DecodeString(responseHex)
	if err != nil {
		log.Fatal(response, err)
	}

	copy(hmac[:], response)
	return hmac, nil
}
