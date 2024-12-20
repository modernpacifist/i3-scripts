package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd := exec.Command("notify-send", "--expire-time=3000", "hello world")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(output))
}
