package change_monitor_brightness

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"
)

const (
	configFilename string      = "~/.ScreenDim.json"
	defaultPerms   os.FileMode = 0644
	minBrightness  float64     = 0
	maxBrightness  float64     = 1
)

type Config struct {
	Path       string  `json:"-"`
	Brightness float64 `json:"brightness"`
}

func Create() (Config, error) {
	absolutePath, err := ResolveAbsolutePath(configFilename)
	if err != nil {
		return Config{}, fmt.Errorf("resolving absolute path: %w", err)
	}

	return Config{
		Path:       absolutePath,
		Brightness: 0,
	}, nil
}

func (conf *Config) Dump() error {
	jsonData, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	if err := os.WriteFile(conf.Path, jsonData, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

func (conf *Config) Load() error {
	file, err := os.Open(conf.Path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("opening config file: %w", err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	if err := json.Unmarshal(content, conf); err != nil {
		return fmt.Errorf("unmarshaling JSON: %w", err)
	}
	return nil
}

func (conf *Config) UpdateBrightness(newBrightness float64) error {
	if newBrightness < minBrightness || newBrightness > maxBrightness {
		return fmt.Errorf("brightness must be between %.1f and %.1f", minBrightness, maxBrightness)
	}
	conf.Brightness = newBrightness
	return nil
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
