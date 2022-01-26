package watt

import (
	"log"
	"os"
	"os/exec"
)

func open_url(url string) {
	cmd := exec.Command("xdg-open", url)
	log.Println("\n\n\nOpen URL with: ", cmd)

	var err error
	var output []byte
	if os.Getenv("GOCAP_DEMO") == "" {
		output, err = cmd.CombinedOutput()
	}
	if err != nil {
		log.Printf("browser output: %s \nfrom error: %s  ", string(output), err)
	}
}
