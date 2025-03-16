package main

import (
	"log"

	killContainer "github.com/modernpacifist/i3-scripts-go/internal/i3operations/kill_container"
)

func main() {
	if err := killContainer.Execute(); err != nil {
		log.Fatal(err)
	}
}
