package rename_workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	common "github.com/modernpacifist/i3-scripts-go/internal/i3scripts"
	"go.i3wm.org/i3/v4"
)

type WorkspaceNames struct {
	Workspaces map[string]string `json:"workspaces"`
}

const configFileName = "workspace_names.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, configFileName), nil
}

func loadWorkspaceNames() (*WorkspaceNames, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty workspace names if file doesn't exist
			return &WorkspaceNames{Workspaces: make(map[string]string)}, nil
		}
		return nil, err
	}

	var wsNames WorkspaceNames
	if err := json.Unmarshal(data, &wsNames); err != nil {
		return nil, err
	}

	if wsNames.Workspaces == nil {
		wsNames.Workspaces = make(map[string]string)
	}

	return &wsNames, nil
}

func saveWorkspaceNames(wsNames *WorkspaceNames) error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(wsNames, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func promptUserWorkspaceName(defaultName string) (string, error) {
	var promptMessage string
	if defaultName != "" {
		promptMessage = fmt.Sprintf("Rename workspace to: [%s] ", defaultName)
	} else {
		promptMessage = "Rename workspace to: "
	}

	for {
		userInput, err := common.Runi3Input(promptMessage, 0)
		if err != nil {
			return "", err
		}

		// If user just pressed enter without input
		if userInput == "" {
			// User just pressed enter, use default if available
			if defaultName != "" {
				return defaultName, nil
			}
			// Otherwise, keep prompting
			continue
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

	// Load workspace names from JSON
	wsNames, err := loadWorkspaceNames()
	if err != nil {
		return err
	}

	// Get stored name for this workspace index, if any
	wsIndexStr := strconv.Itoa(int(wsIndex))
	storedName := wsNames.Workspaces[wsIndexStr]

	// Prompt user with the stored name as default
	newName, err := promptUserWorkspaceName(storedName)
	if err != nil {
		return err
	}

	// TODO: add check if newName is spaces
	// BUG: there is no way to clear the ws name currently
	if newName == "" {
		cmd = fmt.Sprintf("rename workspace to %s", currentWsName)
	} else {
		cmd = fmt.Sprintf("rename workspace to %d:%s", wsIndex, replaceSpacesWithUnderscore(newName))
		// Store the new name in the JSON
		wsNames.Workspaces[wsIndexStr] = newName
		if err := saveWorkspaceNames(wsNames); err != nil {
			return err
		}
	}

	if _, err := i3.RunCommand(cmd); err != nil {
		return err
	}

	return nil
}
