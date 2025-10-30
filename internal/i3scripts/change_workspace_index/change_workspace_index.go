package change_workspace_index

import (
	"errors"
	"fmt"
	"os"
	"strings"

	common "github.com/modernpacifist/i3-scripts-go/internal/i3scripts"
	"go.i3wm.org/i3/v4"
)

func checkIndexAvailability(wsNum int64, numbers []i3.Workspace) {
	for _, num := range numbers {
		if num.Num == wsNum {
			common.NotifySend(3, fmt.Sprintf("Index %d already occupied", wsNum))
			os.Exit(0)
		}
	}
}

func Execute(newWsIndex int64) error {
	var currentWsName string
	var newWsName string
	var err error

	existingWs, err := common.GetWorkspaces()
	if err != nil {
		return err
	}

	checkIndexAvailability(newWsIndex, existingWs)

	currentWs, err := common.GetFocusedWorkspace()
	if err != nil {
		return err
	}

	currentWsName = currentWs.Name

	parts := strings.Split(currentWsName, ":")

	switch len(parts) {
	case 1:
		newWsName = fmt.Sprintf("%d", newWsIndex)
	case 2:
		newWsName = fmt.Sprintf("%d:%s", newWsIndex, parts[1])
	default:
		return errors.New("Invalid workspace name")
	}

	if err := common.RunRenameWorkspaceCommand(newWsName); err != nil {
		return err
	}

	return nil
}
