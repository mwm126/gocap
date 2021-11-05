package cap

// #cgo CFLAGS: -g -Wall -Werror -I/usr/include/ykpers-1/ -I/usr/share/mingw-w64/include/
// #cgo LDFLAGS: -L/usr/x86_64-w64-mingw32/lib -Wl,-Bstatic -lykpers-1 -lyubikey -Wl,-Bdynamic -lsetupapi -lhid
// #include <stdlib.h>
// #include "yk.h"
import "C"

import (
	"encoding/hex"
	"errors"
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
	var buffer [1000]C.char

	argv := [5](*C.char){
		C.CString(""),
		C.CString("-1"),
		C.CString("-Y"),
		C.CString("-x"),
		C.CString(chal)}
	code := C.the_main(5, &argv[0], &buffer[0])
	if code != 0 {
		err := errors.New("Error from get_otp")
		return otp, err
	}
	out_s := strings.TrimSpace(C.GoString(&buffer[0]))
	if len(out_s) != 32 {
		return otp, errors.New("invalid challenge response")
	}
	otp, err := modhexDecode(out_s)
	return otp, err
}

func run_yk_hmac(chal string) ([20]byte, error) {
	var hmac [20]byte
	var buffer [1000]C.char
	argv := [5](*C.char){
		C.CString(""),
		C.CString("-2"),
		C.CString("-H"),
		C.CString("-x"),
		C.CString(chal),
	}
	code := C.the_main(5, &argv[0], &buffer[0])
	if code != 0 {
		err := errors.New("Error from get_otp")
		return hmac, err
	}
	out_s := strings.TrimSpace(C.GoString(&buffer[0]))
	if len(out_s) != 40 {
		return hmac, errors.New("Bad response from yubikey: " + out_s)
	}
	hmac_s, err := hex.DecodeString(out_s)
	copy(hmac[:], hmac_s)
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
	var result [16]byte
	if len(m) != 32 {
		return result, errors.New("invalid challenge")
	}
	for ii, value := range m {
		hexstring[ii] = from_mod[value]
	}
	result_s, err := hex.DecodeString(string(hexstring[:]))
	copy(result[:], result_s)
	return result, err
}
