package config

import (
	"fmt"
	"log"
	"os/user"
	"strings"
)

func ResolveFilepathHomedir(filename string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(fmt.Sprintf("~/%s", filename), "~", usr.HomeDir, 1)
}

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/user"
// 	"strings"

// 	"go.i3wm.org/i3/v4"
// )

// const (
// 	// possibly not needed
// 	ConfigDirectory string = "~/.config/"
// )

// func GenericDump[T any](data T, filename string) error {
// 	jsonData, err := json.MarshalIndent(data, "", "\t")
// 	if err != nil {
// 		return err
// 	}

// 	return os.WriteFile(filename, jsonData, 0644)
// }

// func ResolveAbsoluteFilepath(filename string) string {
// 	usr, err := user.Current()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return strings.Replace(fmt.Sprintf("~/%s", filename), "~", usr.HomeDir, 1)
// }

// func ReadConfig(filename string) (any, error) {
// 	return nil, nil
// }
