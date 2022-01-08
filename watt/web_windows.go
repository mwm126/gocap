package watt

import (
	"log"
	"os/exec"
)

func open_url(url string) {
	cmd := exec.Command("start", url)
	log.Println("\n\n\nOpen URL with: ", cmd)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Println("browser output: ", string(output))
		log.Println("browser error: ", err)
	}

}
