package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/RyanTalbot/bartle/internal/templates"
	"github.com/spf13/cobra"
)

var (
	initStyle string
	initForce bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize bartle for the current project",
	Long: `Create a .bartle.yaml in the repository root with sensible defaults.
Edit the file in your editor after generation.`,
	Example: `
  bartle init
  bartle init -s jira
  bartle init -s custom -f  # combine flags separately (-s jira -f)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Ensure we’re inside a git repo
		target := repoConfigPath()
		if target == "" {
			return fmt.Errorf("not inside a git repository (run `git init` first)")
		}

		// Ensure dir exists, really shouldn't ever hit this.
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("create config directory: %w", err)
		}

		// If a file exists, stop unless --force is provided
		if _, err := os.Stat(target); err == nil && !initForce {
			return fmt.Errorf("%s already exists (use --force to overwrite)", target)
		}

		// TODO: We could add the creation of a backup file here if force is used.

		style := strings.ToLower(initStyle)
		validStyles := map[string]bool{"conventional": true, "jira": true, "custom": true}
		if !validStyles[style] {
			return fmt.Errorf("invalid --style %q (allowed: conventional|jira|custom)", initStyle)
		}

		// Render the correct template with defaults
		tmplStr := pickInitTemplate(style)

		tpl, err := template.New("cfg").Parse(tmplStr)
		if err != nil {
			return fmt.Errorf("parse template: %w", err)
		}
		var buf bytes.Buffer
		if err := tpl.Execute(&buf, defaultTemplateData()); err != nil {
			return fmt.Errorf("render template: %w", err)
		}

		// Write the config file
		if err := os.WriteFile(target, buf.Bytes(), 0o644); err != nil {
			return fmt.Errorf("write config: %w", err)
		}

		// Success message
		fmt.Println("✅ Wrote", target)
		fmt.Println("Tip: run `bartle install-hook` to enforce commit checks locally.")
		fmt.Println("Next: open .bartle.yaml in your editor to customize rules.")

		// TODO: Post init actions could happen here i.e. install hook, interactive mode.

		return nil
	},
}

func init() {
	// sensible defaults for flags
	initStyle = "conventional"

	initCmd.Flags().StringVarP(&initStyle, "style", "s", initStyle, "style: conventional|jira|custom")
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "overwrite existing config if it already exists")

	rootCmd.AddCommand(initCmd)
}

func repoConfigPath() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return filepath.Join(dir, ".bartle.yaml")
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

type templateData struct {
	AIEnabled      bool
	Model          string
	APIKey         string
	ScopeRequired  bool
	MaxLen         int
	LowercaseStart bool
	AutoApply      bool
	BlockOnFail    bool
}

func defaultTemplateData() templateData {
	return templateData{
		AIEnabled:      false,
		Model:          "gpt-5",
		APIKey:         "env:OPENAI_API_KEY",
		ScopeRequired:  true,
		MaxLen:         72,
		LowercaseStart: false,
		AutoApply:      false,
		BlockOnFail:    true,
	}
}

func pickInitTemplate(style string) string {
	switch strings.ToLower(style) {
	case "jira":
		return templates.Jira
	case "custom":
		return templates.Custom
	default:
		return templates.Conventional
	}
}
