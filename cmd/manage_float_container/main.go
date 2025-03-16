package main

import (
	"log"

	manageFloatContainer "github.com/modernpacifist/i3-scripts-go/internal/i3operations/manage_float_container"

	"github.com/spf13/cobra"
)

var (
	restoreFlag string
	showFlag    string
	updateFlag  string
	saveFlag    bool
)

var rootCmd = &cobra.Command{
	Use:   "manage-float-container",
	Short: "Manage floating containers in i3",
	RunE: func(cmd *cobra.Command, args []string) error {
		return manageFloatContainer.Execute(restoreFlag, showFlag, updateFlag, saveFlag)
	},
}

func init() {
	rootCmd.Flags().StringVar(&restoreFlag, "restore", "", "Specify the mark to restore")
	rootCmd.Flags().StringVar(&showFlag, "show", "", "Specify the mark to show")
	rootCmd.Flags().StringVar(&updateFlag, "update", "", "Specify the mark to update")
	rootCmd.Flags().BoolVar(&saveFlag, "save", false, "Save the current container state")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
