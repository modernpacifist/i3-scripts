package main

import (
	"log"

	swapWorkspaces "github.com/modernpacifist/i3-scripts-go/internal/i3operations/swap_workspaces"
)

func main() {
	if err := swapWorkspaces.Execute(); err != nil {
		log.Fatal(err)
	}
}
