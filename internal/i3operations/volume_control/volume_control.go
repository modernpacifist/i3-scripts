package volume_control

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"strconv"
	"strings"

	common "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
)

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

func ToggleVolume() {
	if _, err := exec.Command("bash", "-c", "pactl set-sink-mute @DEFAULT_SINK@ toggle").Output(); err != nil {
		log.Fatal(err)
	}

	common.NotifySend(1.5, "VolumeControl: toggled")
}

func RoundVolume() {
	currentVolume := getCurrentVolume()
	roundedVolume := math.Round(currentVolume/5) * 5

	if _, err := exec.Command("bash", "-c", fmt.Sprintf("pactl set-sink-volume @DEFAULT_SINK@ %.f%%", roundedVolume)).Output(); err != nil {
		log.Fatal(err)
	}

	common.NotifySend(1.5, fmt.Sprintf("VolumeControl: rounded to %.f%%", roundedVolume))
}

func ChangeVolumeLevel(changeValue string) {
	_, err := exec.Command("bash", "-c", fmt.Sprintf("pactl set-sink-volume @DEFAULT_SINK@ %s%%", changeValue)).Output()
	if err != nil {
		log.Fatalf("%s: pactl is not installed on this system", err)
	}
	common.NotifySend(1.5, fmt.Sprintf("VolumeControl: %s%%", changeValue))
}
