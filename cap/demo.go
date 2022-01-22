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
	if strings.Contains(req, "ps auxnww|grep Xvnc|grep -v grep") {
		return `8227 27248  0.0  0.0  92620 70572 ?        S    Aug03   0:234 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :123 -desktop TurboVNC: login03:5 (the_user) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5905 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1`
	}
	if strings.Contains(req, "vncserver -geometry x -otp -novncauth -nohttpd") {
		return ""
	}
	log.Println("Unknown request: ", req)
	panic("FIXME")
}
