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
	NormalizeValue(output i3.Output, nodeConf config.NodeConfig) int64
}

// TopResizeStrategy implements ResizeStrategy for top direction
type TopResizeStrategy struct{}

func (s *TopResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow height %d px, move container up %d px", value, value))
}

func (s *TopResizeStrategy) NormalizeValue(output i3.Output, nodeConf config.NodeConfig) int64 {
	// if val == 0 && nodeConf.Node.Rect.Y > output.Rect.Y {
	if nodeConf.Node.Rect.Y-defaultStatusBarHeight >= output.Rect.Y {
		return nodeConf.Node.Rect.Y - defaultStatusBarHeight
	}
	return -(nodeConf.Node.Rect.Y - defaultStatusBarHeight)
}

type BottomResizeStrategy struct{}

func (s *BottomResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow height %d px", value))
}

func (s *BottomResizeStrategy) NormalizeValue(output i3.Output, nodeConf config.NodeConfig) int64 {
	if nodeConf.Node.Rect.Y+nodeConf.Node.Rect.Height == output.Rect.Height {
		return -(output.Rect.Height - nodeConf.Node.Rect.Height)
	}
	return -(nodeConf.Node.Rect.Y + nodeConf.Node.Rect.Height - output.Rect.Height)
}

type RightResizeStrategy struct{}

func (s *RightResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow width %d px", value))
}

func (s *RightResizeStrategy) NormalizeValue(output i3.Output, nodeConf config.NodeConfig) int64 {
	if nodeConf.Node.Rect.X+nodeConf.Node.Rect.Width == output.Rect.Width {
		return -(output.Rect.Width - nodeConf.Node.Rect.Width)
	}
	return -(nodeConf.Node.Rect.X + nodeConf.Node.Rect.Width - output.Rect.Width)
}

type LeftResizeStrategy struct{}

func (s *LeftResizeStrategy) Resize(value int64) error {
	return common.RunI3Command(fmt.Sprintf("resize grow width %d px, move container left %d px", value, value))
}

func (s *LeftResizeStrategy) NormalizeValue(output i3.Output, nodeConf config.NodeConfig) int64 {
	if nodeConf.Node.Rect.X == output.Rect.X {
		return -(output.Rect.Width - nodeConf.Node.Rect.Width)
	}
	return -(nodeConf.Node.Rect.X - output.Rect.X)
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
	loadedNodeConf, exists := conf.Nodes[focusedNodeConfigIdentifier]
	if !exists {
		conf.Nodes[focusedNodeConfigIdentifier] = config.NodeConfig{
			Node: focusedNode,
		}
		loadedNodeConf = conf.Nodes[focusedNodeConfigIdentifier]
	}

	// instantly update config file with new data
	conf.Nodes[focusedNodeConfigIdentifier] = config.NodeConfig{
		Node: focusedNode,
	}

	distanceToTop, distanceToBottom, distanceToRight, distanceToLeft := getScreenMargins(focusedOutput, focusedNode)

	nodeConf := conf.Nodes[focusedNodeConfigIdentifier]

	nodeConf.DistanceToTop = distanceToTop
	nodeConf.DistanceToBottom = distanceToBottom
	nodeConf.DistanceToRight = distanceToRight
	nodeConf.DistanceToLeft = distanceToLeft

	conf.Nodes[focusedNodeConfigIdentifier] = nodeConf

	if err := conf.Dump(); err != nil {
		return err
	}

	resizeContext, err := NewResizeContext(arg)
	if err != nil {
		return err
	}

	var resizeValue int64
	switch arg {
	case DirectionTop:
		resizeValue = resizeContext.strategy.NormalizeValue(focusedOutput, loadedNodeConf)
	case DirectionBottom:
		resizeValue = resizeContext.strategy.NormalizeValue(focusedOutput, loadedNodeConf)
	case DirectionRight:
		resizeValue = resizeContext.strategy.NormalizeValue(focusedOutput, loadedNodeConf)
	case DirectionLeft:
		resizeValue = resizeContext.strategy.NormalizeValue(focusedOutput, loadedNodeConf)
	}

	return resizeContext.strategy.Resize(resizeValue)
}
