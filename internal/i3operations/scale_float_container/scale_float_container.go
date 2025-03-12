package scale_float_container

import (
	"fmt"
	"os"

	config "github.com/modernpacifist/i3-scripts-go/internal/config/scale_float_container"
	"github.com/modernpacifist/i3-scripts-go/internal/i3operations"
	"go.i3wm.org/i3/v4"
)

// func (wc config.NodeConfig) resizePlusY() {
// 	val := wc.Y - config.CONFIG.StatusBarHeight
// 	if val == 0 {
// 		wc.ResizedPlusYFlag = true
// 		val = -wc.PreviousPlusYValue
// 	}

// 	// TODO: in case of reverse resize here must be a check for resized flag <29-10-23, modernpacifist> //
// 	i3.RunCommand(fmt.Sprintf("resize grow height %d px, move container up %d px", val, val))

// 	wc.Y -= val
// 	if wc.Y == configs.CONFIG.StatusBarHeight {
// 		wc.ResizedPlusYFlag = true
// 	} else {
// 		wc.ResizedPlusYFlag = false
// 	}

// 	// TODO: need to check if resulting resize will be out of bounds <29-10-23, modernpacifist> //
// 	// if the window exceeds the boundaries of focused output, then do not update this value
// 	wc.PreviousPlusYValue = val
// }

// func (wc config.NodeConfig) resizeMinusY() {
// 	val := globalMonitorDimensions.Height - wc.Height - wc.Y
// 	if val == 0 {
// 		wc.ResizedMinusYFlag = true
// 		val = -wc.PreviousMinusYValue
// 	}

// 	i3.RunCommand(fmt.Sprintf("resize grow height %d px", val))

// 	// TODO: in case of reverse resize here must be a check for resized flag <29-10-23, modernpacifist> //
// 	wc.Height += val
// 	if wc.Height == globalMonitorDimensions.Height-wc.Y {
// 		wc.ResizedMinusYFlag = true
// 	} else {
// 		wc.ResizedMinusYFlag = false
// 	}

// 	wc.PreviousMinusYValue = val
// }

// func (wc *WindowConfig) resizePlusX() {
// 	val := globalMonitorDimensions.Width - wc.X - wc.Width
// 	if val == 0 {
// 		wc.ResizedPlusXFlag = true
// 		// TODO: maybe calculate here <29-10-23, modernpacifist> //
// 		val = -wc.PreviousPlusXValue
// 	}

// 	i3.RunCommand(fmt.Sprintf("resize grow width %d px", val))

// 	wc.Width += val
// 	if wc.Width == globalMonitorDimensions.Width-wc.X {
// 		wc.ResizedPlusXFlag = true
// 	} else {
// 		wc.ResizedPlusXFlag = false
// 	}

// 	wc.PreviousPlusXValue = val
// }

// func (wc *WindowConfig) resizeMinusX() {
// 	val := wc.X
// 	if val == 0 {
// 		wc.ResizedMinusXFlag = true
// 		val = -wc.PreviousMinusXValue
// 	}

// 	_, err := i3.RunCommand(fmt.Sprintf("resize grow width %d px, move container left %d px", val, val))
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 		return
// 	}

// 	wc.X -= val
// 	if wc.X == 0 {
// 		wc.ResizedMinusXFlag = true
// 	} else {
// 		wc.ResizedMinusXFlag = false
// 	}

// 	wc.PreviousMinusXValue = val
// }

// func (wc *WindowConfig) resetGeometry() {
// 	node, exists := config.CONFIG.Windows[wc.Mark]
// 	if exists == false {
// 		os.Exit(0)
// 	}

// 	cmd := fmt.Sprintf("move absolute position %d %d, resize set %d %d", node.X, node.Y, node.Width, node.Height)
// 	i3.RunCommand(cmd)
// }

// func (wc *WindowConfig) resizeHorizontally(val int) {
// 	i3.RunCommand(fmt.Sprintf("resize grow width %d px, move container left %d px", val, val/2))
// }

func getPreviousResizeValues(node i3.Node, config config.Config) map[string]int64 {
	resMap := make(map[string]int64)

	// do not need linear search in this, simplifying the config
	for _, n := range config.Nodes {
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

func resolveResizedFlags(node i3.Node, output i3.Output, flagname string) bool {
	switch flagname {
	case "plusY":
		// don't believe I wrote this
		return node.Rect.Y == 64
	case "minusY":
		return node.Rect.Height == output.Rect.Height-node.Rect.Y
	case "plusX":
		return node.Rect.Width == output.Rect.Width-node.Rect.X
	case "minusX":
		return node.Rect.X == 0
	}

	return false
}

func NodeConfigConstructor(node i3.Node, output i3.Output, c config.Config) config.NodeConfig {
	nodeID := node.Window

	previousResizeValuesMap := getPreviousResizeValues(node, c)
	plusY, _ := previousResizeValuesMap["plusY"]
	minusY, _ := previousResizeValuesMap["minusY"]
	plusX, _ := previousResizeValuesMap["plusX"]
	minusX, _ := previousResizeValuesMap["minusX"]

	// TODO: is the node does not contain a mark, just use a container id <17-11-23, modernpacifist> //
	// nodeMark := getNodeMark(node)
	// if nodeMark == "" {
	// 	nodeMark = strconv.FormatInt(node.Window, 10)
	// }

	return config.NodeConfig{
		ID:                  nodeID,
		ResizedPlusYFlag:    resolveResizedFlags(node, output, "plusY"),
		ResizedMinusYFlag:   resolveResizedFlags(node, output, "minusY"),
		ResizedPlusXFlag:    resolveResizedFlags(node, output, "plusX"),
		ResizedMinusXFlag:   resolveResizedFlags(node, output, "minusX"),
		Marks:               i3operations.GetNodeMark(node),
		X:                   node.Rect.X,
		Y:                   node.Rect.Y,
		Width:               node.Rect.Width,
		Height:              node.Rect.Height,
		PreviousPlusYValue:  plusY,
		PreviousMinusYValue: minusY,
		PreviousPlusXValue:  plusX,
		PreviousMinusXValue: minusX,
	}
}

func Execute(arg string) error {
	focusedOutput, err := i3operations.GetFocusedOutput()
	if err != nil {
		return err
	}

	configInstance, err := config.Create(focusedOutput)
	if err != nil {
		return err
	}

	focusedNode, err := i3operations.GetFocusedNode()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", configInstance)
	fmt.Printf("%+v\n", focusedNode)

	configInstance.Nodes[focusedNode.Window] = NodeConfigConstructor(focusedNode, focusedOutput, configInstance)

	if err := configInstance.Dump(); err != nil {
		return err
	}

	os.Exit(0)

	// focusedWindow := config.WindowConfigConstructor(focusedNode)

	// switch arg {
	// case "w":
	// 	focusedWindow.resizePlusY()
	// case "s":
	// 	focusedWindow.resizeMinusY()
	// case "d":
	// 	focusedWindow.resizePlusX()
	// case "a":
	// 	focusedWindow.resizeMinusX()
	// // this flag for debug purposes, remove later
	// case "f":
	// // TODO: reset the size of the floating window to its default <29-10-23, modernpacifist> //
	// case "r":
	// 	focusedWindow.resetGeometry()
	// }

	// // TODO: this also must be in switch<15-11-23, modernpacifist> //
	// if *widen != 0 {
	// 	focusedWindow.resizeHorizontally(*widen)
	// }

	// config.CONFIG.Update(focusedWindow)
	// config.CONFIG.Dump()
	if err := configInstance.Dump(); err != nil {
		return err
	}

	return nil
}
