package watt

import (
	"fmt"
	"log"
	"os/exec"
	"x/sys/windows"
)

func SpiceCmd(localPort uint) (*exec.Cmd, error) {

	log.Println("Looking for remote-viewer in PATH")

	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Classes\VirtViewer.vvfile\shell\open\command`, registry.QUERY_VALUE)
	if err != nil {
		log.Println("Could not find VirtViewer in Registry")
		return nil, err
	}
	defer k.Close()

	s, _, err := k.GetStringValue("SystemRoot")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Windows system root is %q\n", s)
	log.Fatal("WIP")
	virt_viewer = s

	return exec.Command(
		virt_viewer,
		fmt.Sprintf("spice://127.0.0.1::%d", localPort),
	), nil
}

func RunSpice(localPort uint) {
	cmd, err := SpiceCmd(localPort)
	if err != nil {
		log.Println("Could not run Spice")
	}
	log.Println("\n\n\nRunSpice: ", cmd)
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Println("spiceviewer output: ", string(output))
		log.Println("spiceviewer error: ", err)
	}

}
