package main

import (
	"os"
	"os/exec"
)

func tmuxRunning() bool {
	return exec.Command("pgrep", "tmux").Run() == nil
}

func inTmux() bool {
	return os.Getenv("TMUX") != ""
}

func sessionExists(name string) bool {
	return exec.Command("tmux", "has-session", "-t="+name).Run() == nil
}

func switchSession() error {
	out, err := runCmd("tmux", "list-sessions", "-F", "#{session_name}")
	if err != nil {
		return err
	}

	sessions := splitOutput(out)
	if len(sessions) == 0 {
		return nil
	}

	items := makeItems(sessions, func(string) string { return "tmux session" })
	selected, err := selectItem(items, "Switch Session")
	if err != nil || selected == "" {
		return err
	}

	return exec.Command("tmux", "switch-client", "-t", selected).Run()
}
