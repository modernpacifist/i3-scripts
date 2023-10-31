package main

import (
	"fmt"
	"log"

	"go.i3wm.org/i3/v4"
)

var i3Tree i3.Tree

const Mark = "last"

func getFocusedNode() *i3.Node {
	node := i3Tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Focused == true
	})

	return node
}

// TODO: must get a container with mark _last <31-10-23, modernpacifist> //
func getInactiveNode() *i3.Node {
	node := i3Tree.Root.FindChild(func(n *i3.Node) bool {
		for _, m := range n.Marks {
			return m == Mark
		}
		return false
	})

	return node
}

func focusPreviousContainer(focusedNode, inactiveNode *i3.Node) {
	i3.RunCommand(fmt.Sprintf("[con_mark=\"^%s$\"] focus", Mark))
	i3.RunCommand(fmt.Sprintf("[con_id=%d] mark --add %s", focusedNode.ID, Mark))
}

func getI3Tree() i3.Tree {
	tree, err := i3.GetTree()
	if err != nil {
		log.Fatal(err)
	}

	return tree
}

func init() {
	i3Tree = getI3Tree()
}

func main() {
	focusedNode := getFocusedNode()

	focusedInactiveNode := getInactiveNode()

	focusPreviousContainer(focusedNode, focusedInactiveNode)
}
