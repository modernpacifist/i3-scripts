package keyboard_layout

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	i3operations "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
)

func getCurrentKeyboardLayout() string {
	out, err := exec.Command("bash", "-c", "setxkbmap -query | awk '/layout/{print $2}'").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSuffix(string(out), "\n")
}

func setKeyboardLayout(layout string) {
	_, err := exec.Command("bash", "-c", fmt.Sprintf("setxkbmap %s", layout)).Output()
	if err != nil {
		log.Fatal(err)
	}
}

func cycle(layoutsArray []string) {
	current_layout := getCurrentKeyboardLayout()
	initIndex := -1

	for i, value := range layoutsArray {
		if current_layout == value {
			initIndex = i
			break
		}
	}

	length := len(layoutsArray)
	if initIndex == length-1 && initIndex != -1 {
		setKeyboardLayout(layoutsArray[0])
		i3operations.NotifySend(0.5, "Kb layout: " + layoutsArray[0])
	} else {
		setKeyboardLayout(layoutsArray[(initIndex+1)%length])
		i3operations.NotifySend(0.5, "Kb layout: " + layoutsArray[(initIndex+1)%length])
	}
}

func Execute(layoutsArray []string) {
	cycle(layoutsArray)
}
