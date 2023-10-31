package main

import (
	"fmt"
	"os"
	"flag"

	//"go.i3wm.org/i3/v4"
)

func main() {
	windowPosition := flag.Int("pos", 0, "New container preset of the window")
	flag.Parse()

	if *windowPosition == 0 {
		os.Exit(0)
	}

	fmt.Println(*windowPosition)
}
