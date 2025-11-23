package cmd

import (
	"bytes"
	"io"
	"testing"

	"github.com/RyanTalbot/bartle/internal/version"
	"github.com/spf13/cobra"
)

// helper to run a cobra command and capture output
func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := bytes.NewBufferString("")
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err := cmd.Execute()
	out, _ := io.ReadAll(buf)
	return string(out), err
}

func TestVersionCommandVariants(t *testing.T) {
	// stub version info for deterministic output
	origGetVersion := getVersion
	getVersion = func() version.Info {
		return version.Info{
			Version: "v0.0.0-dev",
			Commit:  "HEAD",
			Date:    "unknown",
			BuiltBy: "local",
		}
	}
	defer func() { getVersion = origGetVersion }()

	tests := []struct {
		name    string
		args    []string
		wantOut string
	}{
		{
			name:    "human readable (bartle version)",
			args:    []string{},
			wantOut: "bartle v0.0.0-dev (commit HEAD, built unknown, by local)\n",
		},
		{
			name:    "short only shorthand (bartle version -s)",
			args:    []string{"-s"},
			wantOut: "v0.0.0-dev\n",
		},
		{
			name:    "short only non-shorthand (bartle version --short)",
			args:    []string{"--short"},
			wantOut: "v0.0.0-dev\n",
		},
		{
			name: "full JSON shorthand (bartle version -j)",
			args: []string{"-j"},
			// json.Encoder adds a trailing newline
			wantOut: `{"version":"v0.0.0-dev","commit":"HEAD","date":"unknown","builtBy":"local"}
`,
		},
		{
			name: "full JSON non-shorthand (bartle version --json)",
			args: []string{"--json"},
			// json.Encoder adds a trailing newline
			wantOut: `{"version":"v0.0.0-dev","commit":"HEAD","date":"unknown","builtBy":"local"}
`,
		},
		{
			name: "minimal JSON shorthand (bartle version -js)",
			args: []string{"-js"},
			wantOut: `{"version":"v0.0.0-dev"}
`,
		},
		{
			name: "minimal JSON separated shorthand (bartle version -j -s)",
			args: []string{"-j", "-s"},
			wantOut: `{"version":"v0.0.0-dev"}
`,
		},
		{
			name: "minimal JSON separated shorthand alternate (bartle version -s -j)",
			args: []string{"-j", "-s"},
			wantOut: `{"version":"v0.0.0-dev"}
`,
		},
		{
			name: "minimal JSON separated non-shorthand (bartle version --json --short)",
			args: []string{"--json", "--short"},
			wantOut: `{"version":"v0.0.0-dev"}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := VersionCommand()

			got, err := executeCommand(cmd, tt.args...)
			if err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			if got != tt.wantOut {
				t.Fatalf("output mismatch\nwant: %q\ngot:  %q", tt.wantOut, got)
			}
		})
	}
}
