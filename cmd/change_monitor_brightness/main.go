package main

import (
	"github.com/spf13/cobra"

	changeMonitorBrightness "github.com/modernpacifist/i3-scripts-go/internal/i3operations/change_monitor_brightness"
)

var rootCmd = &cobra.Command{
	Use:   "change-brightness",
	Short: "Change monitor brightness by a specified value",
	Run: func(cmd *cobra.Command, args []string) {
		change, err := cmd.Flags().GetFloat64("change")
		if err != nil || !cmd.Flags().Changed("change") {
			cmd.Help()
			return
		}
		changeMonitorBrightness.Execute(change)
	},
}

func init() {
	rootCmd.Flags().Float64P("change", "c", 0, "Amount to change brightness (positive or negative float)")
	rootCmd.MarkFlagRequired("change")
}

func main() {
	rootCmd.Execute()
}
