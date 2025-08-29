package main

import (
	"reflect"
	"testing"
)

func TestCommandExists(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		expected bool
	}{
		{"existing command", "echo", true},
		{"non-existing command", "nonexistentcmd123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := commandExists(tt.cmd)
			if result != tt.expected {
				t.Errorf("commandExists(%s) = %v, want %v", tt.cmd, result, tt.expected)
			}
		})
	}
}

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name        string
		cmd         string
		args        []string
		expectError bool
	}{
		{"successful command", "echo", []string{"hello"}, false},
		{"command with error", "false", []string{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := runCmd(tt.cmd, tt.args...)
			if tt.expectError && err == nil {
				t.Errorf("runCmd(%s, %v) expected error but got none", tt.cmd, tt.args)
			}
			if !tt.expectError && err != nil {
				t.Errorf("runCmd(%s, %v) unexpected error: %v", tt.cmd, tt.args, err)
			}
			if !tt.expectError && output == "" {
				t.Errorf("runCmd(%s, %v) expected output but got empty string", tt.cmd, tt.args)
			}
		})
	}
}

func TestMakeItems(t *testing.T) {
	data := []string{"item1", "item2", "item3"}
	descFunc := func(s string) string { return "desc for " + s }

	items := makeItems(data, descFunc)

	if len(items) != len(data) {
		t.Errorf("makeItems() returned %d items, want %d", len(items), len(data))
	}

	for i, item := range items {
		if item == nil {
			t.Errorf("item %d is nil", i)
		}
	}
}

func TestSplitOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty string", "", nil},
		{"single line", "line1", []string{"line1"}},
		{"multiple lines", "line1\nline2\nline3", []string{"line1", "line2", "line3"}},
		{"with trailing newline", "line1\nline2\n", []string{"line1", "line2", ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitOutput(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("splitOutput(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitiseSessionName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"alphanumeric", "session123", "session123"},
		{"with spaces", "my session", "my_session"},
		{"with special chars", "my@session#test", "my_session_test"},
		{"with dots", "my.session.test", "my_session_test"},
		{"consecutive underscores", "my__session", "my_session"},
		{"leading/trailing underscores", "_session_", "session"},
		{"empty string", "", "session"},
		{"too long", string(make([]byte, 300)), "session"},
		{"only special chars", "@#$%", "session"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitiseSessionName(tt.input)
			if result != tt.expected {
				t.Errorf("sanitiseSessionName(%q) = %q, want %q", tt.input, result, tt.expected)
			}

			if len(result) > 250 {
				t.Errorf("sanitiseSessionName(%q) result too long: %d chars", tt.input, len(result))
			}
			if result == "" {
				t.Errorf("sanitiseSessionName(%q) returned empty string", tt.input)
			}
		})
	}
}

func TestConfirmAction(t *testing.T) {
	originalAssumeYes := assumeYes
	defer func() { assumeYes = originalAssumeYes }()

	assumeYes = true
	result := confirmAction("Test message")
	if !result {
		t.Error("confirmAction() should return true when assumeYes is true")
	}

	assumeYes = false
}
