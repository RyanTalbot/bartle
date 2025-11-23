package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/RyanTalbot/bartle/internal/version"
	"github.com/spf13/cobra"
)

var (
	versionJSON  bool
	versionShort bool

	// indirection so tests can override this
	getVersion = version.Get
)

func VersionCommand() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Output bartle version information",
		Long: `Output bartle version information.

By default, prints human-readable text:
  bartle v0.1.0 (commit abc123, built 2025-08-22T19:22:33Z, by goreleaser)

Use -s for just the version number, -j for JSON, or -js for minimal JSON.`,
		Example: `
  bartle version      # human-readable
  bartle version -s   # short only
  bartle version -j   # full JSON
  bartle version -js  # minimal JSON`,
		RunE: func(cmd *cobra.Command, args []string) error {
			info := getVersion()

			switch {
			case versionJSON && versionShort:
				minimalInfo := struct {
					Version string `json:"version"`
				}{Version: info.Version}

				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetEscapeHTML(false)
				return encoder.Encode(minimalInfo)

			case versionJSON:
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetEscapeHTML(false)
				return encoder.Encode(info)

			case versionShort:
				_, err := fmt.Fprintln(cmd.OutOrStdout(), info.Version)
				return err

			default:
				_, err := fmt.Fprintf(cmd.OutOrStdout(),
					"bartle %s (commit %s, built %s, by %s)\n",
					info.Version, info.Commit, info.Date, info.BuiltBy,
				)
				return err
			}
		},
	}

	versionCmd.Flags().BoolVarP(&versionJSON, "json", "j", false, "print version information as JSON")
	versionCmd.Flags().BoolVarP(&versionShort, "short", "s", false, "print only the version number")

	return versionCmd
}

func init() {
	rootCmd.AddCommand(VersionCommand())
}
