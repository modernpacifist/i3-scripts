package main

import (
	"log"
	"os"

	common "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
	swapOps "github.com/modernpacifist/i3-scripts-go/internal/i3operations/swap_workspaces"
)

func main() {
	userInput, err := swapOps.GetWorkspaceIndexFromUser()
	if err != nil || userInput == -1 {
		log.Fatal(err)
	}

	targetWorkspace, err := common.GetWorkspaceByIndex(userInput)
	if err != nil {
		log.Fatal(err)
	}

	focusedWorkspace, err := common.GetFocusedWorkspace()
	if err != nil {
		log.Fatal(err)
	}

	if err := swapOps.SwapWorkspaces(focusedWorkspace, targetWorkspace); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
