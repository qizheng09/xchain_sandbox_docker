package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop sandbox enviroment",
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: @DhunterAO
		fmt.Println("Now stopping the sandbox enviroment.Please wait...")
	},
}
