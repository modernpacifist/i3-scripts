package utils

import (
	"fmt"
	"log"
	"os/exec"
)

func NotifySend(seconds float32, msg string) {
	_, err := exec.Command("bash", "-c", fmt.Sprintf("notify-send --expire-time=%.f \"%s\"", seconds*1000, msg)).Output()
	// TODO: probably should catch this error via defer <04-12-23, modernpacifist> //
	if err != nil {
		log.Println(err)
	}
}
