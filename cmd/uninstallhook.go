package cmd

import (
	"fmt"

	"github.com/RyanTalbot/bartle/internal/hooks"
	"github.com/spf13/cobra"
)

var (
	uninstallForce bool
)

var uninstallHookCmd = &cobra.Command{
	Use:   "uninstall-hook",
	Short: "Uninstall Bartle's git commit-msg hook",
	Long: `Removes .git/hooks/commit-msg if it was installed by Bartle.

By default, only removes hooks that contain Bartle's marker.
Use --force to remove any existing commit-msg hook (a backup is created).`,
	Example: `
  bartle uninstall-hook
  bartle uninstall-hook --force`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		repoRoot, err := hooks.RepoRootFromCwd()
		if err != nil {
			return err
		}
		hookPath := hooks.HookPath(repoRoot)

		changed, backupPath, err := hooks.UninstallCommitMsgHook(hookPath, uninstallForce)
		if err != nil {
			return err
		}
		if !changed {
			fmt.Fprintln(cmd.OutOrStdout(), "ℹ️  No commit-msg hook to remove.")
			return nil
		}

		if backupPath != "" {
			fmt.Fprintln(cmd.OutOrStdout(), "✅ Removed hook. Backup saved at", backupPath)
		} else {
			fmt.Fprintln(cmd.OutOrStdout(), "✅ Removed Bartle commit-msg hook.")
		}
		return nil
	},
}

func init() {
	uninstallHookCmd.Flags().BoolVarP(&uninstallForce, "force", "f", false, "remove hook even if not installed by Bartle (creates a backup)")
	rootCmd.AddCommand(uninstallHookCmd)
}
