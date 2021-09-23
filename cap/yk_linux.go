package cap

// #cgo CFLAGS: -g -Wall -Werror -I/usr/include/ykpers-1/
// #cgo LDFLAGS: /usr/lib/x86_64-linux-gnu/libykpers-1.a /usr/lib/x86_64-linux-gnu/libyubikey.a -lusb-1.0
// #include <stdlib.h>
// #include "yk.h"
import "C"

import (
	"encoding/hex"
	"errors"
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

	ch := (*C.uchar)(&challenge[0])
	ohtipi := (*C.uchar)(&otp[0])

	C.get_otp(ch, ohtipi)

	return otp, nil
}

func run_yk_hmac(chal string) ([20]byte, error) {
	var hmac [20]byte

	challenge, err := hex.DecodeString(chal)
	if err != nil {
		return hmac, err
	}
	digest := (*C.uchar)(&challenge[0])
	hmac_c := (*C.uchar)(&hmac[0])

	C.hmac_from_digest(digest, hmac_c)

	return hmac, nil
}
