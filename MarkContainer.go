package main

import (
	"fmt"
	"os"
	"regexp"

	"go.i3wm.org/i3/v4"

	"i3-integration/utils"
)

func main() {
	var mark string = ""

	userInput := utils.Runi3Input("Mark container (press \"f\" to mark with function keys): ", 1)

	switch {
	case regexp.MustCompile("[0-9]").MatchString(userInput):
		mark = userInput

	case userInput == "f":
		mark = "f"
		userInput = utils.Runi3Input("F key chosen, input function index 1-9: ", 1)
		if !regexp.MustCompile("[0-9]").MatchString(userInput) {
			os.Exit(1)
		}
		mark = mark + userInput

	default:
		os.Exit(0)
	}

	i3.RunCommand(fmt.Sprintf("mark %s", mark))
}
