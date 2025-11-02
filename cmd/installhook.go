package cmd

import (
	"fmt"
	"os"

	"github.com/RyanTalbot/bartle/internal/hooks"
	"github.com/spf13/cobra"
)

var (
	hookForce       bool
	hookUseAbsolute bool
)

var installHookCmd = &cobra.Command{
	Use:   "install-hook",
	Short: "Install a git commit-msg hook that runs bartle lint",
	Long: `Installs .git/hooks/commit-msg to enforce your .bartle.yaml rules on every commit.

By default the hook invokes "bartle lint \"$1\"" (using PATH). Use --absolute to
embed the absolute path to the bartle binary in the hook for reliability.`,
	Example: `
  bartle install-hook
  bartle install-hook --absolute
  bartle install-hook --force`,
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		repoRoot, err := hooks.RepoRootFromCwd()
		if err != nil {
			return err
		}

		hookPath := hooks.HookPath(repoRoot)

		var bartleCmd string
		if hookUseAbsolute {
			exe, err := os.Executable()
			if err != nil {
				return fmt.Errorf("resolve bartle path: %w", err)
			}
			bartleCmd = exe
		} else {
			bartleCmd = "bartle"
		}

		// If a hook exists, decide whether we can/should overwrite
		exists, isOurs, err := hooks.CheckExisting(hookPath)
		if err != nil {
			return err
		}
		if exists && !hookForce && !isOurs {
			return fmt.Errorf("%s already exists and is not managed by Bartle (use --force to overwrite)", hookPath)
		}

		if err := hooks.InstallCommitMsgHook(hookPath, bartleCmd); err != nil {
			return err
		}

		fmt.Fprintln(cmd.OutOrStdout(), "✅ Installed Bartle commit-msg hook at", hookPath)
		fmt.Fprintln(cmd.OutOrStdout(), "Commits will now be linted automatically.")
		return nil
	},
}

func init() {
	installHookCmd.Flags().BoolVarP(&hookForce, "force", "f", false, "overwrite an existing hook")
	installHookCmd.Flags().BoolVarP(&hookUseAbsolute, "absolute", "a", false, "embed absolute path to bartle binary in hook")
	rootCmd.AddCommand(installHookCmd)
}
