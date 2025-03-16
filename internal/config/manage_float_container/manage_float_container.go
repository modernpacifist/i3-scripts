package manage_float_container

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"go.i3wm.org/i3/v4"

	common "github.com/modernpacifist/i3-scripts-go/internal/config"
	i3operations "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
)

const (
	configFilename string = "~/.ManageFloatContainer.json"
)

type NodeConfig struct {
	// TODO: add a node id and adress some of the containers by it <13-11-23, modernpacifist> //
	ID     i3.NodeID `json:"ID"`
	X      int64     `json:"X"`
	Y      int64     `json:"Y"`
	Width  int64     `json:"Width"`
	Height int64     `json:"Height"`
	Marks  []string  `json:"Marks"`
	// TODO: having a mark field here is extremely idiotic, since we have a map with keys of marks themselves <13-11-23, modernpacifist> //
}

func NodeConfigConstructor(node i3.Node) NodeConfig {
	return NodeConfig{
		ID:     node.ID,
		X:      node.Rect.X,
		Y:      node.Rect.Y,
		Width:  node.Rect.Width,
		Height: node.Rect.Height,
		Marks:  i3operations.GetNodeMarks(node),
	}
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
