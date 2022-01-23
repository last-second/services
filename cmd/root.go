package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "services",
	Short: "Commands to run the last-second backend services locally",
	Long:  "Commands to run the last-second backend services locally",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
