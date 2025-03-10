package i3operations

import (
	"errors"
	"fmt"

	"github.com/modernpacifist/i3-scripts-go/internal/i3operations"
)

const (
	shittemp_StatusBarOffset int64 = -53

	topLeft   = 1
	topMiddle = 2
	topRight  = 3

	middleLeft   = 4
	middleMiddle = 5
	middleRight  = 6

	bottomLeft   = 7
	bottomMiddle = 8
	bottomRight  = 9
)

type Position struct {
	X int64
	Y int64
}

type OutputGeometry struct {
	Width  int64
	Height int64
}

type NodeGeometry struct {
	X      int64
	Y      int64
	Width  int64
	Height int64
}

func getFocusedOutputGeometry() (OutputGeometry, error) {
	focusedOutput, err := i3operations.GetFocusedOutput()
	if err != nil {
		return OutputGeometry{}, err
	}
	return OutputGeometry{
		Width:  focusedOutput.Rect.Width,
		Height: focusedOutput.Rect.Height,
	}, nil
}

func resolveNewPosition(focusedOutputGeometry OutputGeometry, focusedNodeGeometry NodeGeometry) (int64, int64, error) {
	dummyInput := 1

	positions := map[int]Position{
		topLeft: {
			X: focusedNodeGeometry.X,
			Y: focusedNodeGeometry.Y,
		},
		topMiddle: {
			X: focusedOutputGeometry.Width / 2,
			Y: 0,
		},
		topRight: {
			X: focusedOutputGeometry.Width,
			Y: 0,
		},
		middleLeft: {
			X: 0,
			Y: focusedOutputGeometry.Height / 2,
		},
		middleMiddle: {
			X: focusedOutputGeometry.Width / 2,
			Y: focusedOutputGeometry.Height / 2,
		},
		middleRight: {
			X: focusedOutputGeometry.Width,
			Y: focusedOutputGeometry.Height / 2,
		},
		bottomLeft: {
			X: 0,
			Y: focusedOutputGeometry.Height,
		},
		bottomMiddle: {
			X: focusedOutputGeometry.Width / 2,
			Y: focusedOutputGeometry.Height,
		},
		bottomRight: {
			X: focusedOutputGeometry.Width,
			Y: focusedOutputGeometry.Height,
		},
	}

	pos := positions[dummyInput]

	return pos.X, pos.Y, nil
}

func Execute(arg int) error {
	focusedOutputGeometry, err := getFocusedOutputGeometry()
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", focusedOutputGeometry)

	focusedNode, err := i3operations.GetFocusedNode()
	if err != nil {
		return err
	}

	// possibly check for having 'on' in the focusedNode.Floating as substring
	if focusedNode.Floating != "user_on" {
		return errors.New("focused node is not floating")
	}

	fmt.Printf("%+v\n", focusedNode.Window)

	focusedNodeGeometry := NodeGeometry{
		X:      focusedNode.Rect.X,
		Y:      focusedNode.Rect.Y,
		Width:  focusedNode.Rect.Width,
		Height: focusedNode.Rect.Height,
	}
	fmt.Printf("%+v\n", focusedNodeGeometry)

	newX, newY, err := resolveNewPosition(focusedOutputGeometry, focusedNodeGeometry)
	if err != nil {
		return err
	}

	fmt.Printf("newX: %d, newY: %d\n", newX, newY)

	if err := i3operations.MoveNodeToPosition(focusedNode.Window, newX, newY); err != nil {
		return err
	}

	return nil
}
