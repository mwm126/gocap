package main

import (
	"encoding/hex"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func find_serial() int {
	out, err := exec.Command("ykinfo", "-s", "-q").Output()
	if err != nil {
		log.Fatal(err)
	}
	out_s := strings.TrimSpace(string(out))
	i, err := strconv.Atoi(out_s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func challenge_response(chal string) [32]byte {
	out, err := exec.Command("ykchalresp", "-1", "-Y", "-x", chal).Output()
	if err != nil {
		log.Fatal(err)
	}
	out_s := strings.TrimSpace(string(out))
	log.Println(out_s)
	if len(out_s) != 32 {
		log.Fatal("Invalid ykchalresp")
	}
	resp := modhexDecode(out_s)
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

func challenge_response_hmac(chal string) []byte {
	out, err := exec.Command("ykchalresp", "-2", "-H", "-x", chal).Output()
	if err != nil {
		log.Fatal(err)
	}
	out_s := strings.TrimSpace(string(out))
	log.Println(out_s)
	if len(out_s) != 40 {
		log.Fatal("Invalid ykchalresp")
	}
	resp, err := hex.DecodeString(out_s)
	if err != nil {
		log.Fatal(resp, err)
	}
	return resp
}
