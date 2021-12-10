package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

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
		globalServices = *init
		return nil
	}

	data, err := ioutil.ReadFile("./services.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var services Services
	err = json.Unmarshal(data, &services)
	globalServices = services.Services
	return err
}

func FindServices() ([]Service, error) {
	return globalServices, nil
}
