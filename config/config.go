package config

import (
	"encoding/json"
	"fmt"
	"os"

	"go.i3wm.org/i3/v4"
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
func WindowConfigConstructor(node *i3.Node) config.WindowConfig {
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
