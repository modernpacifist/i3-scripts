package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	resizeFloat "github.com/modernpacifist/i3-scripts-go/internal/i3scripts/resize_float_container"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ResizeFloat",
	Short: "Resize floating window in specified direction",
	Long:  `Resize floating window in i3 window manager using top/bottom/right/left directions`,
}

var resizeCmd = &cobra.Command{
	Use:                "resize [h/j/k/l/w] [value]",
	Short:              "Resize floating window using vim-style directions",
	Long:               `Resize floating window using h (left), j (down), k (up), l (right), w (widen), or w (shrink height) followed by a value (+/-N)`,
	DisableFlagParsing: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("requires exactly 2 arguments: direction (h/j/k/l) and value (+/-N)")
		}
		if !regexp.MustCompile(`^[hjklw]$`).MatchString(args[0]) {
			return fmt.Errorf("invalid direction: use h (left), j (down), k (up), l (right), or w (widen)")
		}
		if !regexp.MustCompile(`^[-+]\d+$`).MatchString(args[1]) {
			return fmt.Errorf("invalid value format: use +N or -N where N is a number")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		resizeValue, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid value: %s", args[1])
		}

		if err := resizeFloat.Execute(args[0], int64(resizeValue)); err != nil {
			return err
		}
		return nil
	},
}

func main() {
	rootCmd.AddCommand(resizeCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
