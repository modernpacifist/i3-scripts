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

	// Direction constants
	DirectionTop    = "top"
	DirectionBottom = "bottom"
	DirectionRight  = "right"
	DirectionLeft   = "left"
)

// ResizeStrategy defines the interface for different resize operations
type ResizeStrategy interface {
	Resize(value int64) error
	// calculateResizeValue(output i3.Output, nodeConf config.NodeConfig) int64
}

// TopResizeStrategy implements ResizeStrategy for top direction
type TopResizeStrategy struct{}

func (s *TopResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow height %d px, move container up %d px", value, value))
}

// func (s *TopResizeStrategy) calculateResizeValue(output i3.Output, nodeConf config.NodeConfig) int64 {
// 	// if val == 0 && nodeConf.Node.Rect.Y > output.Rect.Y {
// 	if nodeConf.Node.Rect.Y-defaultStatusBarHeight >= output.Rect.Y {
// 		return nodeConf.Node.Rect.Y - defaultStatusBarHeight
// 	}
// 	return -(nodeConf.Node.Rect.Y - defaultStatusBarHeight)
// }

type BottomResizeStrategy struct{}

func (s *BottomResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow height %d px", value))
}

// func (s *BottomResizeStrategy) calculateResizeValue(output i3.Output, nodeConf config.NodeConfig) int64 {
// 	if nodeConf.Node.Rect.Y+nodeConf.Node.Rect.Height == output.Rect.Height {
// 		return -(output.Rect.Height - nodeConf.Node.Rect.Height)
// 	}
// 	return -(nodeConf.Node.Rect.Y + nodeConf.Node.Rect.Height - output.Rect.Height)
// }

type RightResizeStrategy struct{}

func (s *RightResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow width %d px", value))
}

// func (s *RightResizeStrategy) calculateResizeValue(output i3.Output, nodeConf config.NodeConfig) int64 {
// 	if nodeConf.Node.Rect.X+nodeConf.Node.Rect.Width == output.Rect.Width {
// 		return -(output.Rect.Width - nodeConf.Node.Rect.Width)
// 	}
// 	return -(nodeConf.Node.Rect.X + nodeConf.Node.Rect.Width - output.Rect.Width)
// }

type LeftResizeStrategy struct{}

func (s *LeftResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow width %d px, move container left %d px", value, value))
}

// func (s *LeftResizeStrategy) calculateResizeValue(output i3.Output, nodeConf config.NodeConfig) int64 {
// 	if nodeConf.Node.Rect.X == output.Rect.X {
// 		return -(output.Rect.Width - nodeConf.Node.Rect.Width)
// 	}
// 	return -(nodeConf.Node.Rect.X - output.Rect.X)
// }

type ResizeContext struct {
	strategy ResizeStrategy
}

func NewResizeContext(direction string) (*ResizeContext, error) {
	var strategy ResizeStrategy
	switch direction {
	case DirectionTop:
		strategy = &TopResizeStrategy{}
	case DirectionBottom:
		strategy = &BottomResizeStrategy{}
	case DirectionRight:
		strategy = &RightResizeStrategy{}
	case DirectionLeft:
		strategy = &LeftResizeStrategy{}
	default:
		return nil, errors.New("invalid direction")
	}
	return &ResizeContext{strategy: strategy}, nil
}

func getScreenMargins(output i3.Output, node i3.Node) (int64, int64, int64, int64) {
	outputRect := output.Rect
	nodeRect := node.Rect

	distanceToTop := nodeRect.Y - defaultStatusBarHeight
	distanceToBottom := outputRect.Y + outputRect.Height - (nodeRect.Y + nodeRect.Height)
	distanceToRight := outputRect.Width - nodeRect.Width
	distanceToLeft := nodeRect.X - outputRect.X

	return distanceToTop, distanceToBottom, distanceToRight, distanceToLeft
}

func Execute(resize_direction string) error {
	focusedOutput, err := common.GetFocusedOutput()
	if err != nil {
		return err
	}

	focusedNode, err := common.GetFocusedNode()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", focusedNode.Rect)

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
	loadedNodeConf, exists := conf.Nodes[focusedNodeConfigIdentifier]
	if !exists {
		conf.Nodes[focusedNodeConfigIdentifier] = config.NodeConfig{
			Node: focusedNode,
		}
		loadedNodeConf = conf.Nodes[focusedNodeConfigIdentifier]
	}

	distanceToTop, distanceToBottom, distanceToRight, distanceToLeft := getScreenMargins(focusedOutput, focusedNode)

	currentNodeConf := config.NodeConfig{
		DistanceToTop:    loadedNodeConf.DistanceToTop,
		DistanceToBottom: loadedNodeConf.DistanceToBottom,
		DistanceToRight:  loadedNodeConf.DistanceToRight,
		DistanceToLeft:   loadedNodeConf.DistanceToLeft,
		Node:             focusedNode,
	}

	var resizeValue int64
	switch resize_direction {
	case DirectionTop:
		if focusedNode.Rect.Y <= defaultStatusBarHeight {
			resizeValue = -loadedNodeConf.DistanceToTop
		} else {
			resizeValue = distanceToTop
		}
		currentNodeConf.DistanceToTop = distanceToTop
	case DirectionBottom:
		if focusedNode.Rect.Y+focusedNode.Rect.Height-defaultStatusBarHeight >= focusedOutput.Rect.Height-defaultStatusBarHeight {
			resizeValue = -loadedNodeConf.DistanceToBottom
		} else {
			resizeValue = distanceToBottom
		}
		currentNodeConf.DistanceToBottom = distanceToBottom
	case DirectionRight:
		if focusedNode.Rect.Width >= focusedOutput.Rect.Width {
			resizeValue = -loadedNodeConf.DistanceToRight
		} else {
			resizeValue = distanceToRight
		}
		currentNodeConf.DistanceToRight = distanceToRight
	case DirectionLeft:
		if focusedNode.Rect.X <= focusedOutput.Rect.Width {
			resizeValue = -loadedNodeConf.DistanceToLeft
		} else {
			resizeValue = distanceToLeft
		}
		currentNodeConf.DistanceToLeft = distanceToLeft
	}

	conf.Nodes[focusedNodeConfigIdentifier] = currentNodeConf

	resizeContext, err := NewResizeContext(resize_direction)
	if err != nil {
		return err
	}

	if err := resizeContext.strategy.Resize(resizeValue); err != nil {
		return err
	}

	return conf.Dump()
}
