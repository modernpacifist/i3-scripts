package main

import (
	"log"
	"os"

	common "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
	renameOps "github.com/modernpacifist/i3-scripts-go/internal/i3operations/rename_workspace"
)

func main() {
	focusedWorkspace, err := common.GetFocusedWorkspace()
	if err != nil {
		log.Fatal(err)
	}

	userInput, err := renameOps.GetWorkspaceNameFromUser()
	if err != nil {
		log.Fatal(err)
	}

	if err := renameOps.Renamei3Workspace(focusedWorkspace.Num, userInput); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
