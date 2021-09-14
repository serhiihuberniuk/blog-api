package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "Cli for blog-api",
	Long: `Cli is a command line app for blog-api.
	Available cache command only.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
