package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/modernpacifist/i3-scripts-go/config"
	"go.i3wm.org/i3/v4"
)

const (
	configFilename string = ".ScaleFloatWindow.json"
)

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
	val := wc.Y - config.CONFIG.StatusBarHeight
	if val == 0 {
		wc.ResizedPlusYFlag = true
		val = -wc.PreviousPlusYValue
	}

	// TODO: in case of reverse resize here must be a check for resized flag <29-10-23, modernpacifist> //
	i3.RunCommand(fmt.Sprintf("resize grow height %d px, move container up %d px", val, val))

	wc.Y -= val
	if wc.Y == config.CONFIG.StatusBarHeight {
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
	node, exists := config.CONFIG.Windows[wc.Mark]
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

func init() {
	configFileLoc := resolveJsonAbsolutePath(configFilename)
	config.CONFIG = JsonConfigConstructor(configFileLoc)

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

	config.CONFIG.Update(focusedWindow)
	config.CONFIG.Dump()
}
