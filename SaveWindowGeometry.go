package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"go.i3wm.org/i3/v4"
)

type NodeConfig struct {
	X      int64  `json:"X"`
	Y      int64  `json:"Y"`
	Width  int64  `json:"Width"`
	Height int64  `json:"Height"`
	Mark   string `json:"Mark"`
}

func nodeConfigConstructor(node *i3.Node) NodeConfig {
	return NodeConfig{
		X:      node.Rect.X,
		Y:      node.Rect.Y,
		Width:  node.Rect.Width,
		Height: node.Rect.Height,
		Mark:   getNodeMark(node),
	}
}

func getI3Tree() i3.Tree {
	tree, err := i3.GetTree()
	if err != nil {
		log.Fatal(err)
	}
	return tree
}

func getFocusedNode() *i3.Node {
	i3Tree := getI3Tree()

	node := i3Tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Focused == true
	})

	if node == nil {
		log.Fatal(errors.New("Could not find focused node"))
	}

	return node
}

func getNodeMark(node *i3.Node) string {
	if len(node.Marks) == 0 {
		return ""
	}
	return node.Marks[0]
}

type Config struct {
	Location string                `json:"-"`
	Nodes    map[string]NodeConfig `json:"Nodes"`
}

func createConfigFile(configFileLoc string) {
	var config Config

	file, err := os.Create(configFileLoc)
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
}

func ConfigConstructor(configFileLoc string) Config {
	_, err := os.Stat(configFileLoc)
	if os.IsNotExist(err) == true {
		createConfigFile(configFileLoc)
	}

	jsonData, err := ioutil.ReadFile(configFileLoc)
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	config.Location = configFileLoc

	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		log.Fatal(err)
	}

	if config.Nodes == nil {
		config.Nodes = make(map[string]NodeConfig)
	}

	return config
}

func (jc *Config) Update(np NodeConfig) {
	jc.Nodes[np.Mark] = np
}

func (jc *Config) Dump() {
	jsonData, err := json.MarshalIndent(jc, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(jc.Location, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func resolveFileAbsolutePath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(fmt.Sprintf("~/%s", filename), "~", usr.HomeDir, 1)
}

func restoreWindowWithParameters(nodeConfig NodeConfig, mark string) {
	cmd := fmt.Sprintf("mark \"%s\", move scratchpad, [con_mark=\"^%s$\"] scratchpad show, move absolute position %d %d, resize set %d %d", mark, mark, nodeConfig.X, nodeConfig.Y - 24, nodeConfig.Width, nodeConfig.Height + 24)
	i3.RunCommand(cmd)
}

func main() {
	var restore string

	flag.StringVar(&restore, "restore", "", "Specify the mark to restore")

	flag.Parse()

	configPath := resolveFileAbsolutePath(".SaveWindowGeometry.json")
	config := ConfigConstructor(configPath)

	if restore != "" {
		value, exists := config.Nodes[restore]
		if exists == true {
			restoreWindowWithParameters(value, restore)
		}

		os.Exit(0)
	}

	focusedNode := getFocusedNode()
	nodeConfig := nodeConfigConstructor(focusedNode)
	config.Update(nodeConfig)
	config.Dump()
}
