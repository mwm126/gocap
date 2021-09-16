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
	CapPort        uint
	DisableNtp     bool
	External_Ips   map[string]string
	Fe261_Ips      map[string]string
	Joule_Ips      map[string]string
	Watt_Ips       map[string]string
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
		CapPort:    62201, // Unlikely to change
		DisableNtp: false,
		External_Ips: map[string]string{
			"alb_admin":            "198.99.249.252",
			"alb_foreign_national": "205.254.146.76",
			"alb_research":         "198.99.249.252",
			"alb_scientific":       "198.99.249.248",
			"external":             "0.0.0.0",
			"mgn_admin":            "204.154.136.254",
			"mgn_foreign_national": "205.254.146.76",
			"mgn_research":         "204.154.136.254",
			"mgn_scientific":       "204.154.139.10",
			"pgh_admin":            "204.154.137.254",
			"pgh_foreign_national": "204.154.137.254",
			"pgh_research":         "204.154.137.254",
			"pgh_scientific":       "198.99.246.197",
			"vpn":                  "199.249.243.253",
		},
		Fe261_Ips: map[string]string{
			"alb_admin":            "204.154.140.51",
			"alb_foreign_national": "204.154.140.51",
			"alb_research":         "204.154.140.51",
			"alb_scientific":       "204.154.140.51",
			"external":             "0.0.0.0",
			"mgn_admin":            "204.154.140.51",
			"mgn_foreign_national": "204.154.140.51",
			"mgn_research":         "204.154.140.51",
			"mgn_scientific":       "204.154.140.51",
			"pgh_admin":            "204.154.140.51",
			"pgh_foreign_national": "204.154.140.51",
			"pgh_research":         "204.154.140.51",
			"pgh_scientific":       "198.99.246.146",
			"vpn":                  "198.99.246.146",
		},
		Joule_Ips: map[string]string{
			"alb_admin":            "204.154.139.11",
			"alb_foreign_national": "204.154.139.11",
			"alb_research":         "204.154.139.11",
			"alb_scientific":       "204.154.139.11",
			"external":             "0.0.0.0",
			"mgn_admin":            "204.154.139.11",
			"mgn_foreign_national": "204.154.139.11",
			"mgn_research":         "204.154.139.11",
			"mgn_scientific":       "204.154.139.11",
			"pgh_admin":            "204.154.139.11",
			"pgh_foreign_national": "204.154.139.11",
			"pgh_research":         "204.154.139.11",
			"pgh_scientific":       "198.99.246.146",
			"vpn":                  "199.249.243.253",
		},
		Watt_Ips: map[string]string{
			"alb_admin":            "204.154.140.10",
			"alb_foreign_national": "204.154.140.10",
			"alb_research":         "204.154.140.10",
			"alb_scientific":       "204.154.140.10",
			"external":             "0.0.0.0",
			"mgn_admin":            "204.154.140.10",
			"mgn_foreign_national": "204.154.140.10",
			"mgn_research":         "204.154.140.10",
			"mgn_scientific":       "204.154.140.10",
			"pgh_admin":            "204.154.140.10",
			"pgh_foreign_national": "204.154.140.10",
			"pgh_research":         "204.154.140.10",
			"pgh_scientific":       "198.99.246.146",
			"vpn":                  "199.249.243.253",
		},
	}
}

func GetConfig() config {
	cfg_path := get_cfg_path()

	conf := NewConfig()

	if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
		log.Println("Config file does not exist; using defaults")
		return conf
	}

	data, err := ioutil.ReadFile(cfg_path)
	if err != nil {
		log.Println("Warning: Could not read config file: ", cfg_path)
		return conf
	}
	toml.Unmarshal(data, &conf)
	return conf
}

func WriteConfig(conf config) {
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

func SaveForwards(fwds []string) {
	cfg := GetConfig()
	cfg.Joule_Forwards = fwds[2:]
	WriteConfig(cfg)
}
