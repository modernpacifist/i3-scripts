package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"

	"go.i3wm.org/i3/v4"
)

const (
	ConfigFilename string = ".ScaleFloatWindow.json"
)

var globalConfig JsonConfig
var globalMonitorDimensions i3.Rect

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

func getNodeMark(node *i3.Node) string {
	if len(node.Marks) == 0 {
		return ""
	}
	return node.Marks[0]
}

func resolveResizedFlags(node *i3.Node, flagname string) bool {
	switch flagname {
	case "plusY":
		return node.Rect.Y == globalConfig.StatusBarHeight
	case "minusY":
		return node.Rect.Height == globalMonitorDimensions.Height-node.Rect.Y
	case "plusX":
		return node.Rect.Width == globalMonitorDimensions.Width-node.Rect.X
	case "minusX":
		return node.Rect.X == 0
	}
	return false
}

func getPreviousResizeValues(node *i3.Node) map[string]int64 {
	resMap := make(map[string]int64)

	for _, n := range globalConfig.Windows {
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
	StatusBarHeight int64                   `json:"StatusBarHeight"`
	Windows         map[string]WindowConfig `json:"Windows"`
}

func createJsonConfigFile(configFileLoc string) {
	var jsonConfig JsonConfig
	// default value
	jsonConfig.StatusBarHeight = 29

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

func (jc *JsonConfig) Update(wc WindowConfig) {
	jc.Windows[wc.Mark] = wc
}

func (jc *JsonConfig) Dump() {
	jsonData, err := json.MarshalIndent(jc, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(jc.Location, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
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

func (wc *WindowConfig) resizePlusY() {
	val := wc.Y - globalConfig.StatusBarHeight
	if val == 0 {
		wc.ResizedPlusYFlag = true
		val = -wc.PreviousPlusYValue
	}

	// TODO: in case of reverse resize here must be a check for resized flag <29-10-23, modernpacifist> //
	i3.RunCommand(fmt.Sprintf("resize grow height %d px, move container up %d px", val, val))

	wc.Y -= val
	if wc.Y == globalConfig.StatusBarHeight {
		wc.ResizedPlusYFlag = true
	} else {
		wc.ResizedPlusYFlag = false
	}

	// TODO: need to check if resulting resize will be out of bounds <29-10-23, modernpacifist> //
	// if the window exceeds the boundaries of focused output, then do not update this value
	wc.PreviousPlusYValue = val
}

func (wc *WindowConfig) resizeMinusY() {
	val := globalMonitorDimensions.Height - wc.Height - wc.Y
	if val == 0 {
		wc.ResizedMinusYFlag = true
		val = -wc.PreviousMinusYValue
	}

	i3.RunCommand(fmt.Sprintf("resize grow height %d px", val))

	// TODO: in case of reverse resize here must be a check for resized flag <29-10-23, modernpacifist> //
	wc.Height += val
	if wc.Height == globalMonitorDimensions.Height-wc.Y {
		wc.ResizedMinusYFlag = true
	} else {
		wc.ResizedMinusYFlag = false
	}

	wc.PreviousMinusYValue = val
}

func (wc *WindowConfig) resizePlusX() {
	val := globalMonitorDimensions.Width - wc.X - wc.Width
	if val == 0 {
		wc.ResizedPlusXFlag = true
		// TODO: maybe calculate here <29-10-23, modernpacifist> //
		val = -wc.PreviousPlusXValue
	}

	i3.RunCommand(fmt.Sprintf("resize grow width %d px", val))

	wc.Width += val
	if wc.Width == globalMonitorDimensions.Width-wc.X {
		wc.ResizedPlusXFlag = true
	} else {
		wc.ResizedPlusXFlag = false
	}

	wc.PreviousPlusXValue = val
}

func (wc *WindowConfig) resizeMinusX() {
	val := wc.X
	if val == 0 {
		wc.ResizedMinusXFlag = true
		val = -wc.PreviousMinusXValue
	}

	i3.RunCommand(fmt.Sprintf("resize grow width %d px, move container left %d px", val, val))

	wc.X -= val
	if wc.X == 0 {
		wc.ResizedMinusXFlag = true
	} else {
		wc.ResizedMinusXFlag = false
	}

	wc.PreviousMinusXValue = val
}

func (wc *WindowConfig) resetGeometry() {
	node, exists := globalConfig.Windows[wc.Mark]
	if exists == false {
		os.Exit(0)
	}

	cmd := fmt.Sprintf("move absolute position %d %d, resize set %d %d", node.X, node.Y, node.Width, node.Height)
	i3.RunCommand(cmd)
}

func (wc *WindowConfig) resizeHorizontally(val int) {
	i3.RunCommand(fmt.Sprintf("resize grow width %d px, move container left %d px", val, val/2))
}

func getPrimaryOutputRect() i3.Rect {
	outputs, err := i3.GetOutputs()
	if err != nil {
		log.Fatal(err)
	}

	for _, output := range outputs {
		if output.Primary && output.Active {
			fmt.Println(output.Rect)
			return output.Rect
		}
	}

	return i3.Rect{}
}

func getI3Tree() i3.Tree {
	tree, err := i3.GetTree()
	if err != nil {
		log.Fatal(err)
	}

	return tree
}

func resolveJsonAbsolutePath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Replace(fmt.Sprintf("~/%s", filename), "~", usr.HomeDir, 1)
}

func init() {
	configFileLoc := resolveJsonAbsolutePath(ConfigFilename)
	globalConfig = JsonConfigConstructor(configFileLoc)

	// TODO: this must not be like this <28-10-23, modernpacifist> //
	globalMonitorDimensions = getPrimaryOutputRect()
}

func main() {
	mode := flag.String("mode", "", "Specify resize mode")
	widen := flag.Int("widen", 0, "Specify resize widen")
	flag.Parse()

	if *mode == "" && *widen == 0 {
		log.Fatal("No arguments were specified")
	}

	focusedNode := getFocusedNode()
	focusedWindow := WindowConfigConstructor(focusedNode)

	switch *mode {
	case "w":
		focusedWindow.resizePlusY()
	case "s":
		focusedWindow.resizeMinusY()
	case "d":
		focusedWindow.resizePlusX()
	case "a":
		focusedWindow.resizeMinusX()
	// this flag for debug purposes, remove later
	case "f":

	// TODO: reset the size of the floating window to its default <29-10-23, modernpacifist> //
	case "r":
		focusedWindow.resetGeometry()
		//default:
		//log.Fatal(errors.New("The -mode flag must be only of w/s/d/a values"))
	}

	// TODO: this also must be in switch<15-11-23, modernpacifist> //
	if *widen != 0 {
		focusedWindow.resizeHorizontally(*widen)
	}

	globalConfig.Update(focusedWindow)
	globalConfig.Dump()
}
