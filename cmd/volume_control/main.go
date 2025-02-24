package main

import (
	"flag"
	"log"
	"os"
	"regexp"

	volumeControl "github.com/modernpacifist/i3-scripts-go/internal/i3operations/volume_control"
)

const MAX_VOLUME = 100

func main() {
	toggle := flag.Bool("toggle", false, "Mute/unmute volume")
	round := flag.Bool("round", false, "Round volume value to closest 5")

	flag.Parse()

	if *toggle && *round {
		log.Println("Arguments toggle and round can't be specified at the same time")
		os.Exit(0)
	}

	if *toggle == true {
		volumeControl.ToggleVolume()
		os.Exit(0)
	}

	if *round == true {
		volumeControl.RoundVolume()
		os.Exit(0)
	}

	userInputDummy := flag.Arg(0)

	regex := regexp.MustCompile(`^[-+]\d+$`)
	if !regex.MatchString(userInputDummy) {
		log.Fatal("Wrong input user format")
	}

	volumeControl.ChangeVolumeLevel(userInputDummy)
}
