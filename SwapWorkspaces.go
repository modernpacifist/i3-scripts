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

func resolveWorkspaceName(wsName string) string {
	re := regexp.MustCompile(`\d+:(.+)`)
	result := re.FindStringSubmatch(wsName)
	if len(result) > 1 {
		return result[1]
	}

	return ""
}

func swapWorkspaces(currentWs i3.Workspace, swapWs i3.Workspace) (err error) {
	a := resolveWorkspaceName(currentWs.Name)
	b := resolveWorkspaceName(swapWs.Name)

	if a != "" {
		a = fmt.Sprintf("%d:%s", swapWs.Num, a)
	} else {
		a = fmt.Sprintf("%d", swapWs.Num)
	}

	if b == "" {
		b = fmt.Sprintf("%d", currentWs.Num)
	} else {
		b = fmt.Sprintf("%d:%s", currentWs.Num, b)
	}

	i3.RunCommand(fmt.Sprintf("rename workspace %s to %s", currentWs.Name, a))
	i3.RunCommand(fmt.Sprintf("rename workspace %s to %s", swapWs.Name, b))

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
