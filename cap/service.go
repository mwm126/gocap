package cap

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	externalip "github.com/glendc/go-external-ip"
)

var services_url = "https://hpc.netl.doe.gov/cap/services.json"

type Network struct {
	ClientExternalAddress string `json:"client_external_address"`
	CapServerAddress      string `json:"cap_server_address"`
}

// The Service represents a server that can be connect to with CAP, such as Joule, Watt, or Vee
type Service struct {
	Name     string
	Forwards []string
	CapPort  uint `json:"capport"`
	SshPort  uint `json:"sshport"`
	Networks map[string]Network
}

// For a given Network (MGN-Admin, PIT-Admin, MGN-Science, etc.) return the external address and CAP server address (for this service)
func (s *Service) FindAddresses(network string) (net.IP, net.IP) {
	network_ips := make(map[string]string)
	external_ips := make(map[string]string)
	for name, val := range s.Networks {
		network_ips[name] = val.CapServerAddress
		external_ips[name] = val.ClientExternalAddress
	}

	var ext_addr, server_addr net.IP
	if network == "external" {
		ext_addr = GetExternalIp()
	} else {
		ext_addr = net.ParseIP(external_ips[network])
	}
	server_addr = net.ParseIP(network_ips[network])
	return ext_addr, server_addr
}

type Services struct {
	Services []Service
}

var globalServices []Service
var demoPort uint

func SetDemoServices(services []Service) {
	globalServices = services
}

func SetDemoPort(port uint) {
	demoPort = port
}

func defaultDemoServices() []Service {
	sshport := demoPort
	capport := sshport // doesn't matter; ignored anyway
	return []Service{{
		Name:    "joule",
		CapPort: uint(capport),
		SshPort: sshport,
		Networks: map[string]Network{
			"external": {
				ClientExternalAddress: "127.0.0.1",
				CapServerAddress:      "127.0.0.1",
			},
		},
	},
		{
			Name:    "watt",
			CapPort: uint(capport),
			SshPort: sshport,
			Networks: map[string]Network{
				"external": {
					ClientExternalAddress: "127.0.0.1",
					CapServerAddress:      "127.0.0.1",
				},
			},
		},
	}
}

func (cm *ConnectionManager) FindServices(uname string) ([]Service, error) {
	if os.Getenv("GOCAP_DEMO") != "" && globalServices == nil {
		globalServices = defaultDemoServices()
	}

	if globalServices != nil {
		return globalServices, nil
	}

	var services Services

	timestamp, err := getNtpTime()
	if err != nil {
		log.Printf("Unable to get NTP time:  %v", err)
		log.Printf("Warning: going to use local time, without checking for NTP offset")
		timestamp = (int32)(time.Now().Unix())
	}
	auth_addr := net.IPv4(1, 1, 1, 1)
	ssh_addr := net.IPv4(2, 2, 2, 2)
	server_addr := net.IPv4(3, 3, 3, 3)
	packet, err := cm.knocker.makePacket(uname, timestamp, auth_addr, ssh_addr, server_addr)
	if err != nil {
		log.Println("Unable to get services:", err)
		return globalServices, err
	}

	response, err := http.Post(services_url, "application/json", bytes.NewReader(packet))
	if err != nil {
		log.Println(err)
		return globalServices, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return globalServices, err
	}

	err = json.Unmarshal(body, &services)
	if err != nil {
		log.Println("Unable to parse services.json: ", err)
		return globalServices, err
	}
	globalServices = services.Services
	return globalServices, nil
}

func GetExternalIp() net.IP {
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		log.Println("Warning: Could not find external IP, ", err)
	}
	return ip
}
