package cap

import (
	"bytes"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func get_cfg_path() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}
	dir := filepath.Dir(exe)
	return filepath.Join(dir, "capclient.toml")
}

type config struct {
	Enable_fe261   bool
	Enable_joule   bool
	Enable_watt    bool
	Joule_Forwards []string
}

func GetConfig() config {
	cfg_path := get_cfg_path()
	log.Println("cfg_path = ", cfg_path)

	var conf config
	md, err := toml.DecodeFile(cfg_path, &conf)
	if err != nil {
		log.Println("Warning: Could not read config file")
	}
	log.Printf("Undecoded keys: %q\n", md.Undecoded())
	return conf
}

func WriteConfig(conf config) {
	cfg_path := get_cfg_path()
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(conf); err != nil {
		log.Fatal(err)
	}
	log.Println("Debug: ", buf.String())
	file, err := os.Create(cfg_path)
	if err != nil {
		log.Println("Warning: Could not save config file: ", cfg_path)
		return
	}
	file.Write(buf.Bytes())
}
