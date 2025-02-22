package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"go.i3wm.org/i3/v4"

	utils "github.com/modernpacifist/i3-scripts-go/pkg/i3utils"
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

func resolveNewWorkspaceName(ws1 i3.Workspace, ws2 i3.Workspace) string {
	result := strings.Split(ws1.Name, ":")
	if len(result) > 1 {
		return fmt.Sprintf("%d:%s", ws2.Num, result[1])
	}

	return fmt.Sprintf("%d", ws2.Num)
}

func swapWorkspaces(currentWs i3.Workspace, swapWs i3.Workspace) (err error) {
	var bufferWsName string = "BufferWs"

	currentWsNewName := resolveNewWorkspaceName(currentWs, swapWs)
	swapWsNewName := resolveNewWorkspaceName(swapWs, currentWs)

	i3.RunCommand(fmt.Sprintf("rename workspace %s to %s", currentWs.Name, bufferWsName))
	i3.RunCommand(fmt.Sprintf("rename workspace %s to %s", swapWs.Name, swapWsNewName))
	i3.RunCommand(fmt.Sprintf("rename workspace %s to %s", bufferWsName, currentWsNewName))

	return
}

func main() {
	userInput := getUserInput()
	if userInput == "" {
		os.Exit(0)
	}
	//userInput := "4"

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
