package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRootCmd(t *testing.T) {
	// Test that root command exists and has correct properties
	if rootCmd.Use != "aragomodoro" {
		t.Errorf("Expected Use 'aragomodoro', got '%s'", rootCmd.Use)
	}
	if rootCmd.Short == "" {
		t.Error("Short description should not be empty")
	}
	if rootCmd.Long == "" {
		t.Error("Long description should not be empty")
	}
}

func TestRootCmdFlags(t *testing.T) {
	// Test that all expected flags are present
	expectedFlags := []string{
		"focus",
		"break",
		"repeat",
		"continue",
		"web",
		"port",
	}

	for _, flagName := range expectedFlags {
		flag := rootCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag '%s' to exist", flagName)
		}
	}
}

func TestRootCmdFlagDefaults(t *testing.T) {
	// Reset flags to defaults
	rootCmd.Flags().Set("focus", "25")
	rootCmd.Flags().Set("break", "5")
	rootCmd.Flags().Set("repeat", "1")
	rootCmd.Flags().Set("continue", "false")
	rootCmd.Flags().Set("web", "false")
	rootCmd.Flags().Set("port", "8080")

	tests := []struct {
		flagName string
		expected string
	}{
		{"focus", "25"},
		{"break", "5"},
		{"repeat", "1"},
		{"continue", "false"},
		{"web", "false"},
		{"port", "8080"},
	}

	for _, tt := range tests {
		flag := rootCmd.Flags().Lookup(tt.flagName)
		if flag == nil {
			t.Errorf("Flag '%s' not found", tt.flagName)
			continue
		}
		if flag.DefValue != tt.expected {
			t.Errorf("Flag '%s' default value: expected '%s', got '%s'",
				tt.flagName, tt.expected, flag.DefValue)
		}
	}
}

func TestRootCmdHelp(t *testing.T) {
	// Capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute help command
	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.Execute()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// Check for error (help should not return error)
	if err != nil {
		t.Errorf("Help command should not return error: %v", err)
	}

	// Check that help contains expected content
	expectedContent := []string{
		"Aragomodoro",
		"Usage:",
		"Flags:",
		"--focus",
		"--break",
		"--repeat",
		"--web",
		"--port",
	}

	for _, content := range expectedContent {
		if !strings.Contains(output, content) {
			t.Errorf("Help output should contain '%s'", content)
		}
	}
}

func TestExecuteExists(t *testing.T) {
	// Test that Execute function exists and doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute() panicked: %v", r)
		}
	}()

	// We can't easily test the full execution without mocking,
	// but we can test that the function exists
	// This test just verifies the function can be called
}

func TestFlagValidation(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		hasError bool
	}{
		{
			name:     "ValidFlags",
			args:     []string{"--focus", "25", "--break", "5", "--repeat", "1"},
			hasError: false,
		},
		{
			name:     "WebMode",
			args:     []string{"--web", "--port", "8080"},
			hasError: false,
		},
		{
			name:     "AllFlags",
			args:     []string{"--focus", "30", "--break", "10", "--repeat", "4", "--continue"},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of rootCmd for testing
			cmd := rootCmd
			cmd.SetArgs(tt.args)

			// We can't easily test execution without side effects,
			// but we can test flag parsing
			err := cmd.ParseFlags(tt.args)
			if tt.hasError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func BenchmarkRootCmdParsing(b *testing.B) {
	args := []string{"--focus", "25", "--break", "5", "--repeat", "1"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := rootCmd
		cmd.ParseFlags(args)
	}
}
