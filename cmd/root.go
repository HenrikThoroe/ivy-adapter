package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ivyagnt",
	Short: "An adapter for the Ivy game server",
	Long: "This is an adapter for the Ivy game server.\n" +
		"It allows you to connect your engine to the server and play games or run tests.\n",
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
