package main

import (
	"log"

	markContainer "github.com/modernpacifist/i3-scripts-go/internal/i3operations/mark_container"
)

func main() {
	if err := markContainer.Execute(); err != nil {
		log.Fatal(err)
	}
}
