package cap

import (
	"log"
	"net"

	"github.com/glendc/go-external-ip"
	"gopkg.in/yaml.v3"
)

type User struct {
	Name       string
	Occupation string
	DisableNTP bool `yaml:"disableNTP"`
	DickSize   int
}

type Network struct {
	DisableNTP        bool   `yaml:"disableNTP"`
	ExternalAddress   string `yaml:"externalAddress"`
	Fe261Address      string `yaml:"fe261Address"`
	JouleAddress      string `yaml:"jouleAddress"`
	ServerPort        int    `yaml:"serverPort"`
	SessionMgtFwdAddr string `yaml:"sessionMgtFwdAddr"`
	WattAddress       string `yaml:"wattAddress"`
}

func getNetworks() map[string]Network {
	data := make(map[string]Network)

	err := yaml.Unmarshal([]byte(networks_yaml), &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
func getNetworkNames() []string {
	netMap := getNetworks()
	names := make([]string, 0, len(netMap))

	for name := range netMap {
		names = append(names, name)
	}
	return names
}

var networks_yaml = `
alb_admin:
  disableNTP: True
  externalAddress: "198.99.249.252"
  fe261Address: "204.154.140.51"
  jouleAddress: "204.154.139.11"
  serverPort: 62201
  sessionMgtFwdAddr: "172.16.0.1"
  wattAddress: "204.154.140.10"

alb_foreign_national:
  disableNTP: True
  externalAddress: "205.254.146.76"
  fe261Address: "204.154.140.51"
  jouleAddress: "204.154.139.11"
  serverPort: 62201
  sessionMgtFwdAddr: "172.16.0.1"
  wattAddress: "204.154.140.10"

alb_research:
  disableNTP: True
  externalAddress: "198.99.249.252"
  fe261Address: "204.154.140.51"
  jouleAddress: "204.154.139.11"
  serverPort: 62201
  sessionMgtFwdAddr: "172.16.0.1"
  wattAddress: "204.154.140.10"

alb_scientific:
  disableNTP: True
  externalAddress: "198.99.249.248"
  fe261Address: "204.154.140.51"
  jouleAddress: "204.154.139.11"
  serverPort: 62201
  sessionMgtFwdAddr: "172.16.0.1"
  wattAddress: "204.154.140.10"

external:
  disableNTP: True
  externalAddress: ""
  fe261Address: "204.154.140.51"
  jouleAddress: "204.154.139.11"
  serverPort: 62201
  sessionMgtFwdAddr: "172.16.0.1"
  wattAddress: "204.154.140.10"

mgn_admin:
  fe261Address: "204.154.140.51"
  wattAddress: "204.154.140.10"
  sessionMgtFwdAddr: "172.16.0.1"
  jouleAddress: "204.154.139.11"
  serverPort: 62201
  externalAddress: "204.154.136.254"
  disableNTP: True

mgn_foreign_national:
  disableNTP: True
  externalAddress: "205.254.146.76"
  fe261Address: "204.154.140.51"
  jouleAddress: "204.154.139.11"
  serverPort: 62201
  sessionMgtFwdAddr: "172.16.0.1"
  wattAddress: "204.154.140.10"

mgn_research:
  disableNTP: True
  externalAddress: "204.154.136.254"
  fe261Address: "204.154.140.51"
  jouleAddress: "204.154.139.11"
  serverPort: 62201
  sessionMgtFwdAddr: "172.16.0.1"
  wattAddress: "204.154.140.10"

mgn_scientific:
  disableNTP: True
  externalAddress: "204.154.139.10"
  fe261Address: "204.154.140.51"
  jouleAddress: "204.154.139.11"
  serverPort: 62201
  sessionMgtFwdAddr: "172.16.0.1"
  wattAddress: "204.154.140.10"

pgh_admin:
  disableNTP: True
  externalAddress: "204.154.137.254"
  fe261Address: "204.154.140.51"
  jouleAddress: "204.154.139.11"
  serverPort: 62201
  sessionMgtFwdAddr: "172.16.0.1"
  wattAddress: "204.154.140.10"

pgh_foreign_national:
  disableNTP: True
  externalAddress: "204.154.137.254"
  fe261Address: "204.154.140.51"
  jouleAddress: "204.154.139.11"
  serverPort: 62201
  sessionMgtFwdAddr: "172.16.0.1"
  wattAddress: "204.154.140.10"

pgh_research:
  disableNTP: True
  externalAddress: "204.154.137.254"
  fe261Address: "204.154.140.51"
  jouleAddress: "204.154.139.11"
  serverPort: 62201
  sessionMgtFwdAddr: "172.16.0.1"
  wattAddress: "204.154.140.10"

pgh_scientific:
  disableNTP: True
  externalAddress: "198.99.246.197"
  fe261Address: "198.99.246.146"
  jouleAddress: "198.99.246.146"
  serverPort: 62201
  sessionMgtFwdAddr: "172.16.0.1"
  wattAddress: "198.99.246.146"

vpn:
  disableNTP: True
  externalAddress: "199.249.243.253"
  fe261Address: "198.99.246.146"
  jouleAddress: "198.99.246.146"
  serverPort: 62201
  sessionMgtFwdAddr: "172.16.0.1"
  wattAddress: "199.249.243.253"
`

func GetExternalIp() net.IP {
	consensus := externalip.DefaultConsensus(nil, nil)
	ip, err := consensus.ExternalIP()
	if err != nil {
		log.Println("Warning: Could not find external IP, ", err)
	}
	return ip
}
