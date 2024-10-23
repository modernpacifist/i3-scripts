package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func getKbdlayout() string {
	out, err := exec.Command("bash", "-c", "setxkbmap -query | awk '/layout/{print $2}'").Output()
	if err != nil {
		log.Fatal(err)
	}
	return strings.TrimSuffix(string(out), "\n")
}

func setKbdlayout(layout string, expireTime float64) {
	_, err := exec.Command("bash", "-c", fmt.Sprintf("setxkbmap %s", layout)).Output()
	if err == nil {
		cmd := fmt.Sprintf("notify-send --expire-time=%.0f \"Kb layout: %s\"", expireTime*1000, layout)
		exec.Command("bash", "-c", cmd).Output()
	}
}

func cycle(layoutsArray []string, expireTime float64) {
	current_layout := getKbdlayout()
	initIndex := -1

	for i, value := range layoutsArray {
		if current_layout == value {
			initIndex = i
			break
		}
	}

	// TODO; wrong exit
	if initIndex == -1 {
		os.Exit(1)
	}

	length := len(layoutsArray)
	if initIndex == length-1 {
		setKbdlayout(layoutsArray[0], expireTime)
	} else {
		setKbdlayout(layoutsArray[(initIndex+1)%length], expireTime)
	}
}

func main() {
	layouts := flag.String("layouts", "", "New index")
	notifyTimeout := flag.Float64("timeout", 0, "Notification timeout")

	flag.Parse()

	if *layouts == "" {
		//panic("layouts flag was not specified")
		log.Fatalln("layouts flag was not specified")
	}

	if *notifyTimeout == 0 {
		//panic("layouts flag was not specified")
		log.Fatalln("timeout flag was not specified")
	}

	layoutsArray := strings.Split(*layouts, "/")

	cycle(layoutsArray, *notifyTimeout)
}
