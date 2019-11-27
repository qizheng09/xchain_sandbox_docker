package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init the sandbox enviroment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Now initializing the sandbox enviroment.Please wait...")
		fmt.Println("Initialize successful!")
	},
}
