package main

import (
	"encoding/json"
	"fmt"
	"os"
	//"flag"
	"errors"
	"io/ioutil"
	"log"
	"os/user"
	"strings"

	"go.i3wm.org/i3/v4"
)

const globalConfigFilename = ".MoveFloatContainer.json"

// TODO: this must not be constant <05-11-23, modernpacifist> //
const globalStatusBarDefaultHeight = 53

type NodeConfig struct {
	ID     i3.NodeID `json:"ID"`
	X      int64     `json:"X"`
	Y      int64     `json:"Y"`
	Width  int64     `json:"Width"`
	Height int64     `json:"Height"`
	Mark   string    `json:"Mark"`
}

func getNodeMark(node *i3.Node) string {
	if len(node.Marks) == 0 {
		return ""
	}
	return node.Marks[0]
}

func nodeConfigConstructor(node *i3.Node) NodeConfig {
	return NodeConfig{
		ID:     node.ID,
		X:      node.Rect.X,
		Y:      node.Rect.Y,
		Width:  node.Rect.Width,
		Height: node.Rect.Height,
		Mark:   getNodeMark(node),
	}
}

type JsonConfig struct {
	Location string                `json:"-"`
	Nodes    map[string]NodeConfig `json:"Nodes"`
}

func createJsonConfigFile(configFileLoc string) {
	var jsonConfig JsonConfig

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

func (jc *JsonConfig) Update(np NodeConfig) {
	jc.Nodes[np.Mark] = np
}

func (jc *JsonConfig) Dump() {
	jsonData, err := json.MarshalIndent(jc, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(jc.Location, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

//function save_preset() {
//updated_json=$(jq ".Preset = $1" $CONFIGFILE)
//echo $updated_json | jq . > $CONFIGFILE
//}

func getPrimaryOutput() []i3.Output {
	var res []i3.Output
	outputs, _ := i3.GetOutputs()
	for _, output := range outputs {
		if output.Active == true {
			res = append(res, output)
		}
	}
	return res
}

func getFocusedNode(i3Tree i3.Tree) *i3.Node {
	node := i3Tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Focused == true
	})

	if node == nil {
		log.Fatal(errors.New("Could not find focused node"))
	}

	return node
}

func getI3Tree() i3.Tree {
	tree, err := i3.GetTree()
	if err != nil {
		log.Fatal(err)
	}
	return tree
}

func resolveFileAbsolutePath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(fmt.Sprintf("~/%s", filename), "~", usr.HomeDir, 1)
}

func main() {
	//windowPosition := flag.Int("pos", 0, "New container preset of the window")
	//windowMark := flag.String("mark", "", "Set specified mark options on the focused container")
	//flag.Parse()

	//if *windowPosition == 0 {
	//log.Fatal("Window position was not specified as argument")
	//os.Exit(0)
	//}

	i3Tree := getI3Tree()
	node := getFocusedNode(i3Tree)
	//fmt.Println(node)
	ndconf := nodeConfigConstructor(node)
	fmt.Println(ndconf)

	//focusedNodeMark := getNodeMark(node)
	//fmt.Println(focusedNodeMark)
	// TODO: if the mark exists, then save the preset to json <05-11-23, modernpacifist> //

	absoluteConfigFilepath := resolveFileAbsolutePath(globalConfigFilename)
	// TODO: need to check if the file does not exist and create it <05-11-23, modernpacifist> //
	fmt.Println(absoluteConfigFilepath)

	createJsonConfigFile(absoluteConfigFilepath)

	//out := getPrimaryOutput()
	//fmt.Println(out)

	//fmt.Println(*windowPosition)
}
