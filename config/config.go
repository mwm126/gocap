package config

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Enable_fe261   bool
	Enable_joule   bool
	Enable_watt    bool
	Joule_Forwards []string
	Watt_Forwards  []string
	Fe261_Forwards []string
	DisableNtp     bool
	YubikeyTimeout uint
}

func NewConfig() Config {
	// Return default configuration

	return Config{
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
		DisableNtp:     false,
		YubikeyTimeout: 0,
	}
}

func GetConfig() Config {
	cfg_path := get_cfg_path()

	conf := NewConfig()

	if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
		log.Println("Config file does not exist; using defaults")
		return conf
	}

	data, err := ioutil.ReadFile(cfg_path)
	if err != nil {
		log.Println("Warning: Could not read Config file: ", cfg_path)
		return conf
	}
	err = toml.Unmarshal(data, &conf)
	if err != nil {
		log.Println("Warning: malformed config file: ", err)
		return conf
	}
	return conf
}

func WriteConfig(conf Config) {
	cfg_path := get_cfg_path()
	buf := bytes.Buffer{}
	encoder := toml.NewEncoder(&buf)
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
	_, err = file.Write(buf.Bytes())
	if err != nil {
		log.Println("Warning: Could not write config file: ", err)
		return
	}
}

func get_cfg_path() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	dir := filepath.Dir(exe)
	return filepath.Join(dir, "capclient.toml")
}

func SaveForwards(fwds []string) {
	cfg := GetConfig()
	cfg.Joule_Forwards = fwds[2:]
	WriteConfig(cfg)
}
