package margin_resize

import (
	"errors"
	"fmt"

	config "github.com/modernpacifist/i3-scripts-go/internal/config/margin_resize"
	common "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
	"go.i3wm.org/i3/v4"
)

const (
	defaultStatusBarHeight = 35
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
	distanceTop := nodeRect.Y - defaultStatusBarHeight
	distanceBottom := outputRect.Y + outputRect.Height - (nodeRect.Y + nodeRect.Height)

	return distanceLeft, distanceRight, distanceTop, distanceBottom
}

func normalizeResizeValue(resizeValue int64, pastNode config.NodeConfig, focusedNode i3.Node) int64 {
	if resizeValue == 0 && pastNode.Node.Rect.Height > 0 {
		resizeValue = pastNode.Node.Rect.Height - focusedNode.Rect.Height
	}
	return resizeValue
}

func Execute(arg string) error {
	focusedOutput, err := common.GetFocusedOutput()
	if err != nil {
		return err
	}

	focusedNode, err := common.GetFocusedNode()
	if err != nil {
		return err
	}

	// double check this later
	if focusedNode.Floating != "user_on" && focusedNode.Floating != "auto_on" {
		return errors.New("node is not floating")
	}

	focusedNodeConfigIdentifier := common.GetNodeMark(focusedNode)
	if focusedNodeConfigIdentifier == "" {
		focusedNodeConfigIdentifier = fmt.Sprintf("%d", focusedNode.Window)
	}

	conf, err := config.Create()
	if err != nil {
		return err
	}

	pastNode, exists := conf.Nodes[focusedNodeConfigIdentifier]
	if !exists {
		conf.Nodes[focusedNodeConfigIdentifier] = config.NodeConfig{
			Node: focusedNode,
		}
		pastNode = conf.Nodes[focusedNodeConfigIdentifier]
	}

	distanceLeft, distanceRight, distanceTop, distanceBottom := getScreenMargins(focusedOutput, focusedNode)

	var resizeValue int64
	switch arg {
	case "top":
		resizeValue = normalizeResizeValue(distanceTop, pastNode, focusedNode)
		increaseHeightToTop(resizeValue)
	case "bottom":
		resizeValue = normalizeResizeValue(distanceBottom, pastNode, focusedNode)
		increaseHeightToBottom(resizeValue)
	case "right":
		resizeValue = normalizeResizeValue(distanceRight, pastNode, focusedNode)
		increaseWidthToRight(resizeValue)
	case "left":
		resizeValue = normalizeResizeValue(distanceLeft, pastNode, focusedNode)
		increaseWidthToLeft(resizeValue)
	default:
		return errors.New("invalid argument")
	}

	if err := conf.Dump(); err != nil {
		return err
	}

	return nil
}
