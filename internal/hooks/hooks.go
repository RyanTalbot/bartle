package hooks

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	// Marker lets us detect whether a hook file was written by Bartle.
	BartleHookMarker = "# BARTLE-HOOK v1"
)

// RepoRootFromCwd walks upward from CWD until it finds a .git directory.
func RepoRootFromCwd() (string, error) {
	start, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working directory: %w", err)
	}
	dir := start
	for {
		if st, err := os.Stat(filepath.Join(dir, ".git")); err == nil && st.IsDir() {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("not inside a git repository (run `git init` first)")
		}
		dir = parent
	}
}

// HookPath returns the commit-msg hook path for a repo root.
func HookPath(repoRoot string) string {
	return filepath.Join(repoRoot, ".git", "hooks", "commit-msg")
}

// CheckExisting reports (exists, isOurs, err) for a hook path.
// isOurs is true if the file contains the Bartle marker.
func CheckExisting(hookPath string) (bool, bool, error) {
	data, err := os.ReadFile(hookPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false, nil
		}
		return false, false, fmt.Errorf("read existing hook: %w", err)
	}
	return true, containsMarker(string(data)), nil
}

func containsMarker(s string) bool {
	return len(s) > 0 && (func() bool {
		for i := 0; i+len(BartleHookMarker) <= len(s); i++ {
			if s[i:i+len(BartleHookMarker)] == BartleHookMarker {
				return true
			}
		}
		return false
	}())
}

// InstallCommitMsgHook writes the hook script atomically and makes it executable.
func InstallCommitMsgHook(hookPath, bartleCmd string) error {
	hooksDir := filepath.Dir(hookPath)
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return fmt.Errorf("create hooks directory: %w", err)
	}

	script := hookScript(bartleCmd)

	tmpPath := hookPath + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(script), 0o755); err != nil {
		return fmt.Errorf("write temp hook: %w", err)
	}
	if err := os.Rename(tmpPath, hookPath); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("activate hook: %w", err)
	}

	if err := os.Chmod(hookPath, 0o755); err != nil {
		return fmt.Errorf("chmod hook: %w", err)
	}

	return nil
}

func hookScript(bartleCmd string) string {
	return fmt.Sprintf(`#!/bin/sh
%s

# Bartle commit-msg hook
%s

# Pass the path to the commit message file to bartle for linting.
# If lint fails, exit non-zero to block the commit.
exec %s lint "$1"
`, "set -e", BartleHookMarker, bartleCmd)
}

func IsOurHook(hookPath string) (bool, error) {
	data, err := os.ReadFile(hookPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("read hook: %w", err)
	}
	return strings.Contains(string(data), BartleHookMarker), nil
}

// UninstallCommitMsgHook removes the commit-msg hook.
// Returns (changed, backupPath, err).
// - If the hook doesn't exist → (false, "", nil)
// - If the hook is ours → delete and return (true, "", nil)
// - If the hook is not ours:
//   - force=false  → error (won't remove other tools' hooks)
//   - force=true   → move to .bak.<timestamp> and return (true, backupPath, nil)
func UninstallCommitMsgHook(hookPath string, force bool) (bool, string, error) {
	// No hook present
	if _, err := os.Stat(hookPath); err != nil {
		if os.IsNotExist(err) {
			return false, "", nil
		}
		return false, "", fmt.Errorf("stat hook: %w", err)
	}

	isOurs, err := IsOurHook(hookPath)
	if err != nil {
		return false, "", err
	}

	// If it's ours, just remove it.
	if isOurs {
		if err := os.Remove(hookPath); err != nil {
			return false, "", fmt.Errorf("remove hook: %w", err)
		}
		return true, "", nil
	}

	// Not ours: respect --force, and back it up.
	if !force {
		return false, "", fmt.Errorf("%s exists and is not managed by Bartle (use --force to remove)", hookPath)
	}

	backupPath := hookPath + ".bak." + time.Now().Format("20060102150405")
	if err := os.Rename(hookPath, backupPath); err != nil {
		return false, "", fmt.Errorf("backup existing hook: %w", err)
	}
	return true, backupPath, nil
}
