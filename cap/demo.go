package cap

import (
	"log"
	"strings"
)

func demoReplies(req string) string {
	if strings.Contains(req, "echo") {
		return ""
	}
	if strings.Contains(req, "hostname") {
		return "demo_hostname"
	}
	if strings.Contains(req, "ping") {
		return "127.0.0.1" // LoginIP
	}

	if strings.Contains(req, "awk") {
		return "123" // UID
	}
	if strings.Contains(req, "openstack project list") {
		return `"ID","Name"
"3880f9d8def44f759324b8881d8fc736","ML_Demo"`
	}
	if strings.Contains(req, "openstack server list") {
		return `"ID","Name","Status","Networks","Image","Flavor"
	"1224e86d-f1b6-414d-ae3b-5475ebd291e4","hopeitworks","ERROR","","Windows 10","ml.gpu.1.bigmem"
	"8b65a1f3-9706-4653-ba8e-57b4307e5ad8","UbuntuEvan","SHUTOFF","ML_Demo_Network=172.16.1.104","Ubuntu-20.04-desktop.base","ml.tiny"
`
	}
	log.Println("Unknown request: ", req)
	panic("FIXME")
}
