package main

import (
	// "errors"
	"fmt"
	// "log"

	// "go.i3wm.org/i3/v4"

	utils "github.com/modernpacifist/i3-scripts-go/pkg/i3utils"
)

func StickNode() {

}

func main() {
	focusedNode := utils.GetFocusedNode()
	focusedWindow := utils.NewWindowConfig(focusedNode)

	fmt.Println(focusedWindow)
}