package change_monitor_brightness

import (
	"encoding/json"
	"os"
	"log"
)

const configFilename = "~/.ScreenDim.json"

type Config struct {
	Brightness float64 `json:"brightness"`
}

func ConstructConfig(brightness float64) Config {
	return Config{
		Brightness: brightness,
	}
}

func (conf *Config) dump() {
	jsonData, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(configFilename, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (conf *Config) read() Config {
	jsonData, err := os.ReadFile(configFilename)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(jsonData, conf)
	if err != nil {
		log.Fatal(err)
	}

	return *conf
}
