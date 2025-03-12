package scale_float_container

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"go.i3wm.org/i3/v4"
)

const (
	configFilename  string      = "~/.ScaleFloatWindow.json"
	defaultPerms    os.FileMode = 0644
	statusBarHeight int         = 53
)

type Config struct {
	Path  string                  `json:"-"`
	Nodes map[string]WindowConfig `json:"nodes"`
}

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

type NodeConfig struct {
	i3.Node
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
		// don't believe I wrote this
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

func JsonConfigConstructor(configFileLoc string) Config {
	var jsonConfig Config

	_, err := os.Stat(configFileLoc)
	if os.IsNotExist(err) == true {
		createJsonConfigFile(configFileLoc)
	}

	jsonData, err := os.ReadFile(configFileLoc)
	if err != nil {
		log.Fatal(err)
	}

	jsonConfig.Path = configFileLoc

	err = json.Unmarshal(jsonData, &jsonConfig)
	if err != nil {
		log.Fatal(err)
	}

	if jsonConfig.Nodes == nil {
		jsonConfig.Nodes = make(map[string]WindowConfig)
	}

	return jsonConfig
}

func (jc Config) Update(wc WindowConfig) {
	jc.Nodes[wc.Mark] = wc
}

func createJsonConfigFile(configFileLoc string) {
	var jsonConfig Config
	// default value
	// jsonConfig.StatusBarHeight = 29

	file, err := os.Create(configFileLoc)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(jsonConfig); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
}

func (jc Config) Dump() {
	jsonData, err := json.MarshalIndent(jc, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(jc.Path, jsonData, defaultPerms); err != nil {
		log.Fatal(err)
	}
}
