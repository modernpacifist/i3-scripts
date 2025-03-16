package manage_float_container

import (
	"errors"
	"fmt"

	config "github.com/modernpacifist/i3-scripts-go/internal/config/manage_float_container"
	common "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
)

func createFloatingContainer(conParams config.NodeConfig, mark string) error {
	cmd := fmt.Sprintf("mark \"%s\", move scratchpad, [con_mark=\"^%s$\"] scratchpad show, move absolute position %d %d, resize set %d %d", mark, mark, conParams.X, conParams.Y, conParams.Width, conParams.Height)
	return common.RunI3Command(cmd)
}

func createFloatingContainerDefault(conParams config.NodeConfig, mark string) error {
	cmd := fmt.Sprintf("mark \"%s\", move scratchpad, [con_mark=\"^%s$\"] scratchpad show, floating enable, resize set %d %d, move position center", mark, mark, conParams.Width, conParams.Height)
	return common.RunI3Command(cmd)
}

func showContainer(conParams config.NodeConfig) error {
	cmd := fmt.Sprintf("[con_id=%d] scratchpad show, move absolute position %d %d, resize set %d %d", conParams.ID, conParams.X, conParams.Y, conParams.Width, conParams.Height)
	return common.RunI3Command(cmd)
}

func Execute(restoreFlag string, showFlag string, updateFlag string, saveFlag bool) error {
	conf, err := config.Create()
	if err != nil {
		return err
	}

	if restoreFlag != "" {
		containerParameters, exists := conf.Nodes[restoreFlag]
		if !exists {
			focusedNode, err := common.GetFocusedNode()
			if err != nil {
				return err
			}
			containerParameters = config.NodeConfigConstructor(focusedNode)

			containerParameters.Width = 2000
			containerParameters.Height = 1000
			containerParameters.Marks = []string{restoreFlag}

			createFloatingContainerDefault(containerParameters, restoreFlag)

			conf.Nodes[restoreFlag] = containerParameters
			conf.Dump()

			return nil
		}
		createFloatingContainer(containerParameters, restoreFlag)
	}

	if showFlag != "" {
		node, exists := conf.Nodes[showFlag]
		if !exists {
			return errors.New("could not find node with mark")
		}
		if err := showContainer(node); err != nil {
			return err
		}
	}

	if updateFlag != "" {
		markedNode, err := common.GetNodeByMark(updateFlag)
		if err != nil {
			return err
		}

		var nodeConfig config.NodeConfig
		_, exists := conf.Nodes[updateFlag]
		if !exists {
			nodeConfig = config.NodeConfigConstructor(markedNode)
			conf.Nodes[updateFlag] = nodeConfig
		}

		nodeConfig = config.NodeConfigConstructor(markedNode)
		conf.Nodes[updateFlag] = nodeConfig
		conf.Dump()
	}

	if saveFlag {
		focusedNode, err := common.GetFocusedNode()
		if err != nil {
			return err
		}
		nodeConfig := config.NodeConfigConstructor(focusedNode)
		conf.Nodes[nodeConfig.Marks[0]] = nodeConfig
		conf.Dump()

		common.NotifySend(0.5, fmt.Sprintf("Saved mark %s", nodeConfig.Marks[0]))
	}

	return nil
}
