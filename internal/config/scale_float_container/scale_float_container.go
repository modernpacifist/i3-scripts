package scale_float_container

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"go.i3wm.org/i3/v4"

	common "github.com/modernpacifist/i3-scripts-go/internal/config"
)

const (
	// configFilename string      = "~/.ScaleFloatWindow.json_bak"
	configFilename string      = "~/.ScaleFloatWindow.json"
	defaultPerms   os.FileMode = 0644
)

type Config struct {
	Path   string                `json:"-"`
	Output i3.Output             `json:"-"`
	Nodes  map[int64]NodeConfig `json:"nodes"`
}

func Create(output i3.Output) (Config, error) {
	absolutePath, err := common.ExpandHomeDir(configFilename)
	if err != nil {
		return Config{}, fmt.Errorf("resolving absolute path: %w", err)
	}

	config := Config{
		Path:   absolutePath,
		Output: output,
		Nodes:  make(map[int64]NodeConfig),
	}

	if _, err := os.Stat(absolutePath); errors.Is(err, os.ErrNotExist) {
		if err := config.Dump(); err != nil {
			return Config{}, fmt.Errorf("creating initial config file: %w", err)
		}
	} else {
		if err := config.Load(); err != nil {
			return Config{}, fmt.Errorf("loading existing config: %w", err)
		}
	}

	return config, nil
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
		return errors.New("config file not found")
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

type NodeConfig struct {
	ID                  int64  `json:"id"`
	ResizedPlusYFlag    bool   `json:"resizedPlusYFlag"`
	ResizedMinusYFlag   bool   `json:"resizedMinusYFlag"`
	ResizedPlusXFlag    bool   `json:"resizedPlusXFlag"`
	ResizedMinusXFlag   bool   `json:"resizedMinusXFlag"`
	X                   int64  `json:"x"`
	Y                   int64  `json:"y"`
	Width               int64  `json:"width"`
	Height              int64  `json:"height"`
	Marks               string `json:"mark"`
	PreviousPlusYValue  int64  `json:"previousPlusYValue"`
	PreviousMinusYValue int64  `json:"previousMinusYValue"`
	PreviousPlusXValue  int64  `json:"previousPlusXValue"`
	PreviousMinusXValue int64  `json:"previousMinusXValue"`
}
