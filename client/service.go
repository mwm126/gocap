package client

import (
	_ "embed"
	"encoding/json"
)

//go:embed services.json
var jsonFile []byte

type Service struct {
	Name          string
	Forwards      []string
	ServerAddress string
	ServerPort    int
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
