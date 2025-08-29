package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func pickSession() error {
	out, err := runCmd("ghq", "list", "-p")
	if err != nil {
		return fmt.Errorf("failed to list GHQ repositories: %w. Make sure GHQ is properly configured and has repositories", err)
	}

	paths := splitOutput(out)
	if len(paths) == 0 {
		return nil
	}

	items := makeItems(paths, filepath.Base)
	selected, err := selectItem(items, "Select Project")
	if err != nil || selected == "" {
		return err
	}

	if info, err := os.Stat(selected); err != nil {
		return fmt.Errorf("selected path does not exist: %s", selected)
	} else if !info.IsDir() {
		return fmt.Errorf("selected path is not a directory: %s", selected)
	}

	selectedName := sanitiseSessionName(strings.ReplaceAll(filepath.Base(selected), ".", "_"))

	if !inTmux() && !tmuxRunning() {
		if !confirmAction(fmt.Sprintf("Create new tmux session '%s' in %s?", selectedName, selected)) {
			return fmt.Errorf("operation cancelled by user")
		}
		return exec.Command("tmux", "new-session", "-s", selectedName, "-c", selected).Run()
	}

	if !sessionExists(selectedName) {
		if !confirmAction(fmt.Sprintf("Create new tmux session '%s' in %s?", selectedName, selected)) {
			return fmt.Errorf("operation cancelled by user")
		}
		if err := exec.Command("tmux", "new-session", "-ds", selectedName, "-c", selected).Run(); err != nil {
			return err
		}
	}

	return exec.Command("tmux", "switch-client", "-t", selectedName).Run()
}
