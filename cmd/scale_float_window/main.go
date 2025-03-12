package main

import (
	"log"

	scaleFloatContainer "github.com/modernpacifist/i3-scripts-go/internal/i3operations/scale_float_container"
)

func main() {
	// mode := flag.String("mode", "", "Specify resize mode")
	// widen := flag.Int("widen", 0, "Specify resize widen")
	// flag.Parse()

	dummy := "w"

	if err := scaleFloatContainer.Execute(dummy); err != nil {
		log.Fatal(err)
	}
}
