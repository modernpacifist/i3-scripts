package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"i3-integration/utils"
)

const MAX_VOLUME = 100

func getCurrentVolume() float64 {
	out, err := exec.Command("bash", "-c", "amixer -D pulse sget Master | grep 'Left:' | awk -F'[][]' '{ print $2 }' | tr -d '%'").Output()
	if err != nil {
		log.Fatal(err)
	}

	num, err := strconv.ParseFloat(strings.TrimSuffix(string(out), "\n"), 64)
	if err != nil {
		log.Fatal(err)
	}

	return num
}

func toggleVolume() {
	_, err := exec.Command("bash", "-c", "pactl set-sink-mute @DEFAULT_SINK@ toggle").Output()
	if err != nil {
		log.Fatal(err)
	}

	utils.NotifySend(1.5, "VolumeControl: toggled")
}

func roundVolume() {
	currentVolume := getCurrentVolume()
	roundedVolume := math.Round(currentVolume/5) * 5

	_, err := exec.Command("bash", "-c", fmt.Sprintf("pactl set-sink-volume @DEFAULT_SINK@ %.f%%", roundedVolume)).Output()
	if err != nil {
		log.Fatal(err)
	}

	utils.NotifySend(1.5, fmt.Sprintf("VolumeControl: rounded to %.f%%", roundedVolume))
}

func changeVolumeLevel(changeValue string) {
	_, err := exec.Command("bash", "-c", fmt.Sprintf("pactl set-sink-volume @DEFAULT_SINK@ %s%%", changeValue)).Output()
	if err != nil {
		log.Fatal(err)
	}
	utils.NotifySend(1.5, fmt.Sprintf("VolumeControl: %s%%", changeValue))
}

func main() {
	toggle := flag.Bool("toggle", false, "Mute/unmute volume")
	round := flag.Bool("round", false, "Round volume value to closest 5")

	flag.Parse()

	if *toggle && *round {
		log.Println("Arguments toggle and round can't be specified at the same time")
		os.Exit(0)
	}

	if *toggle == true {
		toggleVolume()
		os.Exit(0)
	}

	if *round == true {
		roundVolume()
		os.Exit(0)
	}

	userInputDummy := flag.Arg(0)

	regex := regexp.MustCompile(`^[-+]\d+$`)
	if !regex.MatchString(userInputDummy) {
		log.Fatal("Wrong input user format")
	}

	changeVolumeLevel(userInputDummy)
}
