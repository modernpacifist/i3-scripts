package scale_float_container

import (
	"errors"
	"fmt"

	common "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
	"go.i3wm.org/i3/v4"
)

const (
	shit_StatusBarHeight = 35
)

func increaseHeightToTop(value int64) error {
	cmd := fmt.Sprintf("resize grow height %d px, move container up %d px", value, value)
	if err := common.RunI3Command(cmd); err != nil {
		return err
	}

	return nil
}

func increaseHeightToBottom(value int64) error {
	cmd := fmt.Sprintf("resize grow height %d px", value)
	if err := common.RunI3Command(cmd); err != nil {
		return err
	}

	return nil
}

func increaseWidthToLeft(value int64) error {
	cmd := fmt.Sprintf("resize grow width %d px, move container left %d px", value, value)
	if err := common.RunI3Command(cmd); err != nil {
		return err
	}

	return nil
}

func increaseWidthToRight(value int64) error {
	cmd := fmt.Sprintf("resize grow width %d px", value)
	if err := common.RunI3Command(cmd); err != nil {
		return err
	}

	return nil
}

func getScreenMargins(output i3.Output, node i3.Node) (int64, int64, int64, int64) {
	outputRect := output.Rect
	nodeRect := node.Rect

	distanceLeft := nodeRect.X - outputRect.X
	distanceRight := outputRect.X + outputRect.Width - (nodeRect.X + nodeRect.Width)
	distanceTop := nodeRect.Y - shit_StatusBarHeight
	distanceBottom := outputRect.Y + outputRect.Height - (nodeRect.Y + nodeRect.Height)

	return distanceLeft, distanceRight, distanceTop, distanceBottom
}

func Execute(arg string) error {
	focusedOutput, err := common.GetFocusedOutput()
	if err != nil {
		return err
	}
	fmt.Printf("\n%+v\n", focusedOutput.Rect)

	focusedNode, err := common.GetFocusedNode()
	if err != nil {
		return err
	}

	if focusedNode.Floating != "user_on" && focusedNode.Floating != "auto_on" {
		return errors.New("node is not floating")
	}

	fmt.Printf("\nNode Rect: %+v\n", focusedNode.Rect)
	fmt.Printf("\nWindow Rect: %+v\n\n", focusedNode.WindowRect)

	distanceLeft, distanceRight, distanceTop, distanceBottom := getScreenMargins(focusedOutput, focusedNode)
	fmt.Printf("\nDistance to left edge: %d\n", distanceLeft)
	fmt.Printf("\nDistance to right edge: %d\n", distanceRight)
	fmt.Printf("\nDistance to top edge: %d\n", distanceTop)
	fmt.Printf("\nDistance to bottom edge: %d\n\n", distanceBottom)

	switch arg {
	case "w":
		increaseHeightToTop(distanceTop)
	case "s":
		increaseHeightToBottom(distanceBottom)
	case "d":
		increaseWidthToRight(distanceRight)
	case "a":
		increaseWidthToLeft(distanceLeft)
	}

	return nil
}
