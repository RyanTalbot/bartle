package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize bartle for the current project",
	Long: `The init command initializes your project to use bartle and
creates a bartle file in the project root if it doesn't already exist.'`,
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.ReadFile("./bartle.yml")
		if err != nil {
			fmt.Println("Existing bartle file not found. Initializing.")
		}
		fmt.Print(string(file))

		f, err := os.Create("./bartle.yml")
		if err != nil {
			fmt.Println("Couldn't write file.")
		}

		defer f.Close()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
