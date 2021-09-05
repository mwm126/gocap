package cap

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type config struct {
	Enable_fe261   bool
	Enable_joule   bool
	Enable_watt    bool
	Joule_Forwards []string
	Watt_Forwards  []string
	Fe261_Forwards []string
}

func NewConfig() config {
	// Return default configuration

	return config{
		Enable_fe261:   false,
		Enable_joule:   true,
		Enable_watt:    true,
		Joule_Forwards: []string{
			// "???:127.0.0.1:???", // VNC port
		},
		Watt_Forwards: []string{
			"6082:192.168.101.200:6082", // Spice port
			"10080:192.168.101.200:80",  // Web port
		},
		Fe261_Forwards: []string{
			"13389:172.16.1.11:3389",  // RDP 1 port
			"23389:172.16.1.12:3389",  // RDP 2 port
			"13306:172.16.1.13:3306",  // SQL port
			"19392:172.16.1.10:9392",  // OpenVAS port
			"18834:172.16.1.10:8834",  // Nessus port
			"18888:172.16.1.11:18888", // Plexos port
			"1399:172.16.1.11:1399",   // Plexos port 2
		},
	}
}

func GetConfig() config {
	cfg_path := get_cfg_path()

	conf := NewConfig()
	data, err := ioutil.ReadFile(cfg_path)
	if err != nil {
		log.Println("Warning: Could not read config file: ", cfg_path)
	} else {
		toml.Unmarshal(data, &conf)
	}
	return conf
}

func WriteConfig(conf config) {
	cfg_path := get_cfg_path()
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	encoder.ArraysWithOneElementPerLine(true)
	if err := encoder.Encode(conf); err != nil {
		log.Fatal(err)
	}
	log.Println("Saving ", cfg_path, ":\n", buf.String())
	file, err := os.Create(cfg_path)
	if err != nil {
		log.Println("Warning: Could not save config file: ", cfg_path)
		return
	}
	file.Write(buf.Bytes())
}

func get_cfg_path() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	dir := filepath.Dir(exe)
	return filepath.Join(dir, "capclient.toml")
}
