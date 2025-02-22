package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"

	"go.i3wm.org/i3/v4"
)

const (
	ConfigDirectory string = "~/.config/"
)

var (
	CONFIG             JsonConfig
	MONITOR_DIMENSIONS i3.Rect
)

type WindowConfig struct {
	ID                  int64  `json:"id"`
	ResizedPlusYFlag    bool   `json:"resizedPlusYFlag"`
	ResizedMinusYFlag   bool   `json:"resizedMinusYFlag"`
	ResizedPlusXFlag    bool   `json:"resizedPlusXFlag"`
	ResizedMinusXFlag   bool   `json:"resizedMinusXFlag"`
	X                   int64  `json:"x"`
	Y                   int64  `json:"y"`
	Width               int64  `json:"width"`
	Height              int64  `json:"height"`
	Mark                string `json:"mark"`
	PreviousPlusYValue  int64  `json:"previousPlusYValue"`
	PreviousMinusYValue int64  `json:"previousMinusYValue"`
	PreviousPlusXValue  int64  `json:"previousPlusXValue"`
	PreviousMinusXValue int64  `json:"previousMinusXValue"`
}

func getPreviousResizeValues(node *i3.Node) map[string]int64 {
	resMap := make(map[string]int64)

	for _, n := range CONFIG.Windows {
		if n.ID == node.Window {
			resMap["plusY"] = n.PreviousPlusYValue
			resMap["minusY"] = n.PreviousMinusYValue
			resMap["plusX"] = n.PreviousPlusXValue
			resMap["minusX"] = n.PreviousMinusXValue
			break
		}
	}

	return resMap
}

func resolveResizedFlags(node *i3.Node, flagname string) bool {
	switch flagname {
	case "plusY":
		return node.Rect.Y == CONFIG.StatusBarHeight
	case "minusY":
		return node.Rect.Height == MONITOR_DIMENSIONS.Height-node.Rect.Y
	case "plusX":
		return node.Rect.Width == MONITOR_DIMENSIONS.Width-node.Rect.X
	case "minusX":
		return node.Rect.X == 0
	}
	return false
}

func getNodeMark(node *i3.Node) string {
	if len(node.Marks) == 0 {
		return ""
	}
	return node.Marks[0]
}

func WindowConfigConstructor(node *i3.Node) WindowConfig {
	previousResizeValuesMap := getPreviousResizeValues(node)
	plusY, _ := previousResizeValuesMap["plusY"]
	minusY, _ := previousResizeValuesMap["minusY"]
	plusX, _ := previousResizeValuesMap["plusX"]
	minusX, _ := previousResizeValuesMap["minusX"]

	// TODO: is the node does not contain a mark, just use a container id <17-11-23, modernpacifist> //
	nodeMark := getNodeMark(node)
	if nodeMark == "" {
		nodeMark = strconv.FormatInt(node.Window, 10)
	}

	return WindowConfig{
		ID:                  node.Window,
		ResizedPlusYFlag:    resolveResizedFlags(node, "plusY"),
		ResizedMinusYFlag:   resolveResizedFlags(node, "minusY"),
		ResizedPlusXFlag:    resolveResizedFlags(node, "plusX"),
		ResizedMinusXFlag:   resolveResizedFlags(node, "minusX"),
		X:                   node.Rect.X,
		Y:                   node.Rect.Y,
		Width:               node.Rect.Width,
		Height:              node.Rect.Height,
		Mark:                nodeMark,
		PreviousPlusYValue:  plusY,
		PreviousMinusYValue: minusY,
		PreviousPlusXValue:  plusX,
		PreviousMinusXValue: minusX,
	}
}

type JsonConfig struct {
	Location        string                  `json:"-"`
	StatusBarHeight int64                   `json:"statusBarHeight"`
	Windows         map[string]WindowConfig `json:"Windows"`
}

func JsonConfigConstructor(configFileLoc string) JsonConfig {
	_, err := os.Stat(configFileLoc)
	if os.IsNotExist(err) == true {
		createJsonConfigFile(configFileLoc)
	}

	jsonData, err := os.ReadFile(configFileLoc)
	if err != nil {
		log.Fatal(err)
	}

	var jsonConfig JsonConfig
	jsonConfig.Location = configFileLoc

	err = json.Unmarshal(jsonData, &jsonConfig)
	if err != nil {
		log.Fatal(err)
	}

	if jsonConfig.Windows == nil {
		jsonConfig.Windows = make(map[string]WindowConfig)
	}

	return jsonConfig
}

func (jc JsonConfig) Update(wc WindowConfig) {
	jc.Windows[wc.Mark] = wc
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

func (jc JsonConfig) Dump() {
	jsonData, err := json.MarshalIndent(jc, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(jc.Location, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func resolveJsonAbsolutePath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Replace(fmt.Sprintf("~/%s", filename), "~", usr.HomeDir, 1)
}
