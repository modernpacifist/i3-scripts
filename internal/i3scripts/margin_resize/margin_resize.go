package margin_resize

import (
	"errors"
	"fmt"

	config "github.com/modernpacifist/i3-scripts-go/internal/config/margin_resize"
	common "github.com/modernpacifist/i3-scripts-go/internal/i3scripts"
	"go.i3wm.org/i3/v4"
)

const (
	defaultStatusBarHeight = 35

	DirectionTop    = "top"
	DirectionBottom = "bottom"
	DirectionRight  = "right"
	DirectionLeft   = "left"
)

type ResizeStrategy interface {
	Resize(value int64) error
}

// TopResizeStrategy implements ResizeStrategy for top direction
type TopResizeStrategy struct{}

func (s *TopResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow height %d px, move container up %d px", value, value))
}

type BottomResizeStrategy struct{}

func (s *BottomResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow height %d px", value))
}

type RightResizeStrategy struct{}

func (s *RightResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow width %d px", value))
}

type LeftResizeStrategy struct{}

func (s *LeftResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow width %d px, move container left %d px", value, value))
}

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
	distanceToRight := outputRect.Width - (nodeRect.X - outputRect.X + nodeRect.Width)
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
		if focusedOutput.Rect.X + focusedNode.Rect.Width >= focusedOutput.Rect.Width {
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
