package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/RyanTalbot/bartle/internal/config"
	"github.com/RyanTalbot/bartle/internal/lint"
	"github.com/spf13/cobra"
)

var (
	lintMsg string
)

var lintCmd = &cobra.Command{
	Use:   "lint [message-file]",
	Short: "Lint a commit message against .bartle.yaml rules",
	Long: `Validate a commit message against the style and rules defined in .bartle.yaml.

You can pass a message directly with -m/--message, a path to a message file
(e.g. .git/COMMIT_EDITMSG), or pipe a message on stdin.`,
	Example: `
  bartle lint -m "feat(ui): add dropdown"
  bartle lint .git/COMMIT_EDITMSG
  echo "fix(api): handle nil pointer" | bartle lint`,
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true, // don't print usage on lint failures
	SilenceErrors: true, // we've already printed friendly errors
	RunE: func(cmd *cobra.Command, args []string) error {
		msg := strings.TrimSpace(lintMsg)

		if msg == "" && len(args) == 1 {
			b, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("read message file: %w", err)
			}
			msg = strings.TrimSpace(string(b))
			// Strip git comment lines (# ...) commonly found in COMMIT_EDITMSG
			msg = stripGitComments(msg)
		}

		if msg == "" {
			stdinMsg, err := readStdinIfPiped()
			if err != nil {
				return fmt.Errorf("read stdin: %w", err)
			}
			if stdinMsg != "" {
				msg = stdinMsg
			}
		}

		if msg == "" {
			return errors.New("no commit message provided (use -m, a file path, or pipe on stdin)")
		}

		// Normalize CRLF
		msg = strings.ReplaceAll(msg, "\r", "")

		// Load config from repo root
		cfg, _, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		res := lint.ValidateMessage(msg, cfg)
		if res.Valid {
			fmt.Fprintln(cmd.OutOrStdout(), "✅ Commit message is valid!")
			return nil
		}

		fmt.Fprintln(cmd.OutOrStdout(), "❌ Invalid commit message:")
		for _, e := range res.Errors {
			fmt.Fprintln(cmd.OutOrStdout(), e)
		}

		// Return an error to produce non-zero exit code (hooks/CI),
		// but we've already printed the friendly output above.
		return errors.New("lint failed")
	},
}

func init() {
	lintCmd.Flags().StringVarP(&lintMsg, "message", "m", "", "commit message text to lint")
	rootCmd.AddCommand(lintCmd)
}

// readStdinIfPiped returns stdin content if data is piped; otherwise "".
func readStdinIfPiped() (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	// If stdin is a char device, nothing was piped
	if (info.Mode() & os.ModeCharDevice) != 0 {
		return "", nil
	}

	var sb strings.Builder
	reader := bufio.NewReader(os.Stdin)
	for {
		chunk, err := reader.ReadString('\n')
		sb.WriteString(chunk)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
	}

	s := strings.TrimSpace(sb.String())
	s = stripGitComments(s)

	return s, nil
}

// stripGitComments removes lines beginning with '#' which Git places in COMMIT_EDITMSG.
func stripGitComments(s string) string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "#") {
			continue
		}
		out = append(out, line)
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}
