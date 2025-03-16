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
	NormalizeValue(output i3.Output, pastNode config.NodeConfig, currentValue int64) int64
}

// TopResizeStrategy implements ResizeStrategy for top direction
type TopResizeStrategy struct{}

func (s *TopResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow height %d px, move container up %d px", value, value))
}

func (s *TopResizeStrategy) NormalizeValue(output i3.Output, pastNode config.NodeConfig, currentValue int64) int64 {
	if currentValue == 0 {
		return -(output.Rect.Height - pastNode.Node.Rect.Height - defaultStatusBarHeight)
	}
	if currentValue == 0 {
		return -(output.Rect.Height - pastNode.Node.Rect.Height - pastNode.Node.Rect.Y - defaultStatusBarHeight)
	}
	return currentValue
}

type BottomResizeStrategy struct{}

func (s *BottomResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow height %d px", value))
}

func (s *BottomResizeStrategy) NormalizeValue(output i3.Output, pastNode config.NodeConfig, currentValue int64) int64 {
	if currentValue == 0 && pastNode.Node.Rect.Y+pastNode.Node.Rect.Height == output.Rect.Height {
		return -(output.Rect.Height - pastNode.Node.Rect.Height)
	}
	if currentValue == 0 {
		return -(output.Rect.Height - pastNode.Node.Rect.Height - pastNode.Node.Rect.Y)
	}
	return currentValue
}

type RightResizeStrategy struct{}

func (s *RightResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow width %d px", value))
}

func (s *RightResizeStrategy) NormalizeValue(output i3.Output, pastNode config.NodeConfig, currentValue int64) int64 {
	if currentValue == 0 && pastNode.Node.Rect.Width == output.Rect.Width {
		return -(output.Rect.Width - pastNode.Node.Rect.Width)
	}
	if currentValue == 0 {
		return -(output.Rect.Width - pastNode.Node.Rect.Width - (pastNode.Node.Rect.X - output.Rect.X))
	}
	return currentValue
}

type LeftResizeStrategy struct{}

func (s *LeftResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow width %d px, move container left %d px", value, value))
}

func (s *LeftResizeStrategy) NormalizeValue(output i3.Output, pastNode config.NodeConfig, currentValue int64) int64 {
	if currentValue == 0 && pastNode.Node.Rect.Width == output.Rect.Width {
		return -(output.Rect.Width - pastNode.Node.Rect.Width)
	}
	if currentValue == 0 {
		return -(pastNode.Node.Rect.Width)
	}
	return currentValue
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

	distanceTop := nodeRect.Y - defaultStatusBarHeight
	distanceBottom := outputRect.Y + outputRect.Height - (nodeRect.Y + nodeRect.Height)
	distanceRight := outputRect.X + outputRect.Width - (nodeRect.X + nodeRect.Width)
	distanceLeft := nodeRect.X - outputRect.X

	return distanceTop, distanceBottom, distanceRight, distanceLeft
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

	distanceTop, distanceBottom, distanceRight, distanceLeft := getScreenMargins(focusedOutput, focusedNode)

	resizeContext, err := NewResizeContext(arg)
	if err != nil {
		return err
	}

	var resizeValue int64
	switch arg {
	case DirectionTop:
		resizeValue = resizeContext.strategy.NormalizeValue(focusedOutput, pastNodeConfig, distanceTop)
	case DirectionBottom:
		resizeValue = resizeContext.strategy.NormalizeValue(focusedOutput, pastNodeConfig, distanceBottom)
	case DirectionRight:
		resizeValue = resizeContext.strategy.NormalizeValue(focusedOutput, pastNodeConfig, distanceRight)
	case DirectionLeft:
		resizeValue = resizeContext.strategy.NormalizeValue(focusedOutput, pastNodeConfig, distanceLeft)
	}

	return resizeContext.strategy.Resize(resizeValue)
}
