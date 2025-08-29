package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
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

func sanitiseSessionName(name string) string {
	// Only alphanumeric, underscore, hyphen
	reg := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	sanitised := reg.ReplaceAllString(name, "_")

	// Remove consecutive underscores
	reg = regexp.MustCompile(`_+`)
	sanitised = reg.ReplaceAllString(sanitised, "_")

	// Trim underscores
	sanitised = strings.Trim(sanitised, "_")

	// Ensure not empty
	if sanitised == "" {
		sanitised = "session"
	}

	// Limit length to 250 chars
	if len(sanitised) > 250 {
		sanitised = sanitised[:250]
		sanitised = strings.TrimRight(sanitised, "_")
	}

	return sanitised
}

func confirmAction(message string) bool {
	if assumeYes {
		return true
	}
	fmt.Printf("%s (y/N): ", message)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
