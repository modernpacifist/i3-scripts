package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"go.i3wm.org/i3/v4"
)

var i3Tree i3.Tree
var jsonWhitelistPath string

type JsonWhitelist struct {
	Whitelist []string `json:"whitelist"`
	Blacklist []string `json:"blacklist"`
}

func (jw *JsonWhitelist) checkWhitelist(mark string) bool {
	for _, m := range jw.Whitelist {
		if m == mark {
			return true
		}
	}
	return false
}

func readJsonWhitelist(filePath string) JsonWhitelist {
	file, err := os.Open(filePath)
	if err != nil {
		//fmt.Println("Error opening file:", err)
		//return
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		//fmt.Println("Error reading file:", err)
		//return
	}

	var jsonWhitelist JsonWhitelist
	err = json.Unmarshal(content, &jsonWhitelist)
	if err != nil {
		//fmt.Println("Error unmarshaling JSON:", err)
		//return
	}

	return jsonWhitelist
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func getI3Tree() i3.Tree {
	tree, err := i3.GetTree()
	if err != nil {
		log.Fatal(err)
	}

	return tree
}

func runi3Input() (string, error) {
	output, err := exec.Command("bash", "-c", "i3-input -l 1 -P \"Are you sure you want to close this window y/n?: \" | grep -oP \"output = \\K.*\" | tr -d '\n'").Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func getFocusedNode() *i3.Node {
	node := i3Tree.Root.FindFocused(func(n *i3.Node) bool {
		return n.Focused == true
	})

	if node == nil {
		log.Fatal(errors.New("Could not find focused node"))
	}

	return node
}

func getNodeMark(node *i3.Node) string {
	if len(node.Marks) == 0 {
		return ""
	}
	return node.Marks[0]
}

func resolveJsonAbsolutePath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Replace(fmt.Sprintf("~/%s", filename), "~", usr.HomeDir, 1)
}

func init() {
	i3Tree = getI3Tree()
	jsonWhitelistPath = resolveJsonAbsolutePath(".KillContainerConfig.json")
}

func main() {
	var jsonWhitelist JsonWhitelist

	if fileExists(jsonWhitelistPath) == true {
		jsonWhitelist = readJsonWhitelist(jsonWhitelistPath)
	}

	focusedNode := getFocusedNode()

	focusedContainerMark := getNodeMark(focusedNode)
	if focusedContainerMark == "" || jsonWhitelist.checkWhitelist(focusedContainerMark) {
		i3.RunCommand("kill")
		os.Exit(0)
	}

	userConfirmation, err := runi3Input()
	if err != nil {
		log.Fatal(err)
	}

	if userConfirmation == "y" {
		i3.RunCommand("kill")
	}
}
