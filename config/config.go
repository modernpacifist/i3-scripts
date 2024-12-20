package config

import (
	"encoding/json"
	"fmt"
	"os"

	"go.i3wm.org/i3/v4"
)

type JsonConfig struct {
	Location        string                  `json:"-"`
	StatusBarHeight int64                   `json:"statusBarHeight"`
	Windows         map[string]WindowConfig `json:"Windows"`
}

func createJsonConfigFile(configFileLoc string) {
	var jsonConfig JsonConfig
	// default value
	// jsonConfig.StatusBarHeight = 29

	file, err := os.Create(configFileLoc)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(jsonConfig)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
}

type NodeConfig struct {
	i3.Node
}
