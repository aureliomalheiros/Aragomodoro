package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Build the binary for testing
	cmd := exec.Command("go", "build", "-o", "aragomodoro_test", ".")
	err := cmd.Run()
	if err != nil {
		panic("Failed to build test binary: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Cleanup
	os.Remove("aragomodoro_test")
	os.Exit(code)
}

func TestCLIHelp(t *testing.T) {
	cmd := exec.Command("./aragomodoro_test", "--help")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Help command failed: %v", err)
	}

	output := out.String()
	expectedContent := []string{
		"Aragomodoro",
		"Usage:",
		"Flags:",
		"--focus",
		"--break",
		"--web",
	}

	for _, content := range expectedContent {
		if !strings.Contains(output, content) {
			t.Errorf("Help output should contain '%s'", content)
		}
	}
}

func TestCLIVersion(t *testing.T) {
	// Test that the binary runs without crashing
	cmd := exec.Command("./aragomodoro_test", "--help")
	err := cmd.Run()
	if err != nil {
		// Exit code 0 expected for help
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() != 0 {
				t.Errorf("Help command should exit with code 0, got %d", exitError.ExitCode())
			}
		}
	}
}

func TestCLIInvalidFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"InvalidFlag", []string{"--invalid"}},
		{"InvalidFlagValue", []string{"--focus", "invalid"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./aragomodoro_test", tt.args...)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr

			err := cmd.Run()
			if err == nil {
				t.Error("Expected command to fail with invalid flags")
			}
		})
	}
}

func TestWebServerStart(t *testing.T) {
	// Start web server
	cmd := exec.Command("./aragomodoro_test", "--web", "--port", "8082")
	var outBuf, errBuf bytes.Buffer
	var mu sync.Mutex

	// Create thread-safe writers
	cmd.Stdout = &threadSafeBuffer{buf: &outBuf, mu: &mu}
	cmd.Stderr = &threadSafeBuffer{buf: &errBuf, mu: &mu}

	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start web server: %v", err)
	}

	// Give server time to start
	time.Sleep(500 * time.Millisecond)

	// Kill the process
	if cmd.Process != nil {
		cmd.Process.Kill()
	}

	// Wait for process to finish
	cmd.Wait()

	// Check output safely
	mu.Lock()
	output := outBuf.String() + errBuf.String()
	mu.Unlock()

	if !strings.Contains(output, "Starting Aragomodoro web interface") {
		t.Error("Web server should print startup message")
	}
}

// Thread-safe buffer wrapper
type threadSafeBuffer struct {
	buf *bytes.Buffer
	mu  *sync.Mutex
}

func (tsb *threadSafeBuffer) Write(p []byte) (n int, err error) {
	tsb.mu.Lock()
	defer tsb.mu.Unlock()
	return tsb.buf.Write(p)
}

func TestBinaryExists(t *testing.T) {
	// Test that the test binary was created
	_, err := os.Stat("./aragomodoro_test")
	if os.IsNotExist(err) {
		t.Fatal("Test binary should exist")
	}
}

func TestBinaryPermissions(t *testing.T) {
	info, err := os.Stat("./aragomodoro_test")
	if err != nil {
		t.Fatalf("Failed to get binary info: %v", err)
	}

	mode := info.Mode()
	if mode&0111 == 0 {
		t.Error("Binary should be executable")
	}
}

func BenchmarkCLIHelp(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("./aragomodoro_test", "--help")
		cmd.Run()
	}
}
