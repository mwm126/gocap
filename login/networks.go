package login

import (
	"log"
	"net"

	externalip "github.com/glendc/go-external-ip"
)

func GetExternalIp() net.IP {
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		log.Println("Warning: Could not find external IP, ", err)
	}
	return ip
}
