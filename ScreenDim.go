package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	"go.i3wm.org/i3/v4"
)

const ConfigFilename = ".ScreenDim.json"

type JsonStruct struct {
	Brightness float64 `json:"currentBrightness"`
}

type Config struct {
	Path       string  `json:"-"`
	Brightness float64 `json:"Brightness"`
}

func configConstructor(filename string) Config {
	return Config{
		Path:       filename,
		Brightness: 0,
	}
}

func (conf *Config) dump() {
	jsonData, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(conf.Path, jsonData, 0644)
	if err != nil {
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

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	err = json.Unmarshal(content, conf)
	if err != nil {
		log.Fatal("Error unmarshaling JSON:", err)
	}
}

func NotifySend(brightness float64) {
	msg := fmt.Sprintf("Brightness: %.1f", brightness)
	_, err := exec.Command("notify-send", "--expire-time=1000", msg).Output()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getCurrentBrigthnessXrandr(display string) float64 {
	c := fmt.Sprintf("xrandr --verbose --current | grep %s -A5 | tail -n1 | awk -F \": \" '{print $2}'", display)
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

func SetBrightnessLevel(brightness float64) {
	outputs, _ := i3.GetOutputs()
	for _, o := range outputs {
		if o.Active == true {
			c := fmt.Sprintf("xrandr --output %s --brightness %f", o.Name, brightness)
			_, err := exec.Command("bash", "-c", c).Output()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func getPrimaryOutput() string {
	outputs, _ := i3.GetOutputs()
	for _, output := range outputs {
		if output.Primary == true {
			return output.Name
		}
	}
	return ""
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

	return strings.Replace(fmt.Sprintf("~/%s", filename), "~", usr.HomeDir, 1)
}

func main() {
	var floatValue float64
	flag.Float64Var(&floatValue, "f", 0.0, "Brightness value")
	flag.Parse()

	if floatValue == 0 {
		log.Fatal("The brightness level was not specified")
	}

	primaryOutput := getPrimaryOutput()
	if primaryOutput == "" {
		log.Fatal("Could not get primary output")
	}

	absoluteConfigPath := resolveAbsolutePath(ConfigFilename)

	config := configConstructor(absoluteConfigPath)

	_, err := os.Stat(absoluteConfigPath)
	if os.IsNotExist(err) {
		currentBrightness := getCurrentBrigthnessXrandr(primaryOutput)
		config.updateBrightness(currentBrightness)
		config.dump()
	} else {
		config.readFromFile()
	}

	res := resolveBrightnessLevel(config.Brightness, floatValue)

	if res != 0 {
		config.updateBrightness(res)
		SetBrightnessLevel(config.Brightness)
		NotifySend(config.Brightness)
		config.dump()
	}
}
