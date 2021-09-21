package cap

import (
	_ "embed"
	"log"
	"os"
	"os/exec"
	"path"
)

//go:generate go run gen.go

func run_yk_info() ([]byte, error) {
	dir, err := os.MkdirTemp("", "capclient")
	defer os.RemoveAll(dir)
	if err != nil {
		log.Fatal("could not open tempfile", err)
	}

	yki := path.Join(dir, "ykinfo")
	save(yki, []byte("TODO"))

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
		log.Fatal("could not open tempfile", err)
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
	err := os.WriteFile(path, content, 0666)
	if err != nil {
		log.Fatal("could not write ", path, " because: ", err)
	}
}
