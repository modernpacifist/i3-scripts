package move_float_container

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"go.i3wm.org/i3/v4"

	"github.com/modernpacifist/i3-scripts-go/internal/i3scripts"

	common "github.com/modernpacifist/i3-scripts-go/internal/config"
)

const (
	configFilename string      = "~/.MoveFloatContainer.json"
	defaultPerms   os.FileMode = 0644
)

type Config struct {
	Path  string                `json:"-"`
	Nodes map[string]NodeConfig `json:"nodes"`
}

func Create() (Config, error) {
	absolutePath, err := common.ExpandHomeDir(configFilename)
	if err != nil {
		return Config{}, fmt.Errorf("resolving absolute path: %w", err)
	}

	return Config{
		Path:  absolutePath,
		Nodes: make(map[string]NodeConfig),
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

// ConfigEntity
type NodeConfig struct {
	i3.Node
	Marks []string
}

func NodeConfigConstructor(node i3.Node) NodeConfig {
	return NodeConfig{
		Node:  node,
		Marks: i3scripts.GetNodeMarks(node),
	}
}
