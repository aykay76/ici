package cmd
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ici",
	Short: "Local GitHub Actions runner with AI-powered pre-flight analysis",
	Long: `ici (ShiftLeft CI) - A Podman-based local GitHub Actions runner
that executes workflows on your machine before pushing to remote CI.

Features:
  - Parse and execute GitHub Actions workflows locally
  - Run in Podman containers for isolation
  - Catch issues before they reach remote CI
  - Fast feedback loop for development`,
	SilenceUsage: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug output")
}
