package main

import (
	"fmt"
	"os"
	"strings"

	"go.i3wm.org/i3/v4"

	utils "github.com/modernpacifist/i3-scripts-go/pkg/i3utils"
)

func replaceSpacesWithUnderscore(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), " ", "_")
}

func renamei3Workspace(wsIndex int64, newName string) {
	var cmd string

	if newName == "" {
		cmd = fmt.Sprintf("rename workspace to %d", wsIndex)
	} else {
		cmd = fmt.Sprintf("rename workspace to %d:%s", wsIndex, replaceSpacesWithUnderscore(newName))
	}

	i3.RunCommand(cmd)
}

func main() {
	focusedWS, err := utils.GetFocusedWorkspace()
	if err != nil {
		utils.NotifySend(2, fmt.Sprintf("RenameWorkspace: %s", err.Error()))
		os.Exit(1)
	}

	userPromptResult := utils.Runi3Input("Append title: ", 0)

	renamei3Workspace(focusedWS.Num, userPromptResult)
}
