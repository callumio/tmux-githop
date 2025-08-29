package main

import (
	"os/exec"
	"path/filepath"
	"strings"
)

func pickSession() error {
	out, err := runCmd("ghq", "list", "-p")
	if err != nil {
		return err
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

	selectedName := strings.ReplaceAll(filepath.Base(selected), ".", "_")

	if !inTmux() && !tmuxRunning() {
		return exec.Command("tmux", "new-session", "-s", selectedName, "-c", selected).Run()
	}

	if !sessionExists(selectedName) {
		if err := exec.Command("tmux", "new-session", "-ds", selectedName, "-c", selected).Run(); err != nil {
			return err
		}
	}

	return exec.Command("tmux", "switch-client", "-t", selectedName).Run()
}
