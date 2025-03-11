package i3operations

import (
	"fmt"
	"os/exec"

	"github.com/modernpacifist/i3-scripts-go/internal/i3operations"
	"go.i3wm.org/i3/v4"
)

// "Temporary" constant
const (
	shittemp_StatusBarOffset int64 = 26
)

// Position constants
const (
	topLeft      uint8 = 1
	topMiddle    uint8 = 2
	topRight     uint8 = 3
	middleLeft   uint8 = 4
	middleMiddle uint8 = 5
	middleRight  uint8 = 6
	bottomLeft   uint8 = 7
	bottomMiddle uint8 = 8
	bottomRight  uint8 = 9
)

type outputGeometry struct {
	Width        int64
	Height       int64
	WidthOffset  int64
	HeightOffset int64
}

func outputGeometryConstructor(output i3.Output) outputGeometry {
	return outputGeometry{
		Width:        output.Rect.Width,
		Height:       output.Rect.Height,
		WidthOffset:  output.Rect.X,
		HeightOffset: output.Rect.Y,
	}
}

type nodeGeometry struct {
	X           int64
	Y           int64
	Width       int64
	Height      int64
	BorderWidth int64
}

func nodeGeometryConstructor(node i3.Node) nodeGeometry {
	return nodeGeometry{
		X:           node.Rect.X,
		Y:           node.Rect.Y,
		Width:       node.Rect.Width,
		Height:      node.Rect.Height,
		BorderWidth: node.CurrentBorderWidth,
	}
}

type Position struct {
	X int64
	Y int64
}

func resolveNewPosition(outputGeometry outputGeometry, nodeGeometry nodeGeometry) (int64, int64, error) {
	// move to bottom right
	// dummyInput := bottomRight
	dummyInput := bottomMiddle

	positions := map[uint8]Position{
		topLeft: {
			X: nodeGeometry.X,
			Y: nodeGeometry.Y,
		},
		topMiddle: {
			X: outputGeometry.Width / 2,
			Y: 0,
		},
		topRight: {
			X: outputGeometry.Width,
			Y: 0,
		},
		middleLeft: {
			X: 0,
			Y: outputGeometry.Height / 2,
		},
		middleMiddle: {
			X: outputGeometry.Width / 2,
			Y: outputGeometry.Height / 2,
		},
		middleRight: {
			X: outputGeometry.Width,
			Y: outputGeometry.Height / 2,
		},
		bottomLeft: {
			X: outputGeometry.WidthOffset + nodeGeometry.BorderWidth,
			Y: outputGeometry.Height + outputGeometry.HeightOffset - nodeGeometry.Height + nodeGeometry.BorderWidth + shittemp_StatusBarOffset,
		},
		bottomMiddle: {
			X: outputGeometry.WidthOffset + outputGeometry.Width/2 - nodeGeometry.Width/2 + nodeGeometry.BorderWidth,
			Y: outputGeometry.Height + outputGeometry.HeightOffset - nodeGeometry.Height + nodeGeometry.BorderWidth + shittemp_StatusBarOffset,
		},
		bottomRight: {
			X: outputGeometry.Width + outputGeometry.WidthOffset - nodeGeometry.Width + nodeGeometry.BorderWidth,
			Y: outputGeometry.Height + outputGeometry.HeightOffset - nodeGeometry.Height + nodeGeometry.BorderWidth + shittemp_StatusBarOffset,
		},
	}

	pos := positions[dummyInput]

	return pos.X, pos.Y, nil
}

func moveNodeToPosition(nodeId, x, y int64) error {
	cmd := fmt.Sprintf("xdotool windowmove %d %d %d", nodeId, x, y)
	if _, err := exec.Command("bash", "-c", cmd).Output(); err != nil {
		return err
	}

	return nil
}

func Execute(arg int) error {
	focusedOutput, err := i3operations.GetFocusedOutput()
	if err != nil {
		return err
	}
	focusedOutputGeometry := outputGeometryConstructor(focusedOutput)
	fmt.Printf("%+v\n", focusedOutputGeometry)

	focusedNode, err := i3operations.GetFocusedNode()
	if err != nil {
		return err
	}
	focusedNodeGeometry := nodeGeometryConstructor(focusedNode)
	fmt.Printf("%+v\n", focusedNodeGeometry)

	newX, newY, err := resolveNewPosition(focusedOutputGeometry, focusedNodeGeometry)
	if err != nil {
		return err
	}
	fmt.Printf("newX: %d, newY: %d\n", newX, newY)

	if err := moveNodeToPosition(focusedNode.Window, newX, newY); err != nil {
		return err
	}

	return nil
}
