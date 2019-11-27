package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clearCmd)
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the sandbox enviroment",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Now clearing the sandbox enviroment.Please wait...")
		fmt.Println("The enviroment have been cleared!")
	},
}
