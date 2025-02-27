package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	volumeControl "github.com/modernpacifist/i3-scripts-go/internal/i3operations/volume_control"
	"github.com/spf13/cobra"
)

// TODO: read from env <28-02-25, modernpacifist> //
const MAX_VOLUME = 100

var rootCmd = &cobra.Command{
	Use:   "VolumeControl",
	Short: "Control system volume",
	Long:  `A CLI tool to control system volume including mute/unmute and volume adjustment`,
}

var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggle mute/unmute",
	Run: func(cmd *cobra.Command, args []string) {
		volumeControl.ToggleVolume()
	},
}

var roundCmd = &cobra.Command{
	Use:   "round",
	Short: "Round volume to nearest 5",
	Run: func(cmd *cobra.Command, args []string) {
		volumeControl.RoundVolume()
	},
}

var adjustCmd = &cobra.Command{
	Use:                "adjust [+-]<number>",
	Short:              "Adjust volume up or down",
	DisableFlagParsing: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Help()
		}
		regex := regexp.MustCompile(`^[-+]\d+$`)
		fmt.Println(args)
		if !regex.MatchString(args[0]) {
			log.Fatal("Wrong input format. Use +N or -N where N is a number")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("No argument provided")
		}
		volumeControl.AdjustVolume(args[0], MAX_VOLUME)
	},
}

func main() {
	rootCmd.AddCommand(toggleCmd)
	rootCmd.AddCommand(roundCmd)
	rootCmd.AddCommand(adjustCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
