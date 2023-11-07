package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"go.i3wm.org/i3/v4"
)

const Mark = "last"
const ConfigFilename = ".BackAndForthContainer.json"

type Config struct {
	Path               string       `json:"-"`
	PreviousContainers [3]i3.NodeID `json:"PreviousContainersID"`
}

func configConstructor(filename string) Config {
	return Config{
		Path:               filename,
		PreviousContainers: [3]i3.NodeID{0, 0, 0},
	}
}

func (conf *Config) dump() {
	jsonData, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(conf.Path, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (conf *Config) update(node *i3.Node) {
	isDuplicate := false
	for _, value := range conf.PreviousContainers {
		if value == node.ID {
			isDuplicate = true
			break
		}
	}

	if !isDuplicate {
		for i := 0; i < len(conf.PreviousContainers)-1; i++ {
			conf.PreviousContainers[i] = conf.PreviousContainers[i+1]
		}
		conf.PreviousContainers[len(conf.PreviousContainers)-1] = node.ID
	}
}

func (conf *Config) readFromFile() {
	file, err := os.Open(conf.Path)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	err = json.Unmarshal(content, &conf)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}
}

func getI3Tree() i3.Tree {
	tree, err := i3.GetTree()
	if err != nil {
		log.Fatal(err)
	}

	return tree
}

func getFocusedNode(i3Tree i3.Tree) *i3.Node {
	node := i3Tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Focused == true
	})

	return node
}

func getPreviousNode(i3Tree i3.Tree) *i3.Node {
	node := i3Tree.Root.FindChild(func(n *i3.Node) bool {
		for _, m := range n.Marks {
			if m == Mark {
				return true
			}
		}
		return false
	})
	return node
}

func focus(node *i3.Node) {
	i3.RunCommand(fmt.Sprintf("[con_id=%d] focus", node.ID))
}

func mark(node *i3.Node) {
	i3.RunCommand(fmt.Sprintf("[con_id=%d] mark --add %s", node.ID, Mark))
}

func markID(nodeID i3.NodeID) {
	i3.RunCommand(fmt.Sprintf("[con_id=%d] mark --add %s", nodeID, Mark))
}

func RunDaemon(i3Tree i3.Tree, config Config) {
	recv := i3.Subscribe(i3.WindowEventType)
	for recv.Next() {
		ev := recv.Event().(*i3.WindowEvent)
		if ev.Change == "focus" {
			focusedNode := getFocusedNode(i3Tree)
			i3Tree = getI3Tree()
			mark(focusedNode)
			config.update(focusedNode)
			config.dump()
		}
	}
	log.Fatal(recv.Close())
}

func resolveAbsolutePath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Replace(fmt.Sprintf("~/%s", filename), "~", usr.HomeDir, 1)
}

func main() {
	var daemon bool

	flag.BoolVar(&daemon, "daemon", false, "description of the flag")
	flag.Parse()

	i3Tree := getI3Tree()

	absoluteConfigPath := resolveAbsolutePath(ConfigFilename)

	config := configConstructor(absoluteConfigPath)
	_, err := os.Stat(absoluteConfigPath)
	if os.IsNotExist(err) {
		config.dump()
	} else {
		config.readFromFile()
	}

	if daemon == true {
		RunDaemon(i3Tree, config)
		os.Exit(0)
	}

	focusedNode := getFocusedNode(i3Tree)

	var previousNode *i3.Node
	previousNode = getPreviousNode(i3Tree)
	if previousNode == nil {
		lastFocusedContainerID := config.PreviousContainers[len(config.PreviousContainers)-1]
		markID(lastFocusedContainerID)
		i3Tree = getI3Tree()
		previousNode = getPreviousNode(i3Tree)
	}

	mark(focusedNode)
	focus(previousNode)
}
