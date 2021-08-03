package cap

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
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
	yfile, err := ioutil.ReadFile("networks.yaml")
	if err != nil {
		log.Fatal(err)
	}

	data := make(map[string]Network)
	err = yaml.Unmarshal(yfile, &data)
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
