package main

import (
	"log"
	"os"

	renameWorkspace "github.com/modernpacifist/i3-scripts-go/internal/i3operations/rename_workspace"
)

func main() {
	if err := renameWorkspace.Execute(); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
