package login

import (
	_ "embed"
	"encoding/json"
)

//go:embed services.json
var jsonFile []byte

type Network struct {
	ClientExternalAddress string `json:"client_external_address"`
	CapServerAddress      string `json:"cap_server_address"`
}
type Service struct {
	Name     string
	Forwards []string
	CapPort  uint `json:"port"`
	Networks map[string]Network
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
	err := json.Unmarshal(jsonFile, &services)
	globalServices = services.Services
	return err
}

func FindServices() ([]Service, error) {
	return globalServices, nil
}
