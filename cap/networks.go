package cap

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Network struct {
	DisableNTP        bool
	ExternalAddress   string
	Fe261Address      string
	JouleAddress      string
	ServerPort        int
	SessionMgtFwdAddr string
	WattAddress       string
}

func getNetworks() map[string]Network {

	yfile, err := ioutil.ReadFile("/home/mark/repos/gocap/networks.yaml")
	if err != nil {
		log.Fatal(err)
	}

	networks := make(map[string]Network)
	// data := make(map[string]Network)
	// data := make(map[string]interface{})
	err2 := yaml.Unmarshal(yfile, &networks)
	if err2 != nil {
		log.Fatal(err2)
	}

	// for k, v := range data {
	// 	switch val := v.(type) {
	// 	case map[interface{}]interface{}:
	// 		networks[k] = Network{
	// 			disableNTP:        true, // val["disableNTP"].(bool),
	// 			externalAddress:   net.ParseIP(val["externalAddress"].(string)),
	// 			fe261Address:      net.ParseIP(val["fe261Address"].(string)),
	// 			jouleAddress:      net.ParseIP(val["jouleAddress"].(string)),
	// 			serverPort:        val["serverPort"].(int),
	// 			sessionMgtFwdAddr: net.ParseIP(val["sessionMgtFwdAddr"].(string)),
	// 			wattAddress:       net.ParseIP(val["wattAddress"].(string)),
	// 		}
	// 		log.Println(k, val) // is []string
	// 		log.Println("DOG...!   ", val["externalAddress"])
	// 		log.Println("DOG...!   ", val["jouleAddress"])
	// 		log.Println("DOG...!   ", val["serverPort"])
	// 		log.Println("DOG...!   ", val["BOOBS"])
	// 	default:
	// 		log.Fatalf("Type unaccounted for: %+v\n", v)
	// 	}
	// }
	return networks

}
