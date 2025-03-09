package i3operations

import (
	"fmt"

	"github.com/modernpacifist/i3-scripts-go/internal/i3operations"
)

// type NodeConfig struct {
// 	i3.Node
// 	Marks []string
// }
// func NodeConfigConstructor(node *i3.Node) NodeConfig {
// 	return NodeConfig{
// 		Node:  *node,
// 		Marks: i3operations.GetNodeMarks(node),
// 	}
// }

func Execute(arg int) {
	// i3Tree := i3operations.GetI3Tree()
	// node := i3operations.GetFocusedNode(i3Tree)
	// //fmt.Println(node)
	// ndconf := containerParametersConstructor(node)
	// fmt.Println(ndconf)
	//focusedNodeMark := getNodeMark(node)
	//fmt.Println(focusedNodeMark)
	// TODO: if the mark exists, then save the preset to json <05-11-23, modernpacifist> //
	// absoluteConfigFilepath := resolveFileAbsolutePath(globalConfigFilename)
	// // TODO: need to check if the file does not exist and create it <05-11-23, modernpacifist> //
	// fmt.Println(absoluteConfigFilepath)
	// createJsonConfigFile(absoluteConfigFilepath)
	// fmt.Println(i3operations.GetOutputs())

	outputs, err := i3operations.GetOutputs()
	if err != nil {
		fmt.Println("Error getting outputs:", err)
		return
	}

	fmt.Println("Outputs: ")
	for _, output := range outputs {
		fmt.Printf("%+v\n", output)
	}
	fmt.Println()
	fmt.Println()

}
