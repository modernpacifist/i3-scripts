package mark_container

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/modernpacifist/i3-scripts-go/internal/i3scripts"
	"go.i3wm.org/i3/v4"
)

func getMarkFromUser() (mark string) {
	var promptMessage string = `Mark container (press "f" to mark with function keys): `

	for {
		userInput, err := i3scripts.Runi3Input(promptMessage, 1)
		if err != nil {
			return ""
		}

		switch {
		case regexp.MustCompile("[0-9]").MatchString(userInput):
			if mark == "f" && userInput == "0" {
				return "f10"
			}
			mark += userInput
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

func Execute() error {
	mark := getMarkFromUser()
	if mark == "" {
		return errors.New("no mark provided")
	}

	matched, err := regexp.MatchString("^(f?[0-9]{1,2})$", mark)
	if err != nil {
		return errors.New("invalid mark format")
	}

	if !matched {
		return errors.New("invalid mark format")
	}

	if _, err := i3.RunCommand(fmt.Sprintf("mark %s", mark)); err != nil {
		return errors.New("failed to mark container")
	}

	return nil
}
