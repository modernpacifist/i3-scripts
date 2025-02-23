package rename_workspace

import (
	"fmt"
	"strings"

	common "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
	"go.i3wm.org/i3/v4"
)

func GetWorkspaceNameFromUser() (string, error) {
	var userInput string
	var promptMessage string = "Rename workspace to: "

	for {
		userInput = common.Runi3Input(promptMessage, 0)

		if userInput != "" {
			return userInput, nil
		}
	}
}

func replaceSpacesWithUnderscore(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), " ", "_")
}

func Renamei3Workspace(wsIndex int64, newName string) error {
	var cmd string

	if newName == "" {
		cmd = fmt.Sprintf("rename workspace to %d", wsIndex)
	} else {
		cmd = fmt.Sprintf("rename workspace to %d:%s", wsIndex, replaceSpacesWithUnderscore(newName))
	}

	if _, err := i3.RunCommand(cmd); err != nil {
		return err
	}

	return nil
}
