package main

import (
	"log"

	"go.i3wm.org/i3/v4"
)

const (
	globalDaemonPort string = ":63334"
)

func getI3Tree() i3.Tree {
	tree, err := i3.GetTree()
	if err != nil {
		log.Fatal(err)
	}

	return tree
}

func getFocusedNode(i3Tree i3.Tree) *i3.Node {
	node := i3Tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Focused == true
	})

	return node
}

// func RunDaemon(i3Tree i3.Tree, config Config) {
func RunDaemon() {
	recv := i3.Subscribe(i3.TickEventType)

	//ln, err := net.Listen("tcp", globalDaemonPort)
	//if err != nil {
	//log.Fatal(err)
	//}
	//defer ln.Close()

	for recv.Next() {
		ev := recv.Event().(*i3.TickEvent)
		//ev := recv.Event().(*i3.TickEvent)
		log.Println(ev)
		//if ev.Change == "focus" {
		//focusedNode := getFocusedNode(i3Tree)
		//i3Tree = getI3Tree()
		//}
	}
	log.Fatal(recv.Close())
}

func main() {
	RunDaemon()
}
