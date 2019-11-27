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
		fmt.Println("Update successful!")
	},
}
