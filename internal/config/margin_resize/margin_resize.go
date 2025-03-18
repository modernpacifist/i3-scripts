package margin_resize

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
	configFilename string      = "~/.MarginResize.json"
	defaultPerms   os.FileMode = 0644
)

type NodeConfig struct {
	i3.Node

	DistanceToTop    int64 `json:"distance_to_top"`
	DistanceToBottom int64 `json:"distance_to_bottom"`
	DistanceToRight  int64 `json:"distance_to_right"`
	DistanceToLeft   int64 `json:"distance_to_left"`
}

type Config struct {
	Path  string                `json:"-"`
	Nodes map[string]NodeConfig `json:"nodes"`
}

func Create() (Config, error) {
	absolutePath, err := common.ExpandHomeDir(configFilename)
	if err != nil {
		return Config{}, fmt.Errorf("resolving absolute path: %w", err)
	}

	config := Config{
		Path:  absolutePath,
		Nodes: make(map[string]NodeConfig),
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
