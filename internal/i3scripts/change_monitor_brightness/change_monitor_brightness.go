package change_monitor_brightness

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"

	config "github.com/modernpacifist/i3-scripts-go/internal/config/change_monitor_brightness"
	common "github.com/modernpacifist/i3-scripts-go/internal/i3scripts"
)

const (
	MinBrightness = 0.1
	MaxBrightness = 1.0
)

func Validate(b float64) error {
	if b < MinBrightness || b > MaxBrightness {
		return fmt.Errorf("brightness %.1f exceeds allowed range [%.1f, %.1f]", b, MinBrightness, MaxBrightness)
	}
	return nil
}

func getCurrentBrightnessXrandr(display string) (float64, error) {
	if !isValidDisplayName(display) {
		return 0, fmt.Errorf("invalid display name: %s", display)
	}

	cmd := exec.Command("xrandr", "--verbose", "--current")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to execute xrandr: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	displayFound := false
	for i, line := range lines {
		if strings.Contains(line, display) {
			displayFound = true
			for j := 1; j < 6 && i+j < len(lines); j++ {
				if strings.Contains(lines[i+j], "Brightness:") {
					parts := strings.Split(lines[i+j], ":")
					if len(parts) != 2 {
						return 0, fmt.Errorf("unexpected xrandr output format")
					}
					brightness, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
					if err != nil {
						return 0, fmt.Errorf("failed to parse brightness value: %w", err)
					}
					return brightness, nil
				}
			}
		}
	}

	if !displayFound {
		return 0, fmt.Errorf("display %s not found", display)
	}

	return 0, fmt.Errorf("brightness information not found for display %s", display)
}

func isValidDisplayName(display string) bool {
	for _, r := range display {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '-' && r != '_' {
			return false
		}
	}

	return true
}

func setBrightnessLevel(brightness float64) error {
	outputs, err := common.GetOutputs()
	if err != nil {
		return fmt.Errorf("could not get outputs: %w", err)
	}

	for _, o := range outputs {
		if o.Active {
			cmd := fmt.Sprintf(`xrandr --output %s --brightness %f`, o.Name, brightness)
			if _, err := exec.Command("bash", "-c", cmd).Output(); err != nil {
				return fmt.Errorf("failed to set brightness for %s: %w", o.Name, err)
			}
		}
	}

	return nil
}

func resolveBrightnessLevel(currentBrightness float64, changeValue float64) (float64, error) {
	num := currentBrightness + changeValue
	resBrightness := math.Round(num*10) / 10

	if err := Validate(resBrightness); err != nil {
		return 0, err
	}

	return resBrightness, nil
}

func Execute(arg float64) {
	primaryOutput, err := common.GetPrimaryOutput()
	if err != nil {
		log.Fatal("Could not get primary output")
	}

	currentBrightness, err := getCurrentBrightnessXrandr(primaryOutput.Name)
	if err != nil {
		log.Fatalf("Error getting current brightness: %v", err)
	}

	config, err := config.Create()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}

	if _, err := os.Stat(config.Path); os.IsNotExist(err) {
		if err = config.UpdateBrightness(currentBrightness); err != nil {
			log.Fatalf("Error updating brightness: %v", err)
		}
		if err = config.Dump(); err != nil {
			log.Fatalf("Error dumping config: %v", err)
		}
	} else {
		if err := config.Load(); err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
	}

	res, err := resolveBrightnessLevel(config.Brightness, arg)
	if err != nil {
		log.Fatalf("Error resolving brightness level: %v", err)
	}

	if res != 0 {
		config.UpdateBrightness(res)
		err = setBrightnessLevel(res)
		if err != nil {
			log.Fatalf("Error setting brightness: %v", err)
		}
		common.NotifySend(1, fmt.Sprintf("Current brightness: %.1f", config.Brightness))
		config.Dump()
	}
}
