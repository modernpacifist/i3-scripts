package resize_float_container

import (
	"errors"
	"fmt"

	common "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
)

func Execute(mode string, value int64) error {
	switch mode {
	case "l":
		return common.RunI3Command(fmt.Sprintf("resize shrink width %d px, move container right %d px", value, value))
	case "j":
		return common.RunI3Command(fmt.Sprintf("resize shrink height %d px, move container down %d px", value, value))
	case "h":
		return common.RunI3Command(fmt.Sprintf("resize grow width %d px, move container left %d px", value, value))
	case "k":
		return common.RunI3Command(fmt.Sprintf("resize grow height %d px, move container up %d px", value, value))
	default:
		return errors.New("invalid argument")
	}
}
