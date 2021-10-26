package cap

// #cgo CFLAGS: -g -Wall -Werror -I/usr/local/include/ykpers-1 -I/usr/local/include
// #cgo LDFLAGS: /usr/local/lib/libykpers-1.a /usr/local/lib/libyubikey.a -framework CoreServices -framework IOKit
// #include <stdlib.h>
// #include "yk.h"
import "C"

import (
	"encoding/hex"
	"errors"
	"log"
	"os/exec"
	"strings"
)

func run_yk_info() (int32, error) {
	serial := C.get_yk_serial()
	if serial < 0 {
		return -1, errors.New("Error getting info from Yubikey")
	}
	return int32(serial), nil
}

func run_yk_chalresp(chal string) ([16]byte, error) {
	var otp [16]byte
	challenge, err := hex.DecodeString(chal)
	if err != nil {
		return [16]byte{}, err
	}

	cmd := exec.Command(
		"ykchalresp",
		"-1",
		"-Y",
		"-x",
		chal)
	out, err := cmd.Output()
	if err != nil {
		log.Println("ykchalresp FAIL: ", err)
		return otp, err
	}
	out_s := strings.TrimSpace(string(out[:]))
	if len(out_s) != 32 {
		return otp, errors.New("invalid challenge response")
	}
	otp, err = modhexDecode(out_s)
	return otp, err

	ch := (*C.uchar)(&challenge[0])
	ohtipi := (*C.uchar)(&otp[0])

	code := C.get_otp(ch, ohtipi)
	if code != 0 {
		err = errors.New("Error from get_otp")
	}

	return otp, err
}

func run_yk_hmac(chal string) ([20]byte, error) {
	var hmac [20]byte

	challenge, err := hex.DecodeString(chal)
	if err != nil {
		return hmac, err
	}

	cmd := exec.Command(
		"ykchalresp",
		"-2",
		"-H",
		"-x",
		chal)
	out, err := cmd.Output()
	if err != nil {
		log.Println("ykchalresp HMAC FAIL: ", err)
		return hmac, err
	}
	out_s := strings.TrimSpace(string(out[:]))
	hmac_s, err := hex.DecodeString(out_s)
	if err != nil {
		return hmac, err
	}
	hmac = *(*[20]byte)(hmac_s)
	return hmac, err

	digest := (*C.uchar)(&challenge[0])
	hmac_c := (*C.uchar)(&hmac[0])

	code := C.hmac_from_digest(digest, hmac_c)
	if code != 0 {
		err = errors.New("Error from hmac_from_digest")
	}

	return hmac, err
}

func modhexDecode(m string) ([16]byte, error) {
	from_mod := map[rune]rune{
		'c': '0',
		'b': '1',
		'd': '2',
		'e': '3',
		'f': '4',
		'g': '5',
		'h': '6',
		'i': '7',
		'j': '8',
		'k': '9',
		'l': 'a',
		'n': 'b',
		'r': 'c',
		't': 'd',
		'u': 'e',
		'v': 'f',
	}
	var hexstring [32]rune
	if len(m) != 32 {
		var result [16]byte
		return result, errors.New("invalid challenge")
	}
	for ii, value := range m {
		hexstring[ii] = from_mod[value]
	}
	result, err := hex.DecodeString(string(hexstring[:]))
	return *(*[16]byte)(result), err
}
