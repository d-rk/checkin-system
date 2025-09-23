package cmd

import (
	"testing"
)

func TestMergeCommand(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		args     []string
		expected string
	}{
		{
			name:     "command with no arguments",
			command:  "ls",
			args:     []string{},
			expected: "ls",
		},
		{
			name:     "command with single argument",
			command:  "ls",
			args:     []string{"-l"},
			expected: "ls -l",
		},
		{
			name:     "command with multiple arguments",
			command:  "ls",
			args:     []string{"-l", "-a", "/home"},
			expected: "ls -l -a /home",
		},
		{
			name:     "argument with spaces gets quoted",
			command:  "echo",
			args:     []string{"hello world"},
			expected: `echo "hello world"`,
		},
		{
			name:     "multiple arguments with spaces",
			command:  "cp",
			args:     []string{"file with spaces.txt", "another file.txt"},
			expected: `cp "file with spaces.txt" "another file.txt"`,
		},
		{
			name:     "argument with quotes gets escaped",
			command:  "echo",
			args:     []string{`say "hello"`},
			expected: `echo "say \"hello\""`,
		},
		{
			name:     "mixed arguments with and without spaces",
			command:  "grep",
			args:     []string{"-r", "search term", "/path/no/spaces", "/path with spaces/"},
			expected: `grep -r "search term" /path/no/spaces "/path with spaces/"`,
		},
		{
			name:     "empty argument",
			command:  "test",
			args:     []string{""},
			expected: "test ",
		},
		{
			name:     "argument with only spaces",
			command:  "echo",
			args:     []string{"   "},
			expected: `echo "   "`,
		},
		{
			name:     "complex command with multiple quote scenarios",
			command:  "ssh",
			args:     []string{"-o", "StrictHostKeyChecking=no", "user@host", `echo "remote command" && ls -l`},
			expected: `ssh -o StrictHostKeyChecking=no user@host "echo \"remote command\" && ls -l"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeCommand(tt.command, tt.args)
			if result != tt.expected {
				t.Errorf("mergeCommand(%q, %v) = %q, want %q", tt.command, tt.args, result, tt.expected)
			}
		})
	}
}

func TestMergeCommandEdgeCases(t *testing.T) {
	t.Run("nil args slice", func(t *testing.T) {
		result := mergeCommand("command", nil)
		expected := "command"
		if result != expected {
			t.Errorf("mergeCommand(%q, nil) = %q, want %q", "command", result, expected)
		}
	})

	t.Run("empty command", func(t *testing.T) {
		result := mergeCommand("", []string{"arg1", "arg2"})
		expected := " arg1 arg2"
		if result != expected {
			t.Errorf("mergeCommand(%q, %v) = %q, want %q", "", []string{"arg1", "arg2"}, result, expected)
		}
	})

	t.Run("argument with multiple consecutive quotes", func(t *testing.T) {
		result := mergeCommand("echo", []string{`say ""hello"" world`})
		expected := `echo "say \"\"hello\"\" world"`
		if result != expected {
			t.Errorf("mergeCommand result mismatch: got %q, want %q", result, expected)
		}
	})
}
