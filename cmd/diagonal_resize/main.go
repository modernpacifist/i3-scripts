package main

import (
	"flag"
	"fmt"

	"go.i3wm.org/i3/v4"
)

func resize(resizeValue int) {
	i3.RunCommand(fmt.Sprintf("resize grow width %d px or %d ppt", resizeValue, resizeValue))
	i3.RunCommand(fmt.Sprintf("resize grow height %d px or %d ppt", resizeValue, resizeValue))
	i3.RunCommand(fmt.Sprintf("move container left %d px", resizeValue/2))
	i3.RunCommand(fmt.Sprintf("move container up %d px", resizeValue/2))
}

func main() {
	resizeValue := flag.Int("size", 0, "Resize value")

	flag.Parse()

	resize(*resizeValue)
}
