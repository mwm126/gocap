package login

import (
	_ "embed"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
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

func InitServices(init *[]Service) error {
	if init != nil {
		// override for testing
		globalServices = *init
		return nil
	}

	var services Services

	response, err := http.Get(services_url)
	if err != nil {
		log.Println(err)
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	err = json.Unmarshal(body, &services)
	if err != nil {
		log.Println("Unable to parse services.json: ", err)
		return err
	}
	globalServices = services.Services
	return nil
}

func FindServices() ([]Service, error) {
	return globalServices, nil
}
