package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"

	"i3-scripts-go/utils"

	"go.i3wm.org/i3/v4"
)

const (
	ConfigFilename string = ".ManageFloatContainer.json"
)

type ContainerParameters struct {
	// TODO: add a node id and adress some of the containers by it <13-11-23, modernpacifist> //
	ID     i3.NodeID `json:"ID"`
	X      int64     `json:"X"`
	Y      int64     `json:"Y"`
	Width  int64     `json:"Width"`
	Height int64     `json:"Height"`
	Marks  []string  `json:"Marks"`
	// TODO: having a mark field here is extremely idiotic, since we have a map with keys of marks themselves <13-11-23, modernpacifist> //
}

func getContainerMarks(node *i3.Node) []string {
	// TODO: a bug here if the window contains more than one mark <13-11-23, modernpacifist> //
	if len(node.Marks) == 0 {
		return nil
	}
	return node.Marks
}

func containerParametersConstructor(node *i3.Node) ContainerParameters {
	conMarks := getContainerMarks(node)
	//if conMarks == nil {
	//conMarks = []string
	////log.Fatal("This node does not have marks")
	//}
	//if len(conMarks) == 0 {
	//conMarks = []string
	////log.Fatal("This node does not have marks")
	//}

	return ContainerParameters{
		ID:     node.ID,
		X:      node.Rect.X,
		Y:      node.Rect.Y,
		Width:  node.Rect.Width,
		Height: node.Rect.Height,
		Marks:  conMarks,
	}
}

func containerParametersFactory(node *i3.Node) ContainerParameters {
	conMarks := getContainerMarks(node)
	//if conMarks == nil {
	//conMarks = []string
	////log.Fatal("This node does not have marks")
	//}
	//if len(conMarks) == 0 {
	//conMarks = []string
	////log.Fatal("This node does not have marks")
	//}

	return ContainerParameters{
		ID:     node.ID,
		X:      node.Rect.X,
		Y:      node.Rect.Y,
		Width:  node.Rect.Width,
		Height: node.Rect.Height,
		Marks:  conMarks,
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
	Location string                         `json:"-"`
	Nodes    map[string]ContainerParameters `json:"Nodes"`
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

	jsonData, err := os.ReadFile(configFileLoc)
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
		config.Nodes = make(map[string]ContainerParameters)
	}

	return config
}

func (conf *Config) Update(cp ContainerParameters, mark string) {
	conf.Nodes[mark] = cp
}

func (conf *Config) UpdateID(cp ContainerParameters, mark string) {
	if entry, ok := conf.Nodes[mark]; ok {
		temp := entry
		temp.ID = cp.ID
		conf.Nodes[mark] = temp
	}
}

func (conf *Config) Dump() {
	jsonData, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(conf.Location, jsonData, 0644)
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

func createFloatingContainer(conParams ContainerParameters, mark string) {
	cmd := fmt.Sprintf("mark \"%s\", move scratchpad, [con_mark=\"^%s$\"] scratchpad show, move absolute position %d %d, resize set %d %d", mark, mark, conParams.X, conParams.Y, conParams.Width, conParams.Height)
	i3.RunCommand(cmd)
}

func createFloatingContainerDefault(conParams ContainerParameters, mark string) {
	//cmd := fmt.Sprintf("mark --add \"%s\", floating enable, resize set %d %d, move position center", mark, conParams.Width, conParams.Height+24)
	cmd := fmt.Sprintf("mark \"%s\", move scratchpad, [con_mark=\"^%s$\"] scratchpad show, floating enable, resize set %d %d, move position center", mark, mark, conParams.Width, conParams.Height)
	i3.RunCommand(cmd)
}

func showContainer(conParams ContainerParameters) {
	cmd := fmt.Sprintf("[con_id=%d] scratchpad show, move absolute position %d %d, resize set %d %d", conParams.ID, conParams.X, conParams.Y, conParams.Width, conParams.Height)
	i, e := i3.RunCommand(cmd)
	fmt.Println(i)
	fmt.Println(e)
}

func main() {
	var restoreFlag string
	var showFlag string
	var updateFlag string
	var saveFlag bool

	// TODO: use only one argument and use switch statement <15-11-23, modernpacifist> //
	flag.StringVar(&restoreFlag, "restore", "", "Specify the mark to restore")
	flag.StringVar(&showFlag, "show", "", "Specify the mark to show")
	flag.StringVar(&updateFlag, "update", "", "Specify the mark to show")
	flag.BoolVar(&saveFlag, "save", false, "Specify the mark to show")

	flag.Parse()

	absoluteConfigPath := resolveFileAbsolutePath(ConfigFilename)
	config := ConfigConstructor(absoluteConfigPath)

	if restoreFlag != "" {
		containerParameters, exists := config.Nodes[restoreFlag]
		if exists == false {
			// TODO: restore with some arbitrary values <17-11-23, modernpacifist> //
			focusedNode := getFocusedNode()
			containerParameters = containerParametersConstructor(focusedNode)

			containerParameters.Width = 2000
			containerParameters.Height = 1000
			containerParameters.Marks = []string{restoreFlag}

			createFloatingContainerDefault(containerParameters, restoreFlag)

			config.Update(containerParameters, restoreFlag)
			config.Dump()

			os.Exit(0)
		}
		createFloatingContainer(containerParameters, restoreFlag)
	}

	if showFlag != "" {
		node, exists := config.Nodes[showFlag]
		if exists == false {
			os.Exit(0)
		}
		showContainer(node)
	}

	if updateFlag != "" {
		markedNode := getNodeWithMark(updateFlag)
		fmt.Println(markedNode)
		if markedNode == nil {
			os.Exit(0)
		}

		var nodeConfig ContainerParameters
		// if the node with this mark does not exist add it to config
		configNode, exists := config.Nodes[updateFlag]
		fmt.Println(configNode)
		//nodeConfig, exists := config.Nodes[updateFlag]
		if exists == false {
			nodeConfig = containerParametersConstructor(markedNode)
			//nodeConfig.Mark = updateFlag
			config.Update(nodeConfig, updateFlag)
		}

		nodeConfig = containerParametersConstructor(markedNode)
		//nodeConfig.Mark = updateFlag
		config.UpdateID(nodeConfig, updateFlag)
		config.Dump()
	}

	if saveFlag == true {
		focusedNode := getFocusedNode()
		nodeConfig := containerParametersConstructor(focusedNode)
		config.Update(nodeConfig, nodeConfig.Marks[0])
		config.Dump()

		utils.NotifySend(0.5, fmt.Sprintf("Saved mark %s", nodeConfig.Marks[0]))
	}
}
