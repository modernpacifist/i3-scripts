package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	"go.i3wm.org/i3/v4"
)

var JSONFILE string

type JsonStruct struct {
	Brightness float64 `json:"currentBrightness"`
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

func getCurrentBrigthnessJson() float64 {
	file, err := os.Open(JSONFILE)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var s JsonStruct
	err = json.NewDecoder(file).Decode(&s)
	if err != nil {
		log.Fatal(err)
	}

	return s.Brightness
}

//func changeBrightness(displays []string, currentBrightness float64) {
	//resBrightness := currentBrightness

	//for _, display := range displays {
		//c := fmt.Sprintf("xrandr --output $d --brightness $res_brightness", display, resBrightness)
		//cmd, err := exec.Command("bash", "-c", c)
	//}
//}

func dumpJson(jsonFile string, currentBrightness float64) {
	data := JsonStruct{Brightness: currentBrightness}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatal(err)
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

func resolveJsonAbsolutePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Replace("~/.screen_dim.json", "~", usr.HomeDir, 1)
}

func init() {
	JSONFILE = resolveJsonAbsolutePath()
}

func main() {
	var floatValue float64
	flag.Float64Var(&floatValue, "f", 0.0, "Brightness value")
	flag.Parse()

	fmt.Println(floatValue)

	primaryOutput := getPrimaryOutput()
	if primaryOutput == "" {
		log.Fatal("Could not get primary output")
	}

	currentBrightness := getCurrentBrigthnessXrandr(primaryOutput)

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	absJsonPath := strings.Replace(JSONFILE, "~", usr.HomeDir, 1)
	dumpJson(absJsonPath, currentBrightness)

	currentBrightnessJson := getCurrentBrigthnessJson()
	fmt.Println(currentBrightnessJson)
}
