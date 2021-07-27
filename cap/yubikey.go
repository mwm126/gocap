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
	findSerial() int
	// challengeResponse(chal [32]byte) [32]byte
	challengeResponses(chal [32]byte) [32]byte
	challengeResponseHMAC(chal [32]byte) [32]byte
}

// UsbYubikey implementation (for actual yubikey)
type UsbYubikey struct{}

func (yk *UsbYubikey) findSerial() int {
	return 5417533

	out, err := exec.Command("ykinfo", "-s", "-q").Output()
	if err != nil {
		log.Fatal(err)
	}
	serial := strings.TrimSpace(string(out))
	i, err := strconv.Atoi(serial)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func (yk *UsbYubikey) challengeResponses(chal [32]byte) [32]byte {
	log.Println("lenggggggggggggggth chal ", len(chal), chal)
	var qwer [32]byte
	return qwer
	out, err := exec.Command("ykchalresp", "-1", "-Y", "-x", string(chal[:])).Output()
	if err != nil {
		log.Fatal(err)
	}
	response := strings.TrimSpace(string(out))
	log.Println(response)
	if len(response) != 32 {
		log.Fatal("Invalid ykchalresp")
	}
	resp := modhexDecode(response)
	log.Println("lenggggggggggggggth resp ", len(resp))
	log.Println("lenggggggggggggggth resp ", resp)
	if err != nil {
		log.Fatal(resp, err)
	}
	return resp
}

func modhexDecode(m string) [32]byte {
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
	var h [32]byte
	for i, c := range m {
		h[i] = mod2hex[c]
	}
	return h

}

func (yk *UsbYubikey) challengeResponseHMAC(chal [32]byte) [32]byte {

	var qwer [32]byte
	return qwer

	out, err := exec.Command("ykchalresp", "-2", "-H", "-x", string(chal[:])).Output()
	if err != nil {
		log.Fatal(err)
	}
	response := strings.TrimSpace(string(out))
	log.Println(response)
	if len(response) != 40 {
		log.Fatal("Invalid ykchalresp")
	}
	resp, err := hex.DecodeString(response)
	if err != nil {
		log.Fatal(resp, err)
	}

	var hmac [32]byte
	copy(hmac[:], resp)
	return hmac
}
