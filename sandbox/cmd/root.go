package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type sandbox struct {
}

var (
	rootCmd = &cobra.Command{
		Use:   "xchain-sandbox",
		Short: "A tool to start a xchain network with config params",
		Long:  ``,
	}
)

// Execute the func to start root cmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
