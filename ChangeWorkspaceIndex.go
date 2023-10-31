package main

import (
	"os"
	"errors"
	"flag"
	"fmt"
	"strings"
	"os/exec"

	"go.i3wm.org/i3/v4"
)

func getCurrentWsNum() (int64, string, error) {
	workspaces, _ := i3.GetWorkspaces()

	for _, ws := range workspaces {
		if ws.Focused == true {
			return ws.Num, ws.Name, nil
		}
	}

	return 0, "", errors.New("The error")
}

func getExistingWs() []int64 {
	var res []int64
	workspaces, _ := i3.GetWorkspaces()

	for _, ws := range workspaces {
		res = append(res, ws.Num)
	}

	return res
}

func notifySend(time int, msg string) {
	cmd := exec.Command("notify-send", fmt.Sprintf("--expire-time=%d", time), msg)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
}

func checkOccupiedWsIndex(wsNum int64, numbers []int64) {
	for _, num := range numbers {
		if num == wsNum {
			notifySend(3000, fmt.Sprintf("Index %d already occupied", wsNum))
			os.Exit(0)
		}
	}
}

func main() {
	newWsIndex := flag.Int64("index", -1, "New index")

	flag.Parse()

	if *newWsIndex == -1 {
		panic("Index was not specified")
	}

	existingIndices := getExistingWs()
	checkOccupiedWsIndex(*newWsIndex, existingIndices)

	_, currentWsName, err := getCurrentWsNum()
	// TODO: handle error and exit explicitly <13-10-23, modernpacifist> //
	if err != nil {
		panic("Could not retrieve")
	}

	parts := strings.Split(currentWsName, ":")

	if len(parts) == 2 {
		i3.RunCommand(fmt.Sprintf("rename workspace to %d:%s", *newWsIndex, parts[1]))
	}

	if len(parts) == 1 {
		i3.RunCommand(fmt.Sprintf("rename workspace to %d", *newWsIndex))
	}
}
