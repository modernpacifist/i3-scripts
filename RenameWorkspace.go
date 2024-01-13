package main

import (
	"fmt"
	"os"
	"strings"

	"go.i3wm.org/i3/v4"

	"i3-integration/utils"
)

func replaceSpacesWithUnderscore(s string) string {
	trimmed := strings.TrimSpace(s)
	return strings.ReplaceAll(trimmed, " ", "_")
}

func renamei3Ws(wsIndex int64, newName string) {
	var cmdString string

	if newName == "" {
		cmdString = fmt.Sprintf("rename workspace to %d", wsIndex)
	} else {
		newName = replaceSpacesWithUnderscore(newName)
		cmdString = fmt.Sprintf("rename workspace to %d:%s", wsIndex, newName)
	}

	i3.RunCommand(cmdString)
}

func main() {
	focusedWS, err := utils.GetFocusedWorkspace()
	if err != nil {
		utils.NotifySend(2, fmt.Sprintf("RenameWorkspace: %s", err.Error()))
		os.Exit(1)
	}

	userPromptResult := utils.Runi3Input("Append title: ", 0)

	renamei3Ws(focusedWS.Num, userPromptResult)
}
