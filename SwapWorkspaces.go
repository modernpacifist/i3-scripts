package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"go.i3wm.org/i3/v4"

	"i3-integration/utils"
)


func getUserInput() (mark string) {
	var userInput string
	var promptMessage string = "Swap ws with number: "

	for {
		userInput = utils.Runi3Input(promptMessage, 1)

		switch {
		case regexp.MustCompile("[0-9]").MatchString(userInput):
			mark = userInput
			return

		default:
			return
		}
	}
}

func swapWorkspaces(currentWs i3.Workspace, swapWs i3.Workspace) (err error) {
	var bufferWsName string = "bufferWs"

	i3.RunCommand(fmt.Sprintf("rename workspace %s to %s", currentWs.Name, bufferWsName))
	i3.RunCommand(fmt.Sprintf("rename workspace %s to %s", swapWs.Name, currentWs.Name))
	i3.RunCommand(fmt.Sprintf("rename workspace %s to %s", bufferWsName, swapWs.Name))

	return
}

func main() {
	userInput := getUserInput()
	if userInput == "" {
		os.Exit(0)
	}

	wsIndex, err := strconv.ParseInt(userInput, 10, 64)
	if err != nil {
		os.Exit(0)
	}

	indWs, err := utils.GetWorkspaceByIndex(wsIndex)
	if err != nil {
		os.Exit(0)
	}

	currentWorkspace, err := utils.GetFocusedWorkspace()
	if err != nil {
		os.Exit(0)
	}

	err = swapWorkspaces(currentWorkspace, indWs)
	if err != nil {
		os.Exit(0)
	}
}
