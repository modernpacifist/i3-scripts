package configs

import (
	"fmt"
	"log"
	"os/user"
	"strings"

	"go.i3wm.org/i3/v4"
)

const (
	ConfigDirectory string = "~/.config/"
)

type WindowConfig struct {
	ID                  int64  `json:"id"`
	ResizedPlusYFlag    bool   `json:"resizedPlusYFlag"`
	ResizedMinusYFlag   bool   `json:"resizedMinusYFlag"`
	ResizedPlusXFlag    bool   `json:"resizedPlusXFlag"`
	ResizedMinusXFlag   bool   `json:"resizedMinusXFlag"`
	X                   int64  `json:"x"`
	Y                   int64  `json:"y"`
	Width               int64  `json:"width"`
	Height              int64  `json:"height"`
	Mark                string `json:"mark"`
	PreviousPlusYValue  int64  `json:"previousPlusYValue"`
	PreviousMinusYValue int64  `json:"previousMinusYValue"`
	PreviousPlusXValue  int64  `json:"previousPlusXValue"`
	PreviousMinusXValue int64  `json:"previousMinusXValue"`
}

func NewWindowConfig(node *i3.Node) WindowConfig {
	return WindowConfig{
		ID: node.Window,
		X:  node.Rect.X,
		Y:  node.Rect.Y,
	}
}

func resolveJsonAbsolutePath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Replace(fmt.Sprintf("~/%s", filename), "~", usr.HomeDir, 1)
}
