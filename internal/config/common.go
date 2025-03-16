package config

import (
	"fmt"
	"os/user"
	"strings"
)

func ExpandHomeDir(filename string) (string, error) {
	if !strings.HasPrefix(filename, "~") {
		return filename, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	return strings.Replace(filename, "~", usr.HomeDir, 1), nil
}
