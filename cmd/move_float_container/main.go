package main

import (
	moveFloatContainer "github.com/modernpacifist/i3-scripts-go/internal/i3operations/move_float_container"
)

func main() {
	//windowPosition := flag.Int("pos", 0, "New container preset of the window")
	//windowMark := flag.String("mark", "", "Set specified mark options on the focused container")
	//flag.Parse()

	//if *windowPosition == 0 {
	//log.Fatal("Window position was not specified as argument")
	//os.Exit(0)
	//}

	arg := 1
	moveFloatContainer.Execute(arg)

}
