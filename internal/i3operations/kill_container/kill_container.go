package kill_container

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	common "github.com/modernpacifist/i3-scripts-go/internal/i3operations"
)

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
	var file *os.File
	var err error

	if file, err = os.Open(filePath); err != nil {
		//fmt.Println("Error opening file:", err)
		//return
	}
	defer file.Close()

	var content []byte
	if content, err = ioutil.ReadAll(file); err != nil {
		//fmt.Println("Error reading file:", err)
		//return
	}

	var jsonWhitelist JsonWhitelist
	if err = json.Unmarshal(content, &jsonWhitelist); err != nil {
		//fmt.Println("Error unmarshaling JSON:", err)
		//return
	}

	return jsonWhitelist
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func resolveJsonAbsolutePath(filename string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	return strings.Replace(fmt.Sprintf("~/%s", filename), "~", usr.HomeDir, 1)
}

func Execute() error {
	var jsonWhitelist JsonWhitelist

	jsonWhitelistPath = resolveJsonAbsolutePath(".KillContainerConfig.json")

	if fileExists(jsonWhitelistPath) == true {
		jsonWhitelist = readJsonWhitelist(jsonWhitelistPath)
	}

	focusedNode := common.GetFocusedNode()
	focusedNodeMark := common.GetNodeMark(focusedNode)
	if focusedNodeMark == "" || jsonWhitelist.checkWhitelist(focusedNodeMark) {
		return common.RunKillCommand()
	}

	userConfirmation, err := common.Runi3Input("Kill container?", 1)
	if err != nil {
		return err
	}

	if userConfirmation == "y" {
		return common.RunKillCommand()
	}

	return nil
}
