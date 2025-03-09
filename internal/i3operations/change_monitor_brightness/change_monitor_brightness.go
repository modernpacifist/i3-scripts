package change_monitor_brightness

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	common "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
	"go.i3wm.org/i3/v4"
)

const (
	configFilename string = "~/.ScreenDim.json"
)

type Config struct {
	Path       string  `json:"-"`
	Brightness float64 `json:"brightness"`
}

func configConstructor(path string) Config {
	return Config{
		Path:       path,
		Brightness: 0,
	}
}

func (conf *Config) dump() {
	jsonData, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(conf.Path, jsonData, 0644); err != nil {
		log.Fatal(err)
	}
}

func (conf *Config) updateBrightness(newBrightness float64) {
	conf.Brightness = newBrightness
}

func (conf *Config) readFromFile() {
	file, err := os.Open(conf.Path)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	if err := json.Unmarshal(content, conf); err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}
}

func getCurrentBrigthnessXrandr(display string) float64 {
	c := fmt.Sprintf(`xrandr --verbose --current | grep %s -A5 | tail -n1 | awk -F ": " '{print $2}'`, display)
	cmd, err := exec.Command("bash", "-c", c).Output()
	if err != nil {
		log.Fatal(err)
	}

	num, err := strconv.ParseFloat(strings.TrimRight(string(cmd), "\n"), 64)
	if err != nil {
		log.Fatal(err)
	}

	return num
}

func setBrightnessLevel(brightness float64) {
	outputs, _ := i3.GetOutputs()
	for _, o := range outputs {
		if o.Active {
			cmd := fmt.Sprintf(`xrandr --output %s --brightness %f`, o.Name, brightness)
			if _, err := exec.Command("bash", "-c", cmd).Output(); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func resolveBrightnessLevel(currentBrightness float64, changeValue float64) float64 {
	num := currentBrightness + changeValue
	resBrightness := math.Round(num*10) / 10

	if !(0.1 <= resBrightness && resBrightness <= 1.0) {
		log.Fatal("Brightness exceeds the allowed interval")
	}

	return resBrightness
}

func resolveAbsolutePath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Replace(filename, "~", usr.HomeDir, 1)
}

func Execute(arg float64) {
	primaryOutput, err := common.GetPrimaryOutput()
	if err != nil {
		log.Fatal("Could not get primary output")
	}

	absoluteConfigPath := resolveAbsolutePath(configFilename)

	config := configConstructor(absoluteConfigPath)

	if _, err := os.Stat(absoluteConfigPath); os.IsNotExist(err) {
		currentBrightness := getCurrentBrigthnessXrandr(primaryOutput.Name)
		config.updateBrightness(currentBrightness)
	} else {
		config.readFromFile()
	}

	res := resolveBrightnessLevel(config.Brightness, arg)

	if res != 0 {
		config.updateBrightness(res)
		setBrightnessLevel(config.Brightness)
		common.NotifySend(1, fmt.Sprintf("Current brightness: %.1f", config.Brightness))
		config.dump()
	}
}
