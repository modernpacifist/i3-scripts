package main

import (
	// "errors"
	"fmt"
	// "log"

	// "go.i3wm.org/i3/v4"

	"github.com/modernpacifist/i3-scripts-go/internal/configs"
	"github.com/modernpacifist/i3-scripts-go/internal/i3operations"
)

func StickNode() {

}

func main() {
	focusedNode := i3operations.GetFocusedNode()
	focusedWindow := configs.NewWindowConfig(focusedNode)

	fmt.Println(focusedWindow)
}
