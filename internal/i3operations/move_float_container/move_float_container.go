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

func resolveNewPosition(dummyInput uint8, outputGeometry outputGeometry, nodeGeometry nodeGeometry) (Position, error) {
	topX := outputGeometry.WidthOffset + nodeGeometry.BorderWidth
	topY := nodeGeometry.BorderWidth + 61
	middleX := outputGeometry.WidthOffset +
		outputGeometry.Width/2 -
		nodeGeometry.Width/2 +
		nodeGeometry.BorderWidth
	middleY := outputGeometry.Height/2 -
		nodeGeometry.Height/2 +
		nodeGeometry.BorderWidth +
		shittemp_StatusBarOffset
	bottomX := outputGeometry.Width +
		outputGeometry.WidthOffset -
		nodeGeometry.Width +
		nodeGeometry.BorderWidth
	bottomY := outputGeometry.Height +
		outputGeometry.HeightOffset -
		nodeGeometry.Height +
		nodeGeometry.BorderWidth +
		shittemp_StatusBarOffset

	positions := map[uint8]Position{
		topLeft: {
			X: topX,
			Y: topY,
		},
		topMiddle: {
			X: middleX,
			Y: topY,
		},
		topRight: {
			X: bottomX,
			Y: topY,
		},
		middleLeft: {
			X: topX,
			Y: middleY,
		},
		middleMiddle: {
			X: middleX,
			Y: middleY,
		},
		middleRight: {
			X: bottomX,
			Y: middleY,
		},
		bottomLeft: {
			X: topX,
			Y: bottomY,
		},
		bottomMiddle: {
			X: middleX,
			Y: bottomY,
		},
		bottomRight: {
			X: bottomX,
			Y: bottomY,
		},
	}

	if _, ok := positions[dummyInput]; ok {
		return positions[dummyInput], nil
	}

	return Position{}, fmt.Errorf("invalid input: %d", dummyInput)
}

func moveNodeToPosition(nodeId int64, position Position) error {
	cmd := fmt.Sprintf("xdotool windowmove %d %d %d", nodeId, position.X, position.Y)
	if _, err := exec.Command("bash", "-c", cmd).Output(); err != nil {
		return err
	}

	return nil
}

// story: possibility to move to different output based on the flag like output [DP-1/HDMI-1]
func Execute(arg uint8) error {
	focusedOutput, err := i3operations.GetFocusedOutput()
	if err != nil {
		return err
	}
	focusedOutputGeometry := outputGeometryConstructor(focusedOutput)

	focusedNode, err := i3operations.GetFocusedNode()
	if err != nil {
		return err
	}
	focusedNodeGeometry := nodeGeometryConstructor(focusedNode)

	newPosition, err := resolveNewPosition(arg, focusedOutputGeometry, focusedNodeGeometry)
	if err != nil {
		return err
	}

	if err := moveNodeToPosition(focusedNode.Window, newPosition); err != nil {
		return err
	}

	return nil
}
