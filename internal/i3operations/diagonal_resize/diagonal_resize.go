package diagonal_resize

import (
	"fmt"

	"go.i3wm.org/i3/v4"
)

func resize(resizeValue int) error {
	var err error

	_, err = i3.RunCommand(fmt.Sprintf("resize grow width %d px or %d ppt", resizeValue, resizeValue))
	_, err = i3.RunCommand(fmt.Sprintf("resize grow height %d px or %d ppt", resizeValue, resizeValue))
	_, err = i3.RunCommand(fmt.Sprintf("move container left %d px", resizeValue/2))
	_, err = i3.RunCommand(fmt.Sprintf("move container up %d px", resizeValue/2))

	return err
}

func Execute(resizeValue int) error {
	if err := resize(resizeValue); err != nil {
		
	}

	return nil
}
