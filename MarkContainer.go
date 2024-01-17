package main

import (
	"fmt"
	"os"
	"regexp"

	"go.i3wm.org/i3/v4"

	"i3-integration/utils"
)

func getUserInput() (mark string) {
	var userInput string
	var promptMessage string = "Mark container (press \"f\" to mark with function keys): "

	for {
		userInput = utils.Runi3Input(promptMessage, 1)

		switch {
		case regexp.MustCompile("[0-9]").MatchString(userInput):
			mark = mark + userInput
			return

		case userInput == "f":
			mark = mark + userInput
			promptMessage = "Function key chosen, input function index 1-9: "
			break

		case len(mark) > 1:
			return

		default:
			return
		}
	}
}

func main() {
	var mark string = ""

	mark = getUserInput()
	if mark == "" {
		os.Exit(0)
	}

	matched, err := regexp.MatchString("^(f?[0-9])$", mark)
	if err != nil {
		os.Exit(0)
	}

	if !matched {
		os.Exit(0)
	}

	i3.RunCommand(fmt.Sprintf("mark %s", mark))
}
