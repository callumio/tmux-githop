package main

import (
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func makeItems(data []string, descFunc func(string) string) []list.Item {
	items := make([]list.Item, len(data))
	for i, d := range data {
		items[i] = item{title: d, desc: descFunc(d)}
	}
	return items
}

func splitOutput(out string) []string {
	if out == "" {
		return nil
	}
	return strings.Split(out, "\n")
}
