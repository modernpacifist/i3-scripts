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

func resolveNewPosition(oGeometry outputGeometry, nGeometry nodeGeometry) (int64, int64, error) {
	// move to bottom right
	dummyInput := bottomRight

	positions := map[uint8]Position{
		topLeft: {
			X: nGeometry.X,
			Y: nGeometry.Y,
		},
		topMiddle: {
			X: oGeometry.Width / 2,
			Y: 0,
		},
		topRight: {
			X: oGeometry.Width,
			Y: 0,
		},
		middleLeft: {
			X: 0,
			Y: oGeometry.Height / 2,
		},
		middleMiddle: {
			X: oGeometry.Width / 2,
			Y: oGeometry.Height / 2,
		},
		middleRight: {
			X: oGeometry.Width,
			Y: oGeometry.Height / 2,
		},
		bottomLeft: {
			X: 0,
			Y: oGeometry.Height,
		},
		bottomMiddle: {
			X: oGeometry.Width / 2,
			Y: oGeometry.Height,
		},
		bottomRight: {
			X: oGeometry.Width + oGeometry.WidthOffset - nGeometry.Width + nGeometry.BorderWidth,
			Y: oGeometry.Height + oGeometry.HeightOffset - nGeometry.Height + nGeometry.BorderWidth + shittemp_StatusBarOffset,
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
