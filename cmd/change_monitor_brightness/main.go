package main

import (
	"github.com/spf13/cobra"

	changeMonitorBrightness "github.com/modernpacifist/i3-scripts-go/internal/i3operations/change_monitor_brightness"
)

var rootCmd = &cobra.Command{
	Use:   "change-brightness",
	Short: "Change monitor brightness by a specified value",
	Run: func(cmd *cobra.Command, args []string) {
		change, _ := cmd.Flags().GetFloat64("change")
		changeMonitorBrightness.Execute(change)
	},
}

func init() {
	rootCmd.Flags().Float64P("change", "c", 0.1, "Amount to change brightness (positive or negative float)")
}

func main() {
	rootCmd.Execute()
}
