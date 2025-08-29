package main

import (
	"os"
	"testing"
)

func TestTmuxRunning(t *testing.T) {
	result := tmuxRunning()
	_ = result
}

func TestInTmux(t *testing.T) {
	originalTmux := os.Getenv("TMUX")
	defer func() {
		if originalTmux != "" {
			os.Setenv("TMUX", originalTmux)
		} else {
			os.Unsetenv("TMUX")
		}
	}()

	os.Unsetenv("TMUX")
	result := inTmux()
	if result {
		t.Error("inTmux() should return false when TMUX env var is not set")
	}

	os.Setenv("TMUX", "/tmp/tmux-1000/default,1234,0")
	result = inTmux()
	if !result {
		t.Error("inTmux() should return true when TMUX env var is set")
	}
}

func TestSessionExists(t *testing.T) {
	result := sessionExists("nonexistent-session-12345")
	_ = result

	result = sessionExists("")
	_ = result
}
