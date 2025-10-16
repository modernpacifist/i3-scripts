package rename_workspace

import (
	"fmt"
	"strings"

	common "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
	"go.i3wm.org/i3/v4"
)

func promptUserWorkspaceName() (string, error) {
	var promptMessage string = "Rename workspace to: "

	for {
		userInput, err := common.Runi3Input(promptMessage, 0)
		if err != nil || userInput == "" {
			return "", err
		}

		if userInput != "" {
			return userInput, nil
		}
	}
}

func replaceSpacesWithUnderscore(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), " ", "_")
}

func Execute() error {
	var cmd string

	focusedWorkspace, err := common.GetFocusedWorkspace()
	if err != nil {
		return err
	}

	wsIndex := focusedWorkspace.Num
	currentWsName := focusedWorkspace.Name

	newName, err := promptUserWorkspaceName()
	if err != nil {
		return err
	}

	// TODO: add check if newName is spaces
	// BUG: there is no way to clear the ws name currently
	if newName == "" {
		cmd = fmt.Sprintf("rename workspace to %s", currentWsName)
	} else {
		cmd = fmt.Sprintf("rename workspace to %d:%s", wsIndex, replaceSpacesWithUnderscore(newName))
	}

	if _, err := i3.RunCommand(cmd); err != nil {
		return err
	}

	return nil
}
