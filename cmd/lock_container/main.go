package main

import (
	"log"

	lockContainer "github.com/modernpacifist/i3-scripts-go/internal/i3scripts/lock_container"
)

func main() {
	if err := lockContainer.Execute(); err != nil {
		log.Fatal(err)
	}
}
