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

const (
	ConfigFilename string = ".ManageFloatContainer.json"
)

type NodeConfig struct {
	// TODO: add a node id and adress some of the containers by it <13-11-23, modernpacifist> //
	ID     i3.NodeID `json:"ID"`
	X      int64     `json:"X"`
	Y      int64     `json:"Y"`
	Width  int64     `json:"Width"`
	Height int64     `json:"Height"`
	Marks  []string  `json:"Marks"`
	// TODO: having a mark field here is extremely idiotic, since we have a map with keys of marks themselves <13-11-23, modernpacifist> //
}

// func getNodeMark(node *i3.Node) string {
func getNodeMarks(node *i3.Node) []string {
	// TODO: a bug here if the window contains more than one mark <13-11-23, modernpacifist> //
	if len(node.Marks) == 0 {
		return nil
	}
	return node.Marks
}

func nodeConfigConstructor(node *i3.Node) NodeConfig {
	mark := getNodeMarks(node)
	//if mark == "" {
	if mark == nil {
		log.Fatal("This node does not have marks")
	}

	return NodeConfig{
		ID:     node.ID,
		X:      node.Rect.X,
		Y:      node.Rect.Y,
		Width:  node.Rect.Width,
		Height: node.Rect.Height,
		Marks:  getNodeMarks(node),
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
		log.Fatal(errors.New("Could not get focused node"))
	}

	return node
}

func getNodeWithMark(mark string) *i3.Node {
	i3Tree := getI3Tree()

	node := i3Tree.Root.FindChild(func(n *i3.Node) bool {
		for _, m := range n.Marks {
			if m == mark {
				return true
			}
		}
		return false
	})
	return node
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

func (jc *Config) Update(np NodeConfig, mark string) {
	//jc.Nodes[np.Mark] = np
	jc.Nodes[mark] = np
}

// func (jc *Config) UpdateID(np NodeConfig) {
func (jc *Config) UpdateID(np NodeConfig, mark string) {
	//if entry, ok := jc.Nodes[np.Mark]; ok {
	if entry, ok := jc.Nodes[mark]; ok {
		temp := entry
		temp.ID = np.ID
		//jc.Nodes[np.Mark] = temp
		jc.Nodes[mark] = temp
	}
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
	cmd := fmt.Sprintf("mark \"%s\", move scratchpad, [con_mark=\"^%s$\"] scratchpad show, move absolute position %d %d, resize set %d %d", mark, mark, nodeConfig.X, nodeConfig.Y-24, nodeConfig.Width, nodeConfig.Height+24)
	i3.RunCommand(cmd)
}

func showWindowWithParameters(nodeConfig NodeConfig) {
	cmd := fmt.Sprintf("[con_id=%d] scratchpad show, move absolute position %d %d, resize set %d %d", nodeConfig.ID, nodeConfig.X, nodeConfig.Y, nodeConfig.Width, nodeConfig.Height)
	i, e := i3.RunCommand(cmd)
	fmt.Println(i)
	fmt.Println(e)
}

func main() {
	var restoreFlag string
	var showFlag string
	var updateFlag string
	var resetFlag bool

	// TODO: use only one argument and use switch statement <15-11-23, modernpacifist> //
	flag.StringVar(&restoreFlag, "restore", "", "Specify the mark to restore")
	flag.StringVar(&showFlag, "show", "", "Specify the mark to show")
	flag.StringVar(&updateFlag, "update", "", "Specify the mark to show")
	flag.BoolVar(&resetFlag, "reset", false, "Specify the mark to show")

	flag.Parse()

	absoluteConfigPath := resolveFileAbsolutePath(ConfigFilename)
	config := ConfigConstructor(absoluteConfigPath)

	if restoreFlag != "" {
		value, exists := config.Nodes[restoreFlag]
		if exists == false {
			os.Exit(0)
		}
		restoreWindowWithParameters(value, restoreFlag)
	}

	if showFlag != "" {
		node, exists := config.Nodes[showFlag]
		if exists == false {
			os.Exit(0)
		}
		showWindowWithParameters(node)
	}

	if updateFlag != "" {
		markedNode := getNodeWithMark(updateFlag)
		fmt.Println(markedNode)
		if markedNode == nil {
			os.Exit(0)
		}

		var nodeConfig NodeConfig
		// if the node with this mark does not exist add it to config
		configNode, exists := config.Nodes[updateFlag]
		fmt.Println(configNode)
		//nodeConfig, exists := config.Nodes[updateFlag]
		if exists == false {
			nodeConfig = nodeConfigConstructor(markedNode)
			//nodeConfig.Mark = updateFlag
			config.Update(nodeConfig, updateFlag)
		}

		nodeConfig = nodeConfigConstructor(markedNode)
		//nodeConfig.Mark = updateFlag
		config.UpdateID(nodeConfig, updateFlag)
		config.Dump()
		os.Exit(0)
	}

	focusedNode := getFocusedNode()
	nodeConfig := nodeConfigConstructor(focusedNode)
	config.Update(nodeConfig, nodeConfig.Marks[0])
	config.Dump()
}
