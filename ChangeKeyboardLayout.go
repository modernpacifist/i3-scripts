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

func setKbdlayout(layout string) {
	_, err := exec.Command("bash", "-c", fmt.Sprintf("setxkbmap %s", layout)).Output()
	if err == nil {
		exec.Command("bash", "-c", fmt.Sprintf("notify-send --expire-time=1000 \"Kb layout: %s\"", layout)).Output()
	}
}

func cycle(layoutsArray []string) {
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
		setKbdlayout(layoutsArray[0])
	} else {
		setKbdlayout(layoutsArray[(initIndex+1)%length])
	}
}

func main() {
	layouts := flag.String("layouts", "", "New index")

	flag.Parse()

	if *layouts == "" {
		//panic("layouts flag was not specified")
		log.Fatalln("layouts flag was not specified")
	}

	layoutsArray := strings.Split(*layouts, "/")

	cycle(layoutsArray)
}
