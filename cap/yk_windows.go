package cap

import (
	_ "embed"
	"encoding/hex"
	"log"
	"os"
	"os/exec"
	"path"
)

//go:generate go run gen.go
//go:generate unzip -o -d embeds embeds/ykpers-1.20.0-win64.zip

//go:embed embeds/bin/libjson-c-2.dll
var libjson []byte

//go:embed embeds/bin/libykpers-1-1.dll
var libykpers []byte

//go:embed embeds/bin/libyubikey-0.dll
var libyubikey []byte

//go:embed embeds/bin/ykchalresp.exe
var ykchalresp []byte

//go:embed embeds/bin/ykinfo.exe
var ykinfo []byte

func run_yk_info() ([]byte, error) {
	dir, err := os.MkdirTemp("", "capclient")
	defer os.RemoveAll(dir)
	if err != nil {
		log.Fatal("could not open tempfile", err)
	}

	ykc := path.Join(dir, "ykchalresp.exe")
	yki := path.Join(dir, "ykinfo.exe")
	save(path.Join(dir, "libjson-c-2.dll"), libjson)
	save(path.Join(dir, "libykpers-1-1.dll"), libykpers)
	save(path.Join(dir, "libyubikey-0.dll"), libyubikey)
	save(ykc, ykchalresp)
	save(yki, ykinfo)

	// log.Println(yki, "-s", "-q")
	cmd := exec.Command(yki, "-s", "-q")
	output, err := cmd.Output()
	return output, err
}

func run_yk_chalresp(chal string) ([]byte, error) {
	dir, err := os.MkdirTemp("", "capclient")
	defer os.RemoveAll(dir)
	if err != nil {
		log.Fatal("could not open tempfile", err)
	}

	ykc := path.Join(dir, "ykchalresp.exe")
	yki := path.Join(dir, "ykinfo.exe")
	save(path.Join(dir, "libjson-c-2.dll"), libjson)
	save(path.Join(dir, "libykpers-1-1.dll"), libykpers)
	save(path.Join(dir, "libyubikey-0.dll"), libyubikey)
	save(ykc, ykchalresp)
	save(yki, ykinfo)

	// log.Println(ykc, "-1", "-Y", "-x", chal)
	cmd := exec.Command(ykc, "-1", "-Y", "-x", chal)
	output, err := cmd.Output()
	return output, err
}

func run_yk_hmac(chal string) (string, error) {
	dir, err := os.MkdirTemp("", "capclient")
	defer os.RemoveAll(dir)
	if err != nil {
		log.Println("could not open tempfile", err)
		return "", err
	}

	ykc := path.Join(dir, "ykchalresp.exe")
	yki := path.Join(dir, "ykinfo.exe")
	save(path.Join(dir, "libjson-c-2.dll"), libjson)
	save(path.Join(dir, "libykpers-1-1.dll"), libykpers)
	save(path.Join(dir, "libyubikey-0.dll"), libyubikey)
	save(ykc, ykchalresp)
	save(yki, ykinfo)

	// log.Println(ykc, "-2", "-H", "-x", chal)
	cmd := exec.Command(ykc, "-2", "-H", "-x", chal)
	output, err := cmd.Output()
	return hex.EncodeToString(output), err
}

func save(path string, content []byte) {
	err := os.WriteFile(path, content, 0666)
	if err != nil {
		log.Fatal("could not write ", path, " because: ", err)
	}
}
