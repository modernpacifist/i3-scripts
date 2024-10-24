package main

import (
	"errors"
	"fmt"
	"log"

	"go.i3wm.org/i3/v4"

	"i3-scripts-go/utils"
)

type WindowConfig struct {
	ID     int64  `json:"ID"`
	X      int64  `json:"X"`
	Y      int64  `json:"Y"`
	Width  int64  `json:"Width"`
	Height int64  `json:"Height"`
	Mark   string `json:"Mark"`
	Sticky bool   `json:"Sticky"`
}

func WindowConfigConstructor(node *i3.Node) WindowConfig {
	// TODO: is the node does not contain a mark, just use a container id <17-11-23, modernpacifist> //
	return WindowConfig{
		ID:     node.Window,
		X:      node.Rect.X,
		Y:      node.Rect.Y,
		Width:  node.Rect.Width,
		Height: node.Rect.Height,
		Sticky: node.Sticky,
	}
}

func getFocusedNode() *i3.Node {
	i3Tree := utils.GetI3Tree()

	node := i3Tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Focused == true
	})

	if node == nil {
		log.Fatal(errors.New("Could not find focused node"))
	}

	return node
}

func StickNode() {

}

func main() {
	focusedNode := getFocusedNode()
	focusedWindow := WindowConfigConstructor(focusedNode)

	fmt.Println(focusedWindow)
}
