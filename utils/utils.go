package utils

import (
	"errors"
	"fmt"
	"log"
	"os/exec"

	"go.i3wm.org/i3/v4"
)

func NotifySend(seconds float32, msg string) {
	_, err := exec.Command("bash", "-c", fmt.Sprintf("notify-send --expire-time=%.f \"%s\"", seconds*1000, msg)).Output()
	// TODO: probably should catch this error via defer <04-12-23, modernpacifist> //
	if err != nil {
		log.Println(err)
	}
}

/* i3InputLimit must be set to 0 for unlimited input*/
/* i3InputLimit must be set to 0 for unlimited input*/
// TODO: must change the signature so that the i3-input payload is in the arguments <23-01-24, modernpacifist> //
func Runi3Input(i3PromptMessage string, i3InputLimit int) string {
	cmd := fmt.Sprintf("i3-input -P \"%s\" -l %d | grep -oP \"output = \\K.*\" | tr -d '\n'", i3PromptMessage, i3InputLimit)
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatal(err)
	}

	return string(output)
}

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
		return n.Focused == true
	})

	if node == nil {
		log.Fatal(errors.New("Could not get focused node"))
	}

	return node
}

func GetFocusedNode() *i3.Node {
	i3Tree := GetI3Tree()

	node := i3Tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Focused == true
	})

	if node == nil {
		log.Fatal(errors.New("Could not get focused node"))
	}

	return node
}

func GetFocusedWorkspace() (i3.Workspace, error) {
	o, err := i3.GetWorkspaces()
	if err != nil {
		return i3.Workspace{}, err
	}

	for _, ws := range o {
		if ws.Focused == true {
			return ws, nil
		}
	}

	return i3.Workspace{}, errors.New("Could not get focused workspace")
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

	return i3.Workspace{}, errors.New("Could not get focused workspace")
}

func GetFocusedOutput() (res i3.Output, err error) {
	outputs, err := i3.GetOutputs()
	if err != nil {
		return i3.Output{}, err
	}

	var focusedWs i3.Workspace

	o, err := i3.GetWorkspaces()
	if err != nil {
		return i3.Output{}, errors.New("Could not get focused workspace")
	}

	for _, ws := range o {
		if ws.Focused == true {
			focusedWs = ws
			break
		}
	}

	if focusedWs == (i3.Workspace{}) {
		return i3.Output{}, errors.New("Focused workspace instance is null")
	}

	for _, o := range outputs {
		if o.Active == true && o.CurrentWorkspace == focusedWs.Name {
			return o, nil
		}
	}

	return i3.Output{}, errors.New("Could not get focused output")
}
