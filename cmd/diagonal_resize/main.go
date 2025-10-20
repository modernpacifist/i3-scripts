package main

import (
	"log"
	"regexp"
	"strconv"

	diagonalResize "github.com/modernpacifist/i3-scripts-go/internal/i3scripts/diagonal_resize"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "DiagonalResize",
	Short: "Resize the i3 floating window diagonally",
	Long:  `A CLI tool to resize the i3 floating window diagonally`,
}

var sizeCmd = &cobra.Command{
	Use:                "size [value]",
	Short:              "Resize the i3 floating window diagonally",
	DisableFlagParsing: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Help()
		}
		regex := regexp.MustCompile(`^[-+]\d+$`)
		if !regex.MatchString(args[0]) {
			log.Fatal("Wrong input format. Use +N or -N where N is a number")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("No argument provided")
		}

		resizeValue, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal("Invalid value:", args[0])
		}
		diagonalResize.Execute(resizeValue)
	},
}

func main() {
	rootCmd.AddCommand(sizeCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
