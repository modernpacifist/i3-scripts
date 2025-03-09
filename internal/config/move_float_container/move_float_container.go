package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"go.i3wm.org/i3/v4"

	"github.com/modernpacifist/i3-scripts-go/internal/i3operations"
)

const (
	ConfigFilename string = "~/.MoveFloatContainer.json"
	StatusBarHeight int = 53
)

type NodeConfig struct {
	i3.Node
	Marks []string
}


func NodeConfigConstructor(node *i3.Node) NodeConfig {
	return NodeConfig{
		Node:  *node,
		Marks: i3operations.GetNodeMarks(node),
	}
}

type JsonConfig struct {
	Location string                `json:"-"`
	Nodes    map[string]NodeConfig `json:"Nodes"`
}

// func CreateJsonConfigFile[T any](config T, configFileLoc string) {
// 	file, err := os.Create(configFileLoc)
// 	if err != nil {
// 		log.Fatal("Error creating file:", err)
// 	}
// 	defer file.Close()
// 	encoder := json.NewEncoder(file)
// 	err = encoder.Encode(config)
// 	if err != nil {
// 		fmt.Println("Error encoding JSON:", err)
// 		return
// 	}
// }

func CreateJsonConfigFile(jsonConfig JsonConfig, configFileLoc string) {
	// var jsonConfig JsonConfig

	file, err := os.Create(configFileLoc)
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(jsonConfig)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
}
