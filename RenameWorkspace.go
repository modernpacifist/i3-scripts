package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"go.i3wm.org/i3/v4"
)

func runi3Input() string {
	output, err := exec.Command("bash", "-c", "i3-input -P \"Append title: \" | grep -oP \"output = \\K.*\" | tr -d '\n'").Output()
	if err != nil {
		log.Fatal(err)
	}

	return string(output)
}

func getFocusedWorkspaceNum() int64 {
	var res int64
	o, err := i3.GetWorkspaces()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, ws := range o {
		if ws.Focused == true {
			res = ws.Num
		}
	}

	return res
}

func renamei3Ws(wsIndex int64, newName string) {
	var cmdString string

	if newName == "" {
		cmdString = fmt.Sprintf("rename workspace to %d", wsIndex)
	} else {
		cmdString = fmt.Sprintf("rename workspace to %d:%s", wsIndex, newName)
	}

	i3.RunCommand(cmdString)
}

func main() {
	focusedWS := getFocusedWorkspaceNum()
	if focusedWS == -1 {
		os.Exit(1)
	}

	userPromptResult := runi3Input()

	renamei3Ws(focusedWS, userPromptResult)
}
