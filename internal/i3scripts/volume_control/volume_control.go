package volume_control

import (
	"fmt"
	"log"
	"math"
	"os/exec"
	"regexp"
	"strconv"

	common "github.com/modernpacifist/i3-scripts-go/internal/i3scripts"
)

const (
	NotifySendTimeout = 1.5
	RoundValue        = 5
)

// Optimized volume parsing using regex instead of complex shell pipeline
var volumeRegex = regexp.MustCompile(`(\d+)%`)

func getCurrentVolume() (float64, error) {
	// Simplified command using pactl instead of complex amixer pipeline
	cmd := exec.Command("pactl", "get-sink-volume", "@DEFAULT_SINK@")
	output, err := cmd.Output()
	if err != nil {
		// Fallback to amixer if pactl fails
		return getCurrentVolumeAmixer()
	}

	// Parse volume percentage directly from pactl output
	matches := volumeRegex.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return getCurrentVolumeAmixer() // Fallback
	}

	volume, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return getCurrentVolumeAmixer() // Fallback
	}

	return volume, nil
}

// Fallback method using amixer (optimized)
func getCurrentVolumeAmixer() (float64, error) {
	cmd := exec.Command("amixer", "-D", "pulse", "sget", "Master")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get volume: %w", err)
	}

	// Use regex to extract volume percentage
	matches := volumeRegex.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return 0, fmt.Errorf("could not parse volume from amixer output")
	}

	volume, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, fmt.Errorf("could not convert volume to float: %w", err)
	}

	return volume, nil
}

func ToggleVolume() {
	// Use pactl directly instead of bash wrapper
	cmd := exec.Command("pactl", "set-sink-mute", "@DEFAULT_SINK@", "toggle")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	common.NotifySend(NotifySendTimeout, "VolumeControl: toggled")
}

func RoundVolume() {
	currentVolume, err := getCurrentVolume()
	if err != nil {
		log.Fatal(err)
	}

	roundedVolume := math.Round(currentVolume/RoundValue) * RoundValue

	// Use pactl directly
	cmd := exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", fmt.Sprintf("%.0f%%", roundedVolume))
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	common.NotifySend(NotifySendTimeout, fmt.Sprintf("VolumeControl: rounded to %.0f%%", roundedVolume))
}

func AdjustVolume(changeValue string, maxVolume float64) {
	currentVolume, err := getCurrentVolume()
	if err != nil {
		log.Fatal(err)
	}

	change, err := strconv.ParseFloat(changeValue, 64)
	if err != nil {
		log.Fatal(err)
	}

	newVolume := currentVolume + change
	if newVolume > maxVolume {
		common.NotifySend(NotifySendTimeout, fmt.Sprintf("VolumeControl: cannot adjust volume above %.0f%%", maxVolume))
		return
	}

	// Ensure volume doesn't go below 0
	if newVolume < 0 {
		newVolume = 0
	}

	// Use pactl directly with absolute volume instead of relative change
	cmd := exec.Command("pactl", "set-sink-volume", "@DEFAULT_SINK@", fmt.Sprintf("%.0f%%", newVolume))
	if err := cmd.Run(); err != nil {
		log.Fatalf("pactl command failed: %v (pactl may not be installed)", err)
	}

	common.NotifySend(NotifySendTimeout, fmt.Sprintf("VolumeControl: %s%% (now %.0f%%)", changeValue, newVolume))
}
