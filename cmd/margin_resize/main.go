package main

import (
	"log"

	scaleFloatContainer "github.com/modernpacifist/i3-scripts-go/internal/i3operations/margin_resize"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "MarginResize",
	Short: "Resize window in specified direction",
	Long:  `Resize window in i3 window manager using top/bottom/right/left directions`,
}

var (
	marginResizeCmd = &cobra.Command{
		Use:   "border [direction]",
		Short: "Resize window in specified direction",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			direction := args[0]
			if err := scaleFloatContainer.Execute(direction); err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(marginResizeCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
