package swap_workspaces

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	common "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
	"go.i3wm.org/i3/v4"
)

const (
	bufferWsName = "i3scripts_swap_workspace_buffer_ws"
)

func GetWorkspaceIndexFromUser() (int64, error) {
	var userInput string
	var promptMessage string = "Swap workspace with: "

	for {
		userInput = common.Runi3Input(promptMessage, 1)

		switch {
		case regexp.MustCompile("[0-9]").MatchString(userInput):
			wsIndex, err := strconv.ParseInt(userInput, 10, 64)
			if err != nil {
				return -1, errors.New("failed to parse workspace index")
			}
			return wsIndex, nil

		default:
			return -1, errors.New("invalid workspace index")
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

func SwapWorkspaces(currentWs i3.Workspace, swapWs i3.Workspace) (err error) {
	currentWsName := currentWs.Name
	swapWsName := swapWs.Name

	currentWsNewName := resolveNewWorkspaceName(currentWs, swapWs)
	swapWsNewName := resolveNewWorkspaceName(swapWs, currentWs)

	// "transaction" logic, redo this shit
	if _, err := i3.RunCommand(fmt.Sprintf("rename workspace %s to %s", currentWsName, bufferWsName)); err != nil {
		return fmt.Errorf("failed to rename current workspace: %w", err)
	}

	if _, err := i3.RunCommand(fmt.Sprintf("rename workspace %s to %s", swapWsName, swapWsNewName)); err != nil {
		i3.RunCommand(fmt.Sprintf("rename workspace %s to %s", bufferWsName, currentWsName))
		return fmt.Errorf("failed to rename swap workspace: %w", err)
	}

	if _, err := i3.RunCommand(fmt.Sprintf("rename workspace %s to %s", bufferWsName, currentWsNewName)); err != nil {
		i3.RunCommand(fmt.Sprintf("rename workspace %s to %s", swapWsNewName, swapWsName))
		i3.RunCommand(fmt.Sprintf("rename workspace %s to %s", bufferWsName, currentWsName))
		return fmt.Errorf("failed to complete workspace swap: %w", err)
	}

	return nil
}
