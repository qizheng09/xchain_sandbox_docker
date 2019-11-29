package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update sandbox enviroment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Now updating the sandbox enviroment.Please wait...")
		if err := updateEnv(); err != nil {
			fmt.Println("Update failed", err.Error())
		} else {
			fmt.Println("Update successful!")
		}
	},
}

func updateEnv() error {
	return nil
}
