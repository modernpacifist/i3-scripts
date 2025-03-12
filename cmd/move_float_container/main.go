package main

import (
	"log"
	"strconv"

	moveFloatContainer "github.com/modernpacifist/i3-scripts-go/internal/i3operations/move_float_container"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "MoveFloatContainer",
	Short: "Move the i3 floating window to the bottom right",
	Long:  `A CLI tool to move the i3 floating window to the bottom right`,
}

var positionCmd = &cobra.Command{
	Use:   "position",
	Short: "Move the i3 floating window to the bottom right",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arg, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatalf("Invalid value: %s", err)
		}
		moveFloatContainer.Execute(uint8(arg))
	},
}

func init() {
	rootCmd.AddCommand(positionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
