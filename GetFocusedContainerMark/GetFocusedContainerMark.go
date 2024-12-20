package main

import (
	"errors"
	"fmt"
	"log"

	"go.i3wm.org/i3/v4"
)

var i3Tree i3.Tree

func getFocusedNode() *i3.Node {
	node := i3Tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Focused == true
	})

	if node == nil {
		log.Fatal(errors.New("Could not find focused node"))
	}

	return node
}

func getI3Tree() i3.Tree {
	tree, err := i3.GetTree()
	if err != nil {
		log.Fatal(err)
	}

	return tree
}

func getFocusedNodeMark(node *i3.Node) string {
	fmt.Println(node.Window)
	if len(node.Marks) == 0 {
		return ""
	}
	return node.Marks[0]
}

func init() {
	i3Tree = getI3Tree()
}

func main() {
	node := getFocusedNode()

	mark := getFocusedNodeMark(node)
	fmt.Println(mark)
}
