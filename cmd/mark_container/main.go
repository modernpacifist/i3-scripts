package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/modernpacifist/i3-scripts-go/internal/i3operations"
	"go.i3wm.org/i3/v4"
)

func getMarkFromUser() (mark string) {
	var userInput string
	var promptMessage string = `Mark container (press "f" to mark with function keys): `

	for {
		userInput = i3operations.Runi3Input(promptMessage, 1)

		switch {
		case regexp.MustCompile("[0-9]").MatchString(userInput):
			if mark == "f" && userInput == "0" {
				return "f10"
			}
			mark = mark + userInput
			return

		case userInput == "f":
			mark = mark + userInput
			promptMessage = "Function key chosen, input function index 0-9: "

			continue

		case len(mark) > 1:
			return

		default:
			return
		}
	}
}

func main() {
	mark := getMarkFromUser()
	if mark == "" {
		os.Exit(0)
	}

	matched, err := regexp.MatchString("^(f?[0-9]{1,2})$", mark)
	if err != nil {
		os.Exit(0)
	}

	if !matched {
		os.Exit(0)
	}

	i3.RunCommand(fmt.Sprintf("mark %s", mark))
}
