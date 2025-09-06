package cmd

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of bartle",
	Long:  `The version command gets the current version of bartle.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: version implementation
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
