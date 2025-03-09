package change_monitor_brightness

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
)

const (
	configFilename string = "~/.ScreenDim.json"
)

type Config struct {
	Path       string  `json:"-"`
	Brightness float64 `json:"brightness"`
}

func Create() Config {
	absolutePath, err := ResolveAbsolutePath(configFilename)
	if err != nil {
		log.Fatalf("Error resolving absolute path: %v", err)
	}

	return Config{
		Path:       absolutePath,
		Brightness: 0,
	}
}

func (conf *Config) Dump() {
	jsonData, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(conf.Path, jsonData, 0644); err != nil {
		log.Fatal(err)
	}
}

func (conf *Config) Load() {
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

func (conf *Config) UpdateBrightness(newBrightness float64) {
	conf.Brightness = newBrightness
}

func ResolveAbsolutePath(filename string) (string, error) {
	if !strings.HasPrefix(filename, "~") {
		return filename, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	return strings.Replace(filename, "~", usr.HomeDir, 1), nil
}
