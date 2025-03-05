package i3operations

import (
	"errors"
	"fmt"
	"log"
	"os/exec"

	"go.i3wm.org/i3/v4"
)

func GetI3Tree() i3.Tree {
	tree, err := i3.GetTree()
	if err != nil {
		log.Fatal(err)
	}

	return tree
}

func GetWorkspaces() ([]i3.Workspace, error) {
	o, err := i3.GetWorkspaces()
	if err != nil {
		return []i3.Workspace{}, err
	}

	return o, nil
}

func GetWorkspaceNodes() *i3.Node {
	i3Tree := GetI3Tree()

	node := i3Tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Focused
	})

	if node == nil {
		log.Fatal(errors.New("could not get focused node"))
	}

	return node
}

func GetFocusedWorkspace() (i3.Workspace, error) {
	o, err := i3.GetWorkspaces()
	if err != nil {
		return i3.Workspace{}, err
	}

	for _, ws := range o {
		if ws.Focused {
			return ws, nil
		}
	}

	return i3.Workspace{}, errors.New("could not get focused workspace")
}

func GetOutputs() ([]i3.Output, error) {
	o, err := i3.GetOutputs()
	if err != nil {
		return []i3.Output{}, err
	}

	return o, nil
}

func GetPrimaryOutput() (i3.Output, error) {
	o, err := i3.GetOutputs()
	if err != nil {
		return i3.Output{}, err
	}

	for _, output := range o {
		if output.Primary {
			return output, nil
		}
	}

	return i3.Output{}, errors.New("could not get primary output")
}

func GetFocusedOutput() (res i3.Output, err error) {
	var focusedWs i3.Workspace

	outputs, err := i3.GetOutputs()
	if err != nil {
		return i3.Output{}, err
	}

	o, err := i3.GetWorkspaces()
	if err != nil {
		return i3.Output{}, errors.New("could not get focused workspace")
	}

	for _, ws := range o {
		if ws.Focused {
			focusedWs = ws
			break
		}
	}

	if focusedWs == (i3.Workspace{}) {
		return i3.Output{}, errors.New("cocused workspace instance is null")
	}

	for _, o := range outputs {
		if o.Active && o.CurrentWorkspace == focusedWs.Name {
			return o, nil
		}
	}

	return i3.Output{}, errors.New("could not get focused output")
}

func GetFocusedNode() *i3.Node {
	i3Tree := GetI3Tree()

	node := i3Tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Focused
	})

	if node == nil {
		log.Fatal(errors.New("could not get focused node"))
	}

	return node
}

func GetNodeMark(node *i3.Node) string {
	if len(node.Marks) == 0 {
		return ""
	}
	return node.Marks[0]
}

// TODO: make function accept error level
func NotifySend(seconds float32, msg string) {
	cmd := fmt.Sprintf("notify-send --expire-time=%.f \"%s\"", seconds*1000, msg)
	// TODO: probably should catch this error via defer <04-12-23, modernpacifist> //
	if _, err := exec.Command("bash", "-c", cmd).Output(); err != nil {
		log.Fatal(err)
	}
}

/* i3InputLimit must be set to 0 for unlimited input*/
// TODO: must change the signature so that the i3-input payload is in the arguments <23-01-24, modernpacifist> //
func Runi3Input(promptMessage string, inputLimit int) (string, error) {
	cmd := fmt.Sprintf("i3-input -P \"%s\" -l %d | grep -oP \"output = \\K.*\" | tr -d '\n'", promptMessage, inputLimit)
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func GetWorkspaceByIndex(index int64) (i3.Workspace, error) {
	o, err := i3.GetWorkspaces()
	if err != nil {
		return i3.Workspace{}, err
	}

	for _, ws := range o {
		if ws.Num == index {
			return ws, nil
		}
	}

	return i3.Workspace{}, errors.New("could not get workspace by specified index")
}

func RunKillCommand() error {
	if _, err := i3.RunCommand("kill"); err != nil {
		return err
	}

	return nil
}

func RunRenameWorkspaceCommand(newWsName string) error {
	if _, err := i3.RunCommand(fmt.Sprintf("rename workspace to %s", newWsName)); err != nil {
		return err
	}

	return nil
}
