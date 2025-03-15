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

func normalizeResizeValue(direction string, resizeValue int64, output i3.Output, pastNode config.NodeConfig) int64 {
	if resizeValue == 0 {
		switch direction {
		case "top":
			resizeValue = output.Rect.Height - pastNode.Node.Rect.Height
		case "bottom":
			resizeValue = output.Rect.Height - pastNode.Node.Rect.Height
		case "right":
			resizeValue = output.Rect.Width - pastNode.Node.Rect.Width
		case "left":
			resizeValue = output.Rect.Width - pastNode.Node.Rect.Width
		}
		return -resizeValue
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

	// load past config into memory
	pastNodeConfig, exists := conf.Nodes[focusedNodeConfigIdentifier]
	if !exists {
		conf.Nodes[focusedNodeConfigIdentifier] = config.NodeConfig{
			Node: focusedNode,
		}
		pastNodeConfig = conf.Nodes[focusedNodeConfigIdentifier]
	}

	// instantly update config file with new data
	conf.Nodes[focusedNodeConfigIdentifier] = config.NodeConfig{
		Node: focusedNode,
	}

	if err := conf.Dump(); err != nil {
		return err
	}

	distanceLeft, distanceRight, distanceTop, distanceBottom := getScreenMargins(focusedOutput, focusedNode)

	var resizeValue int64
	switch arg {
	case "top":
		resizeValue = normalizeResizeValue("top", distanceTop, focusedOutput, pastNodeConfig)
		increaseHeightToTop(resizeValue)
	case "bottom":
		resizeValue = normalizeResizeValue("bottom", distanceBottom, focusedOutput, pastNodeConfig)
		increaseHeightToBottom(resizeValue)
	case "right":
		resizeValue = normalizeResizeValue("right", distanceRight, focusedOutput, pastNodeConfig)
		increaseWidthToRight(resizeValue)
	case "left":
		resizeValue = normalizeResizeValue("left", distanceLeft, focusedOutput, pastNodeConfig)
		increaseWidthToLeft(resizeValue)
	default:
		return errors.New("invalid argument")
	}

	return nil
}
