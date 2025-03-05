package main

import (
	"log"
	"strconv"

	changeWorkspaceIndex "github.com/modernpacifist/i3-scripts-go/internal/i3operations/change_workspace_index"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ChangeWorkspaceIndex",
	Short: "Change the current workspace index",
	Long:  `A CLI tool to change the current workspace index`,
}

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Specify the new workspace index",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Index was not specified")
		}

		newWsIndex, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		if err := changeWorkspaceIndex.Execute(newWsIndex); err != nil {
			log.Fatal(err)
		}
	},
}

func main() {
	rootCmd.AddCommand(indexCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
