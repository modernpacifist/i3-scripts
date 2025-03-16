package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	keyboardLayout "github.com/modernpacifist/i3-scripts-go/internal/i3operations/keyboard_layout"
)

var rootCmd = &cobra.Command{
	Use:   "KeyboardLayout",
	Short: "A tool to manage keyboard layouts",
	Long:  `A CLI tool to manage keyboard layouts including cycle through layouts`,
}

var cycleCmd = &cobra.Command{
	Use:   "cycle",
	Short: "Cycle through keyboard layouts",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Help()
		}
		regex := regexp.MustCompile(`^[a-zA-Z]+(?:/[a-zA-Z]+)*$`)
		if !regex.MatchString(args[0]) {
			return fmt.Errorf("wrong input format. Use <layout1>/<layout2>/.../<layoutN>")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("No argument provided")
		}
		layoutsArray := strings.Split(args[0], "/")
		keyboardLayout.Execute(layoutsArray)
	},
}

func main() {
	rootCmd.AddCommand(cycleCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
