package cmd

import (
	"github.com/spf13/cobra"
)

var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint a specific commit message.",
	Long: `The lint command lints the given commit message to validity.
This can be passed in directly or piped in.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: lint implementation
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
}
