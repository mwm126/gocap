package cap

// #cgo CFLAGS: -g -Wall -I/usr/local/include/ykpers-1 -I/usr/local/include
// #cgo LDFLAGS: /usr/local/lib/libykpers-1.a /usr/local/lib/libyubikey.a -framework CoreServices -framework IOKit
// #include "yk_darwin.h"
import (
	"C"
)

import (
	_ "embed"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"os/exec"
	"path"
)

//go:generate go run gen.go
//go:generate unzip -o -d embeds embeds/ykpers-1.20.0-win64.zip

func run_yk_info() (int32, error) {
	serial := C.get_yk_serial()
	if serial < 0 {
		return -1, errors.New("Error getting info from Yubikey")
	}
	return int32(serial), nil
}

func run_yk_chalresp(chal string) ([]byte, error) {
	dir, err := os.MkdirTemp("", "capclient")
	defer os.RemoveAll(dir)
	if err != nil {
		log.Println("Could not make temporary directory")
		return []byte{}, err
	}

	ykc := path.Join(dir, "ykchalresp")
	yki := path.Join(dir, "ykinfo")
	save(ykc, []byte("TODO"))
	save(yki, []byte("TODO"))

	// log.Println(ykc, "-1", "-Y", "-x", chal)

	cmd := exec.Command(ykc, "-1", "-Y", "-x", chal)
	output, err := cmd.Output()
	return output, err
}

func run_yk_hmac(chal string) (string, error) {
	dir, err := os.MkdirTemp("", "capclient")
	defer os.RemoveAll(dir)
	if err != nil {
		log.Println("Could not make temporary directory")
		return "", err
	}

	ykc := path.Join(dir, "ykchalresp")
	yki := path.Join(dir, "ykinfo")
	save(ykc, []byte("TODO"))
	save(yki, []byte("TODO"))

	// log.Println(ykc, "-2", "-H", "-x", chal)
	cmd := exec.Command(ykc, "-2", "-H", "-x", chal)
	output, err := cmd.Output()
	return hex.EncodeToString(output), err
}

func save(path string, content []byte) {
	err := os.WriteFile(path, content, 0755)
	if err != nil {
		log.Println("could not write ", path, " because: ", err)
	}
}
