package i3operations

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"sync"

	"go.i3wm.org/i3/v4"
)

// Cache for frequently accessed data
var (
	treeCache      *i3.Tree
	treeCacheMutex sync.RWMutex
	treeCacheValid bool
)

// OptimizedGetI3Tree returns the i3 tree with caching support
func OptimizedGetI3Tree() (i3.Tree, error) {
	treeCacheMutex.RLock()
	if treeCacheValid && treeCache != nil {
		defer treeCacheMutex.RUnlock()
		return *treeCache, nil
	}
	treeCacheMutex.RUnlock()

	tree, err := i3.GetTree()
	if err != nil {
		return i3.Tree{}, fmt.Errorf("failed to get i3 tree: %w", err)
	}

	treeCacheMutex.Lock()
	treeCache = &tree
	treeCacheValid = true
	treeCacheMutex.Unlock()

	return tree, nil
}

// InvalidateTreeCache invalidates the cached i3 tree
func InvalidateTreeCache() {
	treeCacheMutex.Lock()
	treeCacheValid = false
	treeCacheMutex.Unlock()
}

// OptimizedGetFocusedNode returns the focused node with caching
func OptimizedGetFocusedNode() (*i3.Node, error) {
	tree, err := OptimizedGetI3Tree()
	if err != nil {
		return nil, err
	}

	node := tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Focused
	})

	if node == nil {
		return nil, errors.New("could not find focused node")
	}

	return node, nil
}

// OptimizedGetFocusedWorkspace returns the focused workspace efficiently
func OptimizedGetFocusedWorkspace() (i3.Workspace, error) {
	workspaces, err := i3.GetWorkspaces()
	if err != nil {
		return i3.Workspace{}, fmt.Errorf("failed to get workspaces: %w", err)
	}

	for _, ws := range workspaces {
		if ws.Focused {
			return ws, nil
		}
	}

	return i3.Workspace{}, errors.New("no focused workspace found")
}

// OptimizedRunCommand runs an i3 command with error handling
func OptimizedRunCommand(command string) error {
	_, err := i3.RunCommand(command)
	if err != nil {
		return fmt.Errorf("failed to run i3 command '%s': %w", command, err)
	}
	
	// Invalidate cache after state-changing commands
	InvalidateTreeCache()
	return nil
}

// NotifySendOptimized sends notifications with improved error handling
func NotifySendOptimized(seconds float32, msg string) error {
	// Use notify-send directly instead of bash wrapper
	cmd := exec.Command("notify-send", 
		fmt.Sprintf("--expire-time=%.0f", seconds*1000), 
		msg)
	
	return cmd.Run()
}

// Runi3InputOptimized runs i3-input with improved parsing
func Runi3InputOptimized(promptMessage string, inputLimit int) (string, error) {
	args := []string{"-P", promptMessage}
	if inputLimit > 0 {
		args = append(args, "-l", fmt.Sprintf("%d", inputLimit))
	}
	
	cmd := exec.Command("i3-input", args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("i3-input failed: %w", err)
	}

	// Parse output using regex instead of shell pipeline
	outputRegex := regexp.MustCompile(`output = (.*)`)
	matches := outputRegex.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return "", errors.New("could not parse i3-input output")
	}

	return matches[1], nil
}

// GetWorkspaceByIndexOptimized finds workspace by index efficiently
func GetWorkspaceByIndexOptimized(index int64) (i3.Workspace, error) {
	workspaces, err := i3.GetWorkspaces()
	if err != nil {
		return i3.Workspace{}, fmt.Errorf("failed to get workspaces: %w", err)
	}

	for _, ws := range workspaces {
		if ws.Num == index {
			return ws, nil
		}
	}

	return i3.Workspace{}, fmt.Errorf("workspace with index %d not found", index)
}

// GetNodeByMarkOptimized finds a node by mark with caching
func GetNodeByMarkOptimized(mark string) (*i3.Node, error) {
	tree, err := OptimizedGetI3Tree()
	if err != nil {
		return nil, err
	}

	node := tree.Root.FindChild(func(n *i3.Node) bool {
		for _, m := range n.Marks {
			if m == mark {
				return true
			}
		}
		return false
	})

	if node == nil {
		return nil, fmt.Errorf("node with mark '%s' not found", mark)
	}

	return node, nil
}

// BatchI3Commands executes multiple i3 commands efficiently
func BatchI3Commands(commands []string) error {
	if len(commands) == 0 {
		return nil
	}

	// Join commands with semicolon for batch execution
	batchCommand := ""
	for i, cmd := range commands {
		if i > 0 {
			batchCommand += "; "
		}
		batchCommand += cmd
	}

	err := OptimizedRunCommand(batchCommand)
	if err != nil {
		return fmt.Errorf("batch command failed: %w", err)
	}

	return nil
}